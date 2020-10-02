import static Services.getPolicies
import static Services.waitForViolation

import groups.BAT
import objects.Deployment
import objects.SecretKeyRef
import objects.Volume
import orchestratormanager.OrchestratorTypes
import org.junit.Assume
import org.junit.experimental.categories.Category
import spock.lang.Shared
import spock.lang.Unroll
import services.CreatePolicyService
import io.stackrox.proto.api.v1.PolicyServiceOuterClass.PatchPolicyRequest
import util.Env

class BuiltinPoliciesTest extends BaseSpecification {
    static final private String TRIGGER_MOST = "trigger-most"
    static final private String TRIGGER_ALPINE = "trigger-alpine"
    static final private String TRIGGER_DOCKER_MOUNT = "trigger-docker-mount"
    static final private String TRIGGER_UNSCANNED = "trigger-unscanned"
    static final private String TEST_PASSWORD = "test-password"

    static final private List<Deployment> DEPLOYMENTS = [
            new Deployment()
                    .setName(TRIGGER_MOST)
                    .setImage("us.gcr.io/stackrox-ci/qa/trigger-policy-violations/most:0.19")
                    // For: "Emergency Deployment Annotation"
                    .addAnnotation("admission.stackrox.io/break-glass", "yay")
                    // For: "Secret Mounted as Environment Variable"
                    .addEnvValueFromSecretKeyRef(
                            "TEST_PASSWORD",
                            new SecretKeyRef(key: "password", name: TEST_PASSWORD)
                    )
                    // For: "Mounting Sensitive Host Directories"
                    .addVolume("sensitive", "/etc/", true)
                    // For: Iptables Executed in Privileged Container
                    .setPrivilegedFlag(true),
            // For: "Alpine Linux Package Manager (apk) in Image"
            new Deployment()
                    .setName(TRIGGER_ALPINE)
                    .setImage("us.gcr.io/stackrox-ci/qa/trigger-policy-violations/alpine:0.6"),
    ]
    static final private List<Deployment> NO_WAIT_DEPLOYMENTS = [
            new Deployment()
                    .setName(TRIGGER_DOCKER_MOUNT)
                    .setImage("nginx:latest")
                    .addVolume(new Volume(name: "docker-sock",
                            hostPath: "/var/run/docker.sock",
                            mountPath: "/var/run/docker.sock")),
            new Deployment()
                    .setName(TRIGGER_UNSCANNED)
                    .setImage("non-existent:image"),
    ]

    @Shared
    private List<String> disabledPolicyIds

    def setupSpec() {
        disabledPolicyIds = []
        getPolicies().forEach {
            policy ->
            if (policy.disabled) {
                println "Temporarily enabling a disabled policy for testing: ${policy.name}"
                CreatePolicyService.patchPolicy(
                        PatchPolicyRequest.newBuilder().setId(policy.id).setDisabled(false).build()
                )
                disabledPolicyIds.add(policy.id)
            }
        }

        orchestrator.createSecret(TEST_PASSWORD)

        for (Deployment deployment : NO_WAIT_DEPLOYMENTS) {
            println("Starting ${deployment.name} without waiting for deployment")
            orchestrator.createDeploymentNoWait(deployment)
        }

        orchestrator.batchCreateDeployments(DEPLOYMENTS)
        for (Deployment deployment : DEPLOYMENTS) {
            println("Waiting for ${deployment.name}")
            assert Services.waitForDeployment(deployment)
        }
    }

    def cleanupSpec() {
        disabledPolicyIds.forEach {
            id ->
            println "Re-disabling a policy after test"
            CreatePolicyService.patchPolicy(
                    PatchPolicyRequest.newBuilder().setId(id).setDisabled(true).build()
            )
        }

        for (Deployment deployment : DEPLOYMENTS + NO_WAIT_DEPLOYMENTS) {
            orchestrator.deleteDeployment(deployment)
        }

        orchestrator.deleteSecret(TEST_PASSWORD)
    }

    @Unroll
    @Category([BAT])
    def "Verify policy '#policyName' is triggered"(String policyName, String deploymentName) {
        // ROX-5298 - Policy tests are unreliable on Openshift
        Assume.assumeTrue(Env.mustGetOrchestratorType() != OrchestratorTypes.OPENSHIFT)

        when:
        "An existing policy"
        assert getPolicies().find { it.name == policyName }

        then:
        "Verify Violation for #policyName is triggered"
        assert waitForViolation(deploymentName, policyName, 30)

        where:
        "Data inputs are:"

        policyName                                                   | deploymentName
        // "30-Day Scan Age" <- Not covered
        "ADD Command used instead of COPY"                           | TRIGGER_MOST
        // "Alpine Linux Package Manager (apk) in Image"            | TRIGGER_ALPINE  // ROX-5099 does not trigger
        "Alpine Linux Package Manager Execution"                     | TRIGGER_ALPINE
        // "CAP_SYS_ADMIN capability added" <- Not covered
        "chkconfig Execution"                                        | TRIGGER_MOST
        "Container using read-write root filesystem"                 | TRIGGER_MOST
        "Compiler Tool Execution"                                    | TRIGGER_MOST
        "crontab Execution"                                          | TRIGGER_MOST
        "Cryptocurrency Mining Process Execution"                    | TRIGGER_MOST
        "Curl in Image"                                              | TRIGGER_MOST
        "Emergency Deployment Annotation"                            | TRIGGER_MOST
        "Fixable CVSS >= 6 and Privileged"                           | TRIGGER_MOST
        // "Heartbleed: CVE-2014-0160" <- Not covered
        "Images with no scans"                                       | TRIGGER_UNSCANNED
        // "Improper Usage of Orchestrator Secrets Volume"          | TRIGGER_MOST  // ROX-5098 does not trigger
        "Insecure specified in CMD"                                  | TRIGGER_MOST
        "iptables Execution"                                         | TRIGGER_MOST
        "Iptables Executed in Privileged Container"                  | TRIGGER_MOST
        "Linux Group Add Execution"                                  | TRIGGER_MOST
        "Linux User Add Execution"                                   | TRIGGER_MOST
        "Login Binaries"                                             | TRIGGER_MOST
        "Mount Docker Socket"                                        | TRIGGER_DOCKER_MOUNT
        "Mounting Sensitive Host Directories"                        | TRIGGER_MOST
        "Netcat Execution Detected"                                  | TRIGGER_MOST
        "Network Management Execution"                               | TRIGGER_MOST
        "nmap Execution"                                             | TRIGGER_MOST
        "No resource requests or limits specified"                   | TRIGGER_MOST
        "Password Binaries"                                          | TRIGGER_MOST
        "Process Targeting Cluster Kubelet Endpoint"                 | TRIGGER_MOST
        "Process Targeting Cluster Kubernetes Docker Stats Endpoint" | TRIGGER_MOST
        "Process Targeting Kubernetes Service Endpoint"              | TRIGGER_MOST
        "Process with UID 0"                                         | TRIGGER_MOST
        "Red Hat Package Manager Execution"                          | TRIGGER_MOST
        "Remote File Copy Binary Execution"                          | TRIGGER_MOST
        "Required Annotation: Email"                                 | TRIGGER_MOST
        "Required Annotation: Owner/Team"                            | TRIGGER_MOST
        "Required Image Label"                                       | TRIGGER_MOST
        "Required Label: Owner/Team"                                 | TRIGGER_MOST
        "Secret Mounted as Environment Variable"                     | TRIGGER_MOST
        "Secure Shell (ssh) Port Exposed in Image"                   | TRIGGER_MOST
        "Secure Shell Server (sshd) Execution"                       | TRIGGER_MOST
        "SetUID Processes"                                           | TRIGGER_MOST
        "Shadow File Modification"                                   | TRIGGER_MOST
        "Shell Spawned by Java Application"                          | TRIGGER_MOST
        "systemctl Execution"                                        | TRIGGER_MOST
        "systemd Execution"                                          | TRIGGER_MOST
        "Ubuntu Package Manager Execution"                           | TRIGGER_MOST
        "Wget in Image"                                              | TRIGGER_MOST
    }
}
