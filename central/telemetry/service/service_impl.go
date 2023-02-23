package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	configDS "github.com/stackrox/rox/central/config/datastore"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/telemetry/centralclient"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"google.golang.org/grpc"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		// TODO: ROX-12750 Replace DebugLogs with Administration.
		user.With(permissions.View(resources.DebugLogs)): {
			"/v1.TelemetryService/GetTelemetryConfiguration",
		},
		// TODO: ROX-12750 Replace DebugLogs with Administration.
		user.With(permissions.Modify(resources.DebugLogs)): {
			"/v1.TelemetryService/ConfigureTelemetry",
		},
		user.With(): {
			"/v1.TelemetryService/GetConfig",
		},
		user.With(permissions.Modify(resources.Administration)): {
			"/v1.TelemetryService/Disable",
			"/v1.TelemetryService/Enable",
		},
	})

	errTelemetryDisabled = errox.NotFound.New("telemetry collection is disabled")
	nothing              = &v1.Empty{}
)

type serviceImpl struct {
	v1.UnimplementedTelemetryServiceServer
}

func (s *serviceImpl) RegisterServiceServer(server *grpc.Server) {
	v1.RegisterTelemetryServiceServer(server, s)
}

func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterTelemetryServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func (s *serviceImpl) GetTelemetryConfiguration(ctx context.Context, _ *v1.Empty) (*storage.TelemetryConfiguration, error) {
	return &storage.TelemetryConfiguration{
		Enabled: false,
	}, nil
}

func (s *serviceImpl) ConfigureTelemetry(ctx context.Context, config *v1.ConfigureTelemetryRequest) (*storage.TelemetryConfiguration, error) {
	return &storage.TelemetryConfiguration{Enabled: false}, nil
}

func (s *serviceImpl) GetConfig(ctx context.Context, _ *v1.Empty) (*central.TelemetryConfig, error) {
	cfg := centralclient.InstanceConfig()
	if !cfg.Enabled() {
		return nil, errTelemetryDisabled
	}
	id, err := authn.IdentityFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &central.TelemetryConfig{
		UserId:       cfg.HashUserAuthID(id),
		Endpoint:     cfg.Endpoint,
		StorageKeyV1: cfg.StorageKey,
	}, nil
}

func updateTelemetryEnabled(ctx context.Context, enable bool) error {
	config, err := configDS.Singleton().GetConfig(ctx)
	if err != nil {
		return err
	}
	if config == nil {
		config = &storage.Config{}
	}
	pc := config.GetPublicConfig()
	if pc == nil {
		pc = &storage.PublicConfig{}
		config.PublicConfig = pc
	}
	t := pc.GetTelemetry()
	if t == nil {
		t = &storage.TelemetryConfiguration{Enabled: true}
		pc.Telemetry = t
	}
	if t.Enabled != enable {
		t.Enabled = enable
		err = configDS.Singleton().UpsertConfig(ctx, config)
	}
	return err
}

func (s *serviceImpl) Disable(ctx context.Context, _ *v1.Empty) (*v1.Empty, error) {
	if !centralclient.InstanceConfig().Enabled() {
		return nothing, nil
	}
	centralclient.Disable()

	if err := updateTelemetryEnabled(ctx, false); err != nil {
		return nil, err
	}
	return nothing, nil
}

func (s *serviceImpl) Enable(ctx context.Context, _ *v1.Empty) (*v1.Empty, error) {
	if !centralclient.Enable().Enabled() {
		return nil, errTelemetryDisabled
	}

	if err := updateTelemetryEnabled(ctx, true); err != nil {
		return nil, err
	}
	return nothing, nil
}
