package phonehome

import (
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	erroxGRPC "github.com/stackrox/rox/pkg/errox/grpc"
	"github.com/stackrox/rox/pkg/grpc/authn"
	grpcError "github.com/stackrox/rox/pkg/grpc/errors"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/httputil"
	pkgPH "github.com/stackrox/rox/pkg/telemetry/phonehome"
	"google.golang.org/grpc"
)

const APICallEvent = "API Call"

// RequestParams holds intercepted call parameters.
type RequestParams struct {
	UserAgent string
	UserID    string
	Path      string
	Code      int
	GrpcReq   any
	HttpReq   *http.Request
}

// interceptor is a function which will be called on every API call if none of
// the previous interceptors in the chain returned false.
// An interceptor function may add custom properties to the props map so that
// they appear in the event.
type interceptor func(rp *RequestParams, props map[string]any) bool

var (
	ignoredPaths = []string{"/v1/ping", "/v1/metadata", "/static/"}
	interceptors = map[string][]interceptor{}
)

// AddInterceptorFunc appends the custom list of telemetry interceptors with the
// provided function.
func AddInterceptorFunc(event string, f interceptor) {
	interceptors[event] = append(interceptors[event], f)
}

func track(rp *RequestParams, t pkgPH.Telemeter) {
	for event, is := range interceptors {
		props := map[string]any{
			"Path":       rp.Path,
			"Code":       rp.Code,
			"User-Agent": rp.UserAgent,
		}
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

func getGrpcRequestDetails(ctx context.Context, err error, info *grpc.UnaryServerInfo, req any) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	ri := requestinfo.FromContext(ctx)
	return &RequestParams{
		UserAgent: strings.Join(ri.Metadata.Get("User-Agent"), ", "),
		UserID:    pkgPH.HashUserID(id),
		Path:      info.FullMethod,
		Code:      int(erroxGRPC.RoxErrorToGRPCCode(err)),
		GrpcReq:   req,
	}
}

func getHttpRequestDetails(ctx context.Context, r *http.Request, err error) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	return &RequestParams{
		UserAgent: strings.Join(r.Header.Values("User-Agent"), ", "),
		UserID:    pkgPH.HashUserID(id),
		Path:      r.URL.Path,
		Code:      grpcError.ErrToHTTPStatus(err),
		HttpReq:   r,
	}
}

// getGRPCInterceptor returns an API interceptor function for GRPC requests.
func getGRPCInterceptor(t pkgPH.Telemeter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		rp := getGrpcRequestDetails(ctx, err, info, req)
		go track(rp, t)
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
func getHTTPInterceptor(t pkgPH.Telemeter) httputil.HTTPInterceptor {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			statusTrackingWriter := httputil.NewStatusTrackingWriter(w)
			handler.ServeHTTP(statusTrackingWriter, r)
			rp := getHttpRequestDetails(r.Context(), r, statusCodeToError(statusTrackingWriter.GetStatusCode()))
			go track(rp, t)
		})
	}
}

// MakeInterceptorsWithConfig returns a couple of interceptors initialized with
// the provided configuration.
func MakeInterceptorsWithConfig(cfg *pkgPH.Config, t pkgPH.Telemeter) (grpc.UnaryServerInterceptor, httputil.HTTPInterceptor) {
	AddInterceptorFunc(APICallEvent, func(rp *RequestParams, props map[string]any) bool {
		for _, ip := range ignoredPaths {
			if strings.HasPrefix(rp.Path, ip) {
				return false
			}
		}
		return cfg.APIPaths.Contains("*") || cfg.APIPaths.Contains(rp.Path)
	})

	return getGRPCInterceptor(t), getHTTPInterceptor(t)
}

// MakeInterceptors returns a couple of interceptors initialized with
// configuration and telemeter singletons.
func MakeInterceptors() (grpc.UnaryServerInterceptor, httputil.HTTPInterceptor) {
	return MakeInterceptorsWithConfig(pkgPH.InstanceConfig(), pkgPH.TelemeterSingleton())
}
