package integration_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/grpc/authn"
	authnMocks "github.com/stackrox/rox/pkg/grpc/authn/mocks"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/grpc/util"
	"github.com/stackrox/rox/pkg/probeupload/mocks"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/sensor/common/connection"
	"github.com/stackrox/rox/sensor/kubernetes/client"
	"github.com/stackrox/rox/sensor/kubernetes/sensor"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

type fakeService struct{
	stream central.SensorService_CommunicateServer
}

func (s *fakeService) Communicate(msg central.SensorService_CommunicateServer) error {
	fmt.Println("Sensor communicate with fake central")
	s.stream = msg

	go func() {
		msg, err := s.stream.Recv()
		if err != nil {
			return
		}
		fmt.Printf("message received: %s\n", msg)
	}()
	return nil
}

func createConnectionAndStartServer(t *testing.T) *grpc.ClientConn {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	fakeCentral := &fakeService{}
	server := grpc.NewServer()
	central.RegisterSensorServiceServer(server, fakeCentral)

	go func() {
		utils.IgnoreError(func() error {
			return server.Serve(listener)
		})
	}()
	return conn
}

type fakeConn struct {
	conn *grpc.ClientConn
	stopSig concurrency.ErrorSignal
	okSig concurrency.Signal
	mockedIdentity authn.Identity
}

type fakeIdentiyExtractor struct{
	mockedIdentity authn.Identity
}

func (f *fakeIdentiyExtractor) IdentityForRequest(_ context.Context, _ requestinfo.RequestInfo) (authn.Identity, error) {
	return f.mockedIdentity, nil
}

func (c *fakeConn) MtlsServiceIdExtractor() (authn.IdentityExtractor, error) {
	return &fakeIdentiyExtractor{c.mockedIdentity}, nil
}

func (c *fakeConn) SetCentralConnectionWithRetries(ptr *util.LazyClientConn) {
	ptr.Set(c.conn)
}

func (c *fakeConn) StopSignal() concurrency.ErrorSignal {
	return c.stopSig
}

func (c *fakeConn) OkSignal() concurrency.Signal {
	return c.okSig
}

func (c *fakeConn) signal() bool {
	return c.okSig.Signal()
}

func makeFakeConnectionFactory(c *grpc.ClientConn, mockedIdentity authn.Identity) connection.ConnectionFactory {
	fakeConnection := &fakeConn{
		conn:    c,
		stopSig: concurrency.NewErrorSignal(),
		okSig:   concurrency.NewSignal(),
		mockedIdentity: mockedIdentity,
	}
	fakeConnection.signal()
	return fakeConnection
}

func Test_SensorHelloHandshake(t *testing.T) {
	// 1. Start sensor
	// 2. Assert that sensor hello is sent over a fake gRPC client
	// 3. Mock a central hello message in the fake gRPC client

	conn := createConnectionAndStartServer(t)
	c := fake.NewSimpleClientset()
	fakeInterface := client.InterfaceFromK8s(c)

	mockCtrl := gomock.NewController(t)

	psMock := mocks.NewMockProbeSource(mockCtrl)

	_, err := fakeInterface.Kubernetes().CoreV1().Nodes().Create(context.Background(), &v1.Node{
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

	require.NoError(t, err)

	mockedIdentity := authnMocks.NewMockIdentity(mockCtrl)
	fakeConnectionFactory := makeFakeConnectionFactory(conn, mockedIdentity)
	s, err := sensor.CreateSensor(fakeInterface, nil, fakeConnectionFactory, psMock, true)
	require.NoError(t, err)

	go s.Start()
	time.Sleep(10 * time.Second)
	s.Stop()
}
