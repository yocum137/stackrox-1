package phonehome

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/telemetry/phonehome/mocks"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type interceptorTestSuite struct {
	suite.Suite

	mockTelemeter *mocks.MockTelemeter
}

var _ suite.SetupTestSuite = (*interceptorTestSuite)(nil)

func TestInterceptor(t *testing.T) {
	suite.Run(t, new(interceptorTestSuite))
}

func (s *interceptorTestSuite) SetupTest() {
	s.mockTelemeter = mocks.NewMockTelemeter(gomock.NewController(s.T()))
}

type testRequest struct {
	value string
}

func (s *interceptorTestSuite) TestAddGrpcInterceptor() {
	testRP := &RequestParams{
		Path:      "/v1.Abc",
		Code:      0,
		UserAgent: "test",
		UserID:    "local:test:unauthenticated",
		GrpcReq: &testRequest{
			value: "test value",
		},
	}
	cfg := &Config{
		ClientID: "test",
		Config:   map[string]any{"APIPaths": set.NewFrozenSet(testRP.Path)},
	}

	cfg.AddInterceptorFunc("TestEvent", func(rp *RequestParams, props map[string]any) bool {
		if rp.Path == testRP.Path {
			if tr, ok := rp.GrpcReq.(*testRequest); ok {
				props["Property"] = tr.value
			}
		}
		return true
	})

	s.mockTelemeter.EXPECT().Track("TestEvent", testRP.UserID, map[string]any{
		"Path":       testRP.Path,
		"Code":       testRP.Code,
		"User-Agent": testRP.UserAgent,
		"Property":   "test value",
	}).Times(1)

	cfg.track(testRP, s.mockTelemeter)
}

func (s *interceptorTestSuite) TestAddHttpInterceptor() {
	testRP := &RequestParams{
		Path:      "/v1/abc",
		Code:      200,
		UserAgent: "test",
		UserID:    "local:test:unauthenticated",
	}
	req, err := http.NewRequest(http.MethodPost, "https://test"+testRP.Path+"?test_key=test_value", nil)
	s.NoError(err)
	testRP.HttpReq = req
	cfg := &Config{
		ClientID: "test",
		Config:   map[string]any{"APIPaths": set.NewFrozenSet(testRP.Path)},
	}

	cfg.AddInterceptorFunc("TestEvent", func(rp *RequestParams, props map[string]any) bool {
		if rp.Path == testRP.Path {
			props["Property"] = rp.HttpReq.FormValue("test_key")
		}
		return true
	})

	s.mockTelemeter.EXPECT().Track("TestEvent", testRP.UserID, map[string]any{
		"Path":       testRP.Path,
		"Code":       testRP.Code,
		"User-Agent": testRP.UserAgent,
		"Property":   "test_value",
	}).Times(1)

	cfg.track(testRP, s.mockTelemeter)
}

func (s *interceptorTestSuite) TestGrpcRequestInfo() {
	testRP := &RequestParams{
		UserID:    "local:test:unauthenticated",
		Code:      0,
		UserAgent: "test",
		Path:      "/v1.Test",
	}
	cfg := &Config{
		ClientID: "test",
		Config:   map[string]any{"APIPaths": set.NewFrozenSet(testRP.Path)},
	}

	md := metadata.New(nil)
	md.Set("User-Agent", testRP.UserAgent)
	ctx := peer.NewContext(context.Background(), &peer.Peer{Addr: &net.UnixAddr{Net: "pipe"}})

	rih := requestinfo.NewRequestInfoHandler()
	ctx, err := rih.UpdateContextForGRPC(metadata.NewIncomingContext(ctx, md))
	s.NoError(err)

	rp := cfg.getGrpcRequestDetails(ctx, err, &grpc.UnaryServerInfo{
		FullMethod: testRP.Path,
	}, "request")
	s.Equal(testRP.Path, rp.Path)
	s.Equal(testRP.Code, rp.Code)
	s.Equal(testRP.UserAgent, rp.UserAgent)
	s.Equal(testRP.UserID, rp.UserID)
	s.Equal("request", rp.GrpcReq)
}

func (s *interceptorTestSuite) TestHttpRequestInfo() {
	testRP := &RequestParams{
		UserID:    "local:test:unauthenticated",
		Code:      200,
		UserAgent: "test",
		Path:      "/v1/test",
	}
	cfg := &Config{
		ClientID: "test",
		Config:   map[string]any{"APIPaths": set.NewFrozenSet(testRP.Path)},
	}

	req, err := http.NewRequest(http.MethodPost, "https://test"+testRP.Path+"?test_key=test_value", nil)
	s.NoError(err)
	req.Header.Add("User-Agent", testRP.UserAgent)

	ctx := context.Background()
	rp := cfg.getHttpRequestDetails(ctx, req, err)
	s.Equal(testRP.Path, rp.Path)
	s.Equal(testRP.Code, rp.Code)
	s.Equal(testRP.UserAgent, rp.UserAgent)
	s.Equal(testRP.UserID, rp.UserID)
}
