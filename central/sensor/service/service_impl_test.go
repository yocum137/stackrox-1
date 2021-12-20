package service

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/cluster/datastore/mocks"
	"github.com/stackrox/rox/central/sensor/service/connection"
	"github.com/stackrox/rox/central/sensor/service/pipeline/all"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestIt(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	manager := connection.ManagerSingleton()
	clusterMockDS := mocks.NewMockDataStore(mockCtrl)

	list, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := grpc.NewServer()

	serverSignal := concurrency.Signal{}

	conn, err := grpc.Dial(list.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, conn)

	s := New(manager, all.Singleton(), clusterMockDS)
	assert.NotNil(t, s)

	go func() {
		defer server.Stop()
		err = server.Serve(list)
		require.NoError(t, err)
	}()

	hello := &central.SensorHello{}
	reply := &central.CentralHello{}
	err = conn.Invoke(context.TODO(), "/", hello, reply)
	require.NoError(t, err)

	select {
	case <-serverSignal.WaitC():
	case <-time.Tick(5*time.Second):
		fmt.Println("Stopping server after timeout")
		server.Stop()
	}
}
