package spike

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	roxv1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

var (
	// How well does this work with parallel tests?
	// namespaceName = envconf.RandomName("qa-ns", 16)
	policyName    = envconf.RandomName("test-policy", 16)
	testdata      = os.DirFS("testdata")

	global env.Environment

	namespaceName string
	namespaceObj  *v1.Namespace
	policyService roxv1.PolicyServiceClient
	alertService  roxv1.AlertServiceClient
)

func TestMain(m *testing.M) {
	mainConfig := envconf.New().WithKubeconfigFile("/Users/fbittenc/workspaces/fred-test-artifacts/kubeconfig")
	global = env.NewWithConfig(mainConfig)

	global.BeforeEachFeature(func(ctx context.Context, config *envconf.Config, t *testing.T, _ features.Feature) (context.Context, error) {
		// Setup rox connection
		conn := testutils.GRPCConnectionToCentral(t)
		policyService = roxv1.NewPolicyServiceClient(conn)
		alertService = roxv1.NewAlertServiceClient(conn)

		// Create namespace for each feature
		r, err := resources.New(config.Client().RESTConfig())
		if err != nil {
			return ctx, err
		}

		namespaceName = envconf.RandomName("qa-ns", 16)
		namespaceObj = createTestNs(t, ctx, r, namespaceName)
		return ctx, nil
	})

	global.AfterEachFeature(func(ctx context.Context, config *envconf.Config, t *testing.T, _ features.Feature) (context.Context, error) {
		// Delete namespace for each feature
		r, err := resources.New(config.Client().RESTConfig())
		if err != nil {
			return ctx, err
		}
		if err = r.Delete(ctx, namespaceObj); err != nil {
			return ctx, err
		}
		return ctx, nil
	})

	os.Exit(global.Run(m))
}

func createTestNs(t *testing.T, ctx context.Context, r *resources.Resources, name string) *v1.Namespace {
	nsObj := v1.Namespace{}
	nsObj.Name = name
	if err := r.Create(ctx, &nsObj); err != nil {
		t.Fatal(err)
	}

	return &nsObj
}

func mustGetResources(t *testing.T, cfg *envconf.Config) *resources.Resources {
	r, err := resources.New(cfg.Client().RESTConfig())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func mustDeleteAll(t *testing.T, ctx context.Context, r *resources.Resources, objs ...k8s.Object) {
	for _, d := range objs {
		err := r.Delete(ctx, d)
		if err != nil {
			t.Errorf("failed to cleanup resource %s", d.GetName())
		}
	}
}

func createPolicyIfMissing(req *roxv1.PostPolicyRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	p, err := policyService.PostPolicy(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.Internal && strings.Contains(s.Message(), "Could not add policy due to name validation") {
				return p.Id, nil
			}
		}
		return "", err
	}
	return p.Id, nil
}

func deletePolicy(t *testing.T, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := policyService.DeletePolicy(ctx, &roxv1.ResourceByID{Id: id})
	assert.NoError(t, err)
}

func getAlertsForPolicy(t testutils.T, policyName string) []*storage.ListAlert {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var alerts []*storage.ListAlert
	resp, err := alertService.ListAlerts(ctx, &roxv1.ListAlertsRequest{
		Query: fmt.Sprintf("Policy:%s+Namespace:%s", policyName, namespaceName),
	})

	if err != nil {
		t.Errorf("Failed to get alerts: %s", err)
		t.FailNow()
	}

	alerts = resp.GetAlerts()
	return alerts
}

func deploymentReplicas(t *testing.T) func(object k8s.Object) int32 {
	return func(object k8s.Object) int32 {
		if dep, ok := object.(*appsV1.Deployment); !ok {
			t.Fatal("cannot convert k8s object to deployment")
			return -1
		} else {
			return *dep.Spec.Replicas
		}
	}
}

func applyFile(t *testing.T, ctx context.Context, cfg *envconf.Config, filename string, object k8s.Object) {
	r := mustGetResources(t, cfg)
	if err := decoder.DecodeFile(
		testdata,
		filename,
		object,
		decoder.MutateNamespace(namespaceName),
	); err != nil {
		t.Fatal(err)
	}

	if err := r.Create(ctx, object); err != nil {
		t.Fatal(err)
	}
}

func patchAndWaitForReplicas(t *testing.T, ctx context.Context, cfg *envconf.Config, object k8s.Object, count int32) {
	r, err := resources.New(cfg.Client().RESTConfig())
	if err != nil {
		t.Fatal(err)
	}

	mergePatch, err := json.Marshal(map[string]interface{}{
		"spec": map[string]interface{}{
			"replicas": count,
		},
	})

	if err != nil 	{
		t.Fatal(err)
	}

	err = r.Patch(ctx, object, k8s.Patch{
		PatchType: types.StrategicMergePatchType,
		Data:      mergePatch,
	})

	if err != nil {
		t.Fatal(err)
	}

	if err := wait.For(conditions.New(r).ResourceScaled(object, deploymentReplicas(t), count)); err != nil {
		t.Fatal(err)
	}
}

func Test_NetworkPolicy(t *testing.T) {

	testCases := map[string]struct {
		policyName        string
		policyField       string
		networkPolicyFile string
	}{
		"Should have ingress": {
			policyField:       "Missing Ingress Network Policy",
			networkPolicyFile: "allow-ingress-netpol.yaml",
		},
		"Should have egress": {
			policyField:       "Missing Egress Network Policy",
			networkPolicyFile: "allow-egress-netpol.yaml",
		},
	}

	for name, testCase := range testCases {
		var policyId string
		var deploymentObj appsV1.Deployment
		feat := features.New(name).
			WithLabel("type", "simple").
			Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				applyFile(t, ctx, cfg, "nginx.yaml", &deploymentObj)
				return ctx
			}).
			WithStep("Create Policy", 1, func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
				t.Logf("Running `Create Policy` step")
				var err error
				policyId, err = createPolicyIfMissing(getFakePolicyRequest(testCase.policyField))
				if err != nil {
					t.Fatal(err)
				}
				return ctx
			}).
			Assess("Violation is shown", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
				testutils.Retry(t, 3, 6*time.Second, func(retryT testutils.T) {
					alerts := getAlertsForPolicy(retryT, policyName)
					assert.Len(retryT, alerts, 1)
				})
				return ctx
			}).
			Assess("Violation disappears", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
				var netpol v12.NetworkPolicy
				applyFile(t, ctx, config, testCase.networkPolicyFile, &netpol)
				patchAndWaitForReplicas(t, ctx, config, &deploymentObj, 1)
				testutils.Retry(t, 3, 6*time.Second, func(retryT testutils.T) {
					alerts := getAlertsForPolicy(retryT, policyName)
					assert.Len(retryT, alerts, 0)
				})
				return ctx
			}).
			Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
				deletePolicy(t, policyId)
				return ctx
			}).
			Feature()
		global.Test(t, feat)
	}
}

func getFakePolicyRequest(fieldName string) *roxv1.PostPolicyRequest {
	return &roxv1.PostPolicyRequest{
		Policy: &storage.Policy{
			Id: "",
			LifecycleStages: []storage.LifecycleStage{
				storage.LifecycleStage_DEPLOY,
			},
			Name:               policyName,
			IsDefault:          false,
			CriteriaLocked:     false,
			MitreVectorsLocked: false,
			Severity:           storage.Severity_LOW_SEVERITY,
			Categories:         []string{"Anomalous Activity"},
			PolicySections: []*storage.PolicySection{
				{
					SectionName: "example",
					PolicyGroups: []*storage.PolicyGroup{
						{
							FieldName:       fieldName,
							BooleanOperator: 0,
							Negate:          false,
							Values: []*storage.PolicyValue{
								{
									Value: "true",
								},
							},
						},
					},
				},
			},
		},
		EnableStrictValidation: false,
	}
}
