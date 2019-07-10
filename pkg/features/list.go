package features

var (
	// K8sRBAC is used to enable k8s rbac collection and processing
	// NB: When removing this feature flag, please also remove references to it in .circleci/config.yml
	K8sRBAC = registerFeature("Enable k8s RBAC objects collection and processing", "ROX_K8S_RBAC", true)

	// ScopedAccessControl controls whether scoped access control is enabled.
	// NB: When removing this feature flag, please also remove references to it in .circleci/config.yml
	ScopedAccessControl = registerFeature("Scoped Access Control", "ROX_SCOPED_ACCESS_CONTROL", true)

	// ClientCAAuth will enable authenticating to central via client certificate authentication
	ClientCAAuth = registerFeature("Client Certificate Authentication", "ROX_CLIENT_CA_AUTH", true)

	// PlaintextExposure enables specifying a port under which to expose the UI and API without TLS.
	PlaintextExposure = registerFeature("Plaintext UI & API exposure", "ROX_PLAINTEXT_EXPOSURE", true)
)
