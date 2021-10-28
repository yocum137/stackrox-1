// Package resources lists all resource types used by Central.
package resources

import (
	"sort"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/features"
)

// All resource types that we want to define (for the purposes of enforcing
// API permissions) must be defined here.
// KEEP THE FOLLOWING LIST SORTED IN LEXICOGRAPHIC ORDER.
var (
	Alert                            = newResourceMetadata("Alert", permissions.NamespaceScope)
	Cluster                          = newResourceMetadata("Cluster", permissions.ClusterScope)
	CVE                              = newResourceMetadata("CVE", permissions.NamespaceScope)
	Deployment                       = newResourceMetadata("Deployment", permissions.NamespaceScope)
	Image                            = newResourceMetadata("Image", permissions.NamespaceScope)
	ImageComponent                   = newResourceMetadata("ImageComponent", permissions.NamespaceScope)
	Indicator                        = newResourceMetadata("Indicator", permissions.NamespaceScope)
	K8sRole                          = newResourceMetadata("K8sRole", permissions.NamespaceScope)
	K8sRoleBinding                   = newResourceMetadata("K8sRoleBinding", permissions.NamespaceScope)
	K8sSubject                       = newResourceMetadata("K8sSubject", permissions.NamespaceScope)
	Namespace                        = newResourceMetadata("Namespace", permissions.NamespaceScope)
	NetworkGraph                     = newResourceMetadata("NetworkGraph", permissions.NamespaceScope)
	NetworkPolicy                    = newResourceMetadata("NetworkPolicy", permissions.NamespaceScope)
	Node                             = newResourceMetadata("Node", permissions.ClusterScope)
	Secret                           = newResourceMetadata("Secret", permissions.NamespaceScope)
	ServiceAccount                   = newResourceMetadata("ServiceAccount", permissions.NamespaceScope)
	VulnerabilityManagementRequests  = newResourceMetadataWithFeatureFlag("VulnerabilityManagementRequests", permissions.GlobalScope, features.VulnRiskManagement)
	VulnerabilityManagementApprovals = newResourceMetadataWithFeatureFlag("VulnerabilityManagementApprovals", permissions.GlobalScope, features.VulnRiskManagement)
	VulnerabilityReports             = newResourceMetadataWithFeatureFlag("VulnerabilityReports", permissions.GlobalScope, features.VulnReporting)

	// Internal Resources
	ComplianceOperator = newInternalResourceMetadata("ComplianceOperator", permissions.GlobalScope)

	/* New resources */
	// AuthPlugin, AuthProvider, Group, Licenses, Role, User
	Access = newResourceMetadata("Access", permissions.GlobalScope)

	// AllComments, Config, DebugLogs, NetworkGraphConfig, ProbeUpload,
	// ScannerBundle, ScannerDefinitions, SensorUpgradeConfig, ServiceIdentity
	Administration = newResourceMetadata("Administration", permissions.GlobalScope)

	// Compliance, ComplianceRunSchedule, ComplianceRuns
	Compliance = newResourceMetadata("Compliance", permissions.ClusterScope)

	// This works in combination with scoped Deployment and replaces
	// NetworkBaseline, ProcessWhitelist, Risk, WatchedImage
	DeploymentExtension = newResourceMetadata("DeploymentExtension", permissions.GlobalScope)

	// APIToken, BackupPlugins, ImageIntegration, Notifier
	Integration = newResourceMetadata("Integration", permissions.GlobalScope)

	// Detection, Policy
	Policy                           = newResourceMetadata("Policy", permissions.GlobalScope)

	/* Deprecated */
	// in favour of Access
	AuthPlugin                       = newResourceMetadata("AuthPlugin", permissions.GlobalScope)
	AuthProvider                     = newResourceMetadata("AuthProvider", permissions.GlobalScope)
	Group                            = newResourceMetadata("Group", permissions.GlobalScope)
	Licenses                         = newResourceMetadata("Licenses", permissions.GlobalScope)
	Role                             = newResourceMetadata("Role", permissions.GlobalScope)
	User                             = newResourceMetadata("User", permissions.GlobalScope)
	// in favour of Integration
	APIToken                         = newResourceMetadata("APIToken", permissions.GlobalScope)
	BackupPlugins                    = newResourceMetadata("BackupPlugins", permissions.GlobalScope)
	ImageIntegration                 = newResourceMetadata("ImageIntegration", permissions.GlobalScope)
	Notifier                         = newResourceMetadata("Notifier", permissions.GlobalScope)
	// in favour of Compliance
	ComplianceRunSchedule            = newResourceMetadata("ComplianceRunSchedule", permissions.GlobalScope)
	ComplianceRuns                   = newResourceMetadata("ComplianceRuns", permissions.ClusterScope)
	// in favour of Administration
	AllComments                      = newResourceMetadata("AllComments", permissions.GlobalScope)
	Config                           = newResourceMetadata("Config", permissions.GlobalScope)
	DebugLogs                        = newResourceMetadata("DebugLogs", permissions.GlobalScope)
	NetworkGraphConfig               = newResourceMetadata("NetworkGraphConfig", permissions.GlobalScope)
	ProbeUpload                      = newResourceMetadata("ProbeUpload", permissions.GlobalScope)
	ScannerBundle                    = newResourceMetadata("ScannerBundle", permissions.GlobalScope)
	ScannerDefinitions               = newResourceMetadata("ScannerDefinitions", permissions.GlobalScope)
	SensorUpgradeConfig              = newResourceMetadata("SensorUpgradeConfig", permissions.GlobalScope)
	ServiceIdentity                  = newResourceMetadata("ServiceIdentity", permissions.GlobalScope)
	// in favour of Policy
	Detection                        = newResourceMetadata("Detection", permissions.GlobalScope)
	// in favour of Deployment + DeploymentExtensions
	NetworkBaseline                  = newResourceMetadata("NetworkBaseline", permissions.NamespaceScope)
	ProcessWhitelist                 = newResourceMetadata("ProcessWhitelist", permissions.NamespaceScope)
	Risk                             = newResourceMetadata("Risk", permissions.NamespaceScope)
	// in favour of DeploymentExtensions
	WatchedImage                     = newResourceMetadata("WatchedImage", permissions.GlobalScope)

	resourceToMetadata = make(map[permissions.Resource]permissions.ResourceMetadata)
)

func newResourceMetadata(name permissions.Resource, scope permissions.ResourceScope) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource: name,
		Scope:    scope,
	}
	resourceToMetadata[name] = md
	return md
}

func newResourceMetadataWithFeatureFlag(name permissions.Resource, scope permissions.ResourceScope, flag features.FeatureFlag) permissions.ResourceMetadata {
	md := permissions.ResourceMetadata{
		Resource: name,
		Scope:    scope,
	}
	if flag.Enabled() {
		resourceToMetadata[name] = md
	}
	return md
}

func newInternalResourceMetadata(name permissions.Resource, scope permissions.ResourceScope) permissions.ResourceMetadata {
	return permissions.ResourceMetadata{
		Resource: name,
		Scope:    scope,
	}
}

// ListAll returns a list of all resources.
func ListAll() []permissions.Resource {
	resources := make([]permissions.Resource, 0, len(resourceToMetadata))
	for _, metadata := range ListAllMetadata() {
		resources = append(resources, metadata.Resource)
	}
	return resources
}

// ListAllMetadata returns a list of all resource metadata.
func ListAllMetadata() []permissions.ResourceMetadata {
	metadatas := make([]permissions.ResourceMetadata, 0, len(resourceToMetadata))
	for _, metadata := range resourceToMetadata {
		metadatas = append(metadatas, metadata)
	}
	sort.SliceStable(metadatas, func(i, j int) bool {
		return string(metadatas[i].Resource) < string(metadatas[j].Resource)
	})
	return metadatas
}

// AllResourcesViewPermissions returns a slice containing view permissions for all resource types.
func AllResourcesViewPermissions() []permissions.ResourceWithAccess {
	metadatas := ListAllMetadata()
	result := make([]permissions.ResourceWithAccess, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = permissions.ResourceWithAccess{
			// We want to ensure access to *all* resources, so when using SAC, always perform legacy auth (= enforcement
			// at the global scope) even for cluster- or namespace-scoped resources.
			Resource: permissions.WithLegacyAuthForSAC(metadata, true),
			Access:   storage.Access_READ_ACCESS,
		}
	}
	return result
}

// AllResourcesModifyPermissions returns a slice containing write permissions for all resource types.
func AllResourcesModifyPermissions() []permissions.ResourceWithAccess {
	metadatas := ListAllMetadata()
	result := make([]permissions.ResourceWithAccess, len(metadatas))
	for i, metadata := range metadatas {
		result[i] = permissions.ResourceWithAccess{
			// We want to ensure access to *all* resources, so when using SAC, always perform legacy auth (= enforcement
			// at the global scope) even for cluster- or namespace-scoped resources.
			Resource: permissions.WithLegacyAuthForSAC(metadata, true),
			Access:   storage.Access_READ_WRITE_ACCESS,
		}
	}
	return result
}

// MetadataForResource returns the metadata for the given resource. If the resource is unknown, metadata for this
// resource with global scope is returned.
func MetadataForResource(res permissions.Resource) (permissions.ResourceMetadata, bool) {
	md, found := resourceToMetadata[res]
	if !found {
		md.Resource = res
		md.Scope = permissions.GlobalScope
	}
	return md, found
}
