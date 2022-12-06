package phonehome

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
	erroxGRPC "github.com/stackrox/rox/pkg/errox/grpc"
	"github.com/stackrox/rox/pkg/grpc/authn"
	grpcError "github.com/stackrox/rox/pkg/grpc/errors"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/logging"
	"google.golang.org/grpc"
)

// RequestParams holds intercepted call parameters.
type RequestParams struct {
	UserAgent string
	UserID    string
	Path      string
	Code      int
	GrpcReq   any
	HttpReq   *http.Request
}

var (
	log = logging.LoggerForModule()
	mux = &sync.Mutex{}
)

// AddInterceptorFunc appends the custom list of telemetry interceptors with the
// provided function.
func (cfg *Config) AddInterceptorFunc(event string, f interceptor) {
	mux.Lock()
	defer mux.Unlock()
	if cfg.interceptors == nil {
		cfg.interceptors = make(map[string][]interceptor, 1)
	}
	cfg.interceptors[event] = append(cfg.interceptors[event], f)
}

func (cfg *Config) track(rp *RequestParams, t Telemeter) {
	for event, is := range cfg.interceptors {
		props := map[string]any{}
		ok := true
		for _, interceptor := range is {
			if ok = interceptor(rp, props); !ok {
				break
			}
		}
		if ok {
			t.Track(event, rp.UserID, props)
		}
	}
}

func (cfg *Config) getGrpcRequestDetails(ctx context.Context, err error, info *grpc.UnaryServerInfo, req any) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	ri := requestinfo.FromContext(ctx)
	return &RequestParams{
		UserAgent: strings.Join(ri.Metadata.Get("User-Agent"), ", "),
		UserID:    cfg.HashUserID(id),
		Path:      info.FullMethod,
		Code:      int(erroxGRPC.RoxErrorToGRPCCode(err)),
		GrpcReq:   req,
	}
}

func (cfg *Config) getHttpRequestDetails(ctx context.Context, r *http.Request, err error) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	return &RequestParams{
		UserAgent: strings.Join(r.Header.Values("User-Agent"), ", "),
		UserID:    cfg.HashUserID(id),
		Path:      r.URL.Path,
		Code:      grpcError.ErrToHTTPStatus(err),
		HttpReq:   r,
	}
}

// getGRPCInterceptor returns an API interceptor function for GRPC requests.
func (cfg *Config) getGRPCInterceptor(t Telemeter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		rp := cfg.getGrpcRequestDetails(ctx, err, info, req)
		go cfg.track(rp, t)
		return resp, err
	}
}

func statusCodeToError(code *int) error {
	if code == nil || *code == http.StatusOK {
		return nil
	}
	return errors.Errorf("%d %s", *code, http.StatusText(*code))
}

// getHTTPInterceptor returns an API interceptor function for HTTP requests.
func (cfg *Config) getHTTPInterceptor(t Telemeter) httputil.HTTPInterceptor {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			statusTrackingWriter := httputil.NewStatusTrackingWriter(w)
			handler.ServeHTTP(statusTrackingWriter, r)
			rp := cfg.getHttpRequestDetails(r.Context(), r, statusCodeToError(statusTrackingWriter.GetStatusCode()))
			go cfg.track(rp, t)
		})
	}
}

// MakeInterceptors returns a couple of interceptors.
func (cfg *Config) MakeInterceptors() (grpc.UnaryServerInterceptor, httputil.HTTPInterceptor) {
	t := cfg.TelemeterSingleton()
	return cfg.getGRPCInterceptor(t), cfg.getHTTPInterceptor(t)
}
