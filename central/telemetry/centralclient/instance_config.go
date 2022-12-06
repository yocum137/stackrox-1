package centralclient

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/telemetry/phonehome"
	"github.com/stackrox/rox/pkg/version"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	apiPathsAnnotation = "rhacs.redhat.com/telemetry-apipaths"
	tenantIDLabel      = "rhacs.redhat.com/tenant"
	APICallEvent       = "API Call"
)

var (
	config = &phonehome.Config{
		ClientID: "11102e5e-ca16-4f2b-8d2e-e9e04e8dc531",
	}
	once         sync.Once
	log          = logging.LoggerForModule()
	ignoredPaths = []string{"/v1/ping", "/v1/metadata", "/static/"}
)

func getInstanceConfig() (*phonehome.Config, error) {
	rc, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(rc)
	if err != nil {
		return nil, err
	}
	v, err := clientset.ServerVersion()
	if err != nil {
		return nil, err
	}

	deployments := clientset.AppsV1().Deployments(env.Namespace.Setting())
	central, err := deployments.Get(context.Background(), "central", v1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get central deployment")
	}

	paths, ok := central.GetAnnotations()[apiPathsAnnotation]
	if !ok {
		paths = "*"
	}

	orchestrator := storage.ClusterType_KUBERNETES_CLUSTER.String()
	if env.OpenshiftAPI.BooleanSetting() {
		orchestrator = storage.ClusterType_OPENSHIFT_CLUSTER.String()
	}

	centralID := string(central.GetUID())
	tenantID := central.GetAnnotations()[tenantIDLabel]
	// Consider on-prem central a tenant of itself:
	if tenantID == "" {
		tenantID = centralID
	}

	return &phonehome.Config{
		ClientID: centralID,
		GroupID:  tenantID,
		Properties: map[string]any{
			"Central version":    version.GetMainVersion(),
			"Chart version":      version.GetChartVersion(),
			"Orchestrator":       orchestrator,
			"Kubernetes version": v.GitVersion,
		},
		Config: map[string]any{"APIPaths": set.NewFrozenSet(strings.Split(paths, ",")...)},
	}, nil
}

// InstanceConfig collects the central instance telemetry configuration from
// central Deployment annotations and orchestrator properties. The collected
// data is used for instance identification.
func InstanceConfig() *phonehome.Config {
	once.Do(func() {
		cfg, err := getInstanceConfig()
		if err != nil {
			log.Errorf("Failed to get telemetry configuration: %v. Using hardcoded values.", err)
			return
		}
		config = cfg
		log.Info("Central ID:", config.ClientID)
		log.Info("Tenant ID:", config.GroupID)
		log.Info("API path telemetry enabled for: ", config.Config["APIPaths"])

		config.AddInterceptorFunc(APICallEvent, func(rp *phonehome.RequestParams, props map[string]any) bool {
			for _, ip := range ignoredPaths {
				if strings.HasPrefix(rp.Path, ip) {
					return false
				}
			}
			apipaths := cfg.Config["APIPaths"].(set.FrozenSet[string])
			return apipaths.Contains("*") || apipaths.Contains(rp.Path)
		})

	})
	return config
}
