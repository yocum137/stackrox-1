package service

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/central/cve/common"
	"github.com/stackrox/rox/central/cve/datastore"
	"github.com/stackrox/rox/central/role/resources"
	vulnReqMgr "github.com/stackrox/rox/central/vulnerabilityrequest/manager/requestmgr"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/dackbox/utils/queue"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/and"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"google.golang.org/grpc"
)

var (
	authorizer = func() authz.Authorizer {
		return perrpc.FromMap(map[authz.Authorizer][]string{
			and.And(
				user.With(permissions.Modify(resources.VulnerabilityManagementRequests)),
				user.With(permissions.Modify(resources.VulnerabilityManagementApprovals))): {
				"/v1.CVEService/SuppressCVEs",
				"/v1.CVEService/UnsuppressCVEs",
			},
		})
	}()
)

// serviceImpl provides APIs for CVEs.
type serviceImpl struct {
	cves       datastore.DataStore
	vulnReqMgr vulnReqMgr.Manager
	indexQ     queue.WaitableQueue
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterCVEServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterCVEServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

// SuppressCVEs suppresses CVEs for specific duration or indefinitely.
func (s *serviceImpl) SuppressCVEs(ctx context.Context, request *v1.SuppressCVERequest) (*v1.Empty, error) {
	createdAt := types.TimestampNow()
	if err := s.cves.Suppress(ctx, createdAt, request.GetDuration(), request.GetIds()...); err != nil {
		return nil, err
	}
	// This handles updating image-cve edges and reprocessing affected deployments.
	if err := s.vulnReqMgr.SnoozeVulnerabilityOnRequest(ctx, common.SuppressCVEReqToVulnReq(request, createdAt)); err != nil {
		log.Error(err)
	}
	return &v1.Empty{}, nil
}

// UnsuppressCVEs unsuppresses given CVEs indefinitely.
func (s *serviceImpl) UnsuppressCVEs(ctx context.Context, request *v1.UnsuppressCVERequest) (*v1.Empty, error) {
	if err := s.cves.Unsuppress(ctx, request.GetIds()...); err != nil {
		return nil, err
	}
	// This handles updating image-cve edges and reprocessing affected deployments.
	if err := s.vulnReqMgr.UnSnoozeVulnerabilityOnRequest(ctx, common.UnSuppressCVEReqToVulnReq(request)); err != nil {
		log.Error(err)
	}
	return &v1.Empty{}, nil
}
