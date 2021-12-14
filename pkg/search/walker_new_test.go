package search

import (
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
)

func testCase(t *testing.T, category v1.SearchCategory, prefix string, object interface{}) {
	//newOptions := Walk(category, prefix, object).Original()
	newVisitorOptions := WalkVisitor(category, prefix, object).Original()

	legacyOptions := walkLegacy(category, prefix, object).Original()

	//assert.Equal(t, newOptions, legacyOptions)
	assert.Equal(t, newVisitorOptions, legacyOptions)
}

func TestWalk(t *testing.T) {
	cases := []struct {
		category v1.SearchCategory
		prefix   string
		object   interface{}
	}{
		{category: v1.SearchCategory_DEPLOYMENTS, prefix: "deployment", object: (*storage.Deployment)(nil)},
		//{category: v1.SearchCategory_ALERTS, prefix: "alert", object: (*storage.Alert)(nil)},
		//{category: v1.SearchCategory_PODS, prefix: "pod", object: (*storage.Pod)(nil)},
		//{category: v1.SearchCategory_ROLES, prefix: "k8s_role", object: (*storage.K8SRole)(nil)},
		//{category: v1.SearchCategory_ROLEBINDINGS, prefix: "k8s_role_binding", object: (*storage.K8SRoleBinding)(nil)},
		//{category: v1.SearchCategory_SUBJECTS, prefix: "subject", object: (*storage.Subject)(nil)},
		//{category: v1.SearchCategory_CLUSTERS, prefix: "cluster", object: (*storage.Cluster)(nil)},
		//{category: v1.SearchCategory_NETWORK_ENTITY, prefix: "network_entity", object: (*storage.NetworkEntity)(nil)},
		//{category: v1.SearchCategory_VULNERABILITIES, prefix: "c_v_e", object: (*storage.CVE)(nil)},
		//{category: v1.SearchCategory_REPORT_CONFIGURATIONS, prefix: "report_configuration", object: (*storage.ReportConfiguration)(nil)},
		//{category: v1.SearchCategory_VULN_REQUEST, prefix: "vulnerability_request", object: (*storage.VulnerabilityRequest)(nil)},
		//{category: v1.SearchCategory_NODE_COMPONENT_EDGE, prefix: "nodecomponentedge", object: (*storage.NodeComponentEdge)(nil)},
		//{category: v1.SearchCategory_NAMESPACES, prefix: "namespace_metadata", object: (*storage.NamespaceMetadata)(nil)},
		//{category: v1.SearchCategory_COMPLIANCE_STANDARD, prefix: "standard", object: (*v1.ComplianceStandard)(nil)},
		//{category: v1.SearchCategory_COMPLIANCE_CONTROL, prefix: "control", object: (*v1.ComplianceControl)(nil)},
		//{category: v1.SearchCategory_VULNERABILITIES, prefix: "image.scan.components.vulns", object: (*storage.EmbeddedVulnerability)(nil)},
		//{category: v1.SearchCategory_IMAGE_COMPONENTS, prefix: "image.scan.components", object: (*storage.EmbeddedImageScanComponent)(nil)},
		//{category: v1.SearchCategory_CLUSTER_VULN_EDGE, prefix: "cluster_c_v_e_edge", object: (*storage.ClusterCVEEdge)(nil)},
		//{category: v1.SearchCategory_IMAGE_COMPONENTS, prefix: "image_component", object: (*storage.ImageComponent)(nil)},
		//{category: v1.SearchCategory_SERVICE_ACCOUNTS, prefix: "service_account", object: (*storage.ServiceAccount)(nil)},
		//{category: v1.SearchCategory_IMAGE_COMPONENT_EDGE, prefix: "imagecomponentedge", object: (*storage.ImageComponentEdge)(nil)},
		//{category: v1.SearchCategory_PROCESS_BASELINES, prefix: "process_baseline", object: (*storage.ProcessBaseline)(nil)},
		//{category: v1.SearchCategory_SECRETS, prefix: "secret", object: (*storage.Secret)(nil)},
		//{category: v1.SearchCategory_NODES, prefix: "node", object: (*storage.Node)(nil)},
		//{category: v1.SearchCategory_VULNERABILITIES, prefix: "node.scan.components.vulns", object: (*storage.EmbeddedVulnerability)(nil)},
		//{category: v1.SearchCategory_IMAGE_COMPONENTS, prefix: "node.scan.components", object: (*storage.EmbeddedNodeScanComponent)(nil)},
		//{category: v1.SearchCategory_IMAGE_VULN_EDGE, prefix: "image_c_v_e_edge", object: (*storage.ImageCVEEdge)(nil)},
		//{category: v1.SearchCategory_RISKS, prefix: "risk", object: (*storage.Risk)(nil)},
		//{category: v1.SearchCategory_POLICIES, prefix: "policy", object: (*storage.Policy)(nil)},
		//{category: v1.SearchCategory_COMPONENT_VULN_EDGE, prefix: "component_c_v_e_edge", object: (*storage.ComponentCVEEdge)(nil)},
		//{category: v1.SearchCategory_IMAGES, prefix: "image", object: (*storage.Image)(nil)},
		//{category: v1.SearchCategory_PROCESS_INDICATORS, prefix: "process_indicator", object: (*storage.ProcessIndicator)(nil)},
	}
	for _, c := range cases {
		t.Run(c.category.String(), func(t *testing.T) {
			testCase(t, c.category, c.prefix, c.object)
		})
	}
}
