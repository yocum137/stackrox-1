package integration

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/pkg/errors"
	v12 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stackrox/rox/sensor/kubernetes/client"
	"github.com/stackrox/rox/sensor/kubernetes/sensor"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// 1. Create central fake gRPC connection that will spy on sensors events
// 2. Generate fake k8s API workloads
// 3. Try Ginko library for creating acceptance tests

type StreamMock struct {
	grpc.ClientStream
	ctx context.Context
	sentMessages chan *central.MsgFromSensor
}

func makeStreamMock() *StreamMock {
	return &StreamMock{
		ctx: context.Background(),
		sentMessages: make(chan *central.MsgFromSensor, 10),
	}
}

func (m *StreamMock) Context() context.Context {
	return m.ctx
}

func (m *StreamMock) Send(resp *central.MsgFromSensor) error {
	log.Printf("Sent message")
	m.sentMessages <- resp
	return nil
}

func (m *StreamMock) Recv() (*central.MsgFromSensor, error) {
	log.Printf("Received message")
	resp, more := <-m.sentMessages
	if !more {
		return nil, errors.New("empty")
	}
	return resp, nil
}

type fakeService struct{
	stream central.SensorService_CommunicateServer
}

func (s *fakeService) GetMetadata(ctx context.Context, in *v12.Empty) (*v12.Metadata, error) {
	log.Printf("GetMetadata")
	return &v12.Metadata{
		Version:              "1.2.3",
		BuildFlavor:          "development_build",
		ReleaseBuild:         false,
		LicenseStatus:        0,
	}, nil
}

func (s *fakeService) TLSChallenge(ctx context.Context, in *v12.TLSChallengeRequest) (*v12.TLSChallengeResponse, error) {
	log.Println("TLSChallenge")
	return &v12.TLSChallengeResponse{
		TrustInfoSerialized:  make([]byte, 30),
		Signature:            make([]byte, 30),
	}, nil
}

func (s *fakeService) Communicate(msg central.SensorService_CommunicateServer) error {
	// return makeStreamMock(), nil
	fmt.Println("Sensor communicate with fake central")
	s.stream = msg

	go func() {
		msg, err := s.stream.Recv()
		if err != nil {
			return
		}
		fmt.Println(msg.String())
	}()
	return nil
}

func TestExample(t *testing.T) {
	isolator := envisolator.NewEnvIsolator(t)

	isolator.Setenv("ROX_MTLS_CERT_FILE", "certs/cert.pem")
	isolator.Setenv("ROX_MTLS_KEY_FILE", "certs/key.pem")
	isolator.Setenv("ROX_MTLS_CA_FILE", "certs/ca.pem")
	isolator.Setenv("ROX_MTLS_CA_KEY_FILE", "certs/caKey.pem")

	log.Println("RUNNING TEST")

	defer isolator.RestoreAll()

	// ctx := context.Background()

	fakeCentral := &fakeService{}
	l, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		t.Fatal(err)
	}

	grpcv := grpc.NewServer()
	central.RegisterSensorServiceServer(grpcv, fakeCentral)
	v12.RegisterMetadataServiceServer(grpcv, fakeCentral)

	go func() {
		if err := grpcv.Serve(l); err != nil {
			panic(err)
		}
	}()

	isolator.Setenv("ROX_CENTRAL_ENDPOINT", l.Addr().String())


	c := fake.NewSimpleClientset()
	fakeInterface := client.MustCreateInferfaceFromK8s(c)

	_, err = fakeInterface.Kubernetes().CoreV1().Nodes().Create(context.Background(), &v1.Node{
		Spec:       v1.NodeSpec{
			PodCIDR:            "",
			PodCIDRs:           nil,
			ProviderID:         "",
			Unschedulable:      false,
			Taints:             nil,
		},
		Status:     v1.NodeStatus{
			Capacity:        nil,
			Allocatable:     nil,
			Phase:           "",
			Conditions:      nil,
			Addresses:       nil,
			DaemonEndpoints: v1.NodeDaemonEndpoints{},
			NodeInfo:        v1.NodeSystemInfo{},
			Images:          nil,
			VolumesInUse:    nil,
			VolumesAttached: nil,
			Config:          nil,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		t.Fatal(err)
	}
	
	fakeSensor, err := sensor.CreateSensor(fakeInterface, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	go fakeSensor.Start()
	time.Sleep(30 * time.Second)
	fakeSensor.Stop()
}
