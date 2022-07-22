package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/sensor/tests/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/e2e-framework/klient/k8s"
)

var (
	NginxDeployment  = resource.YamlTestFile{Kind: "Deployment", File: "nginx.yaml"}
	NginxRole        = resource.YamlTestFile{Kind: "Pod", File: "nginx-pod.yaml"}
	NginxRoleBinding = resource.YamlTestFile{Kind: "Service", File: "nginx-service.yaml"}
)

func GetLastMessageWithDeploymentName(messages []*central.MsgFromSensor, n string) *central.MsgFromSensor {
	var lastMessage *central.MsgFromSensor
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].GetEvent().GetDeployment().GetName() == n {
			lastMessage = messages[i]
			break
		}
	}
	return lastMessage
}

func assertLastDeploymentHasPermissionLevel(t *testing.T, messages []*central.MsgFromSensor, permissionLevel storage.PermissionLevel) {
	lastNginxDeploymentUpdate := GetLastMessageWithDeploymentName(messages, "nginx-deployment")
	require.NotNil(t, lastNginxDeploymentUpdate, "should have found a message for nginx-deployment")
	deployment := lastNginxDeploymentUpdate.GetEvent().GetDeployment()
	assert.Equal(
		t,
		deployment.ServiceAccountPermissionLevel,
		permissionLevel,
		fmt.Sprintf("permission level has to be %s", permissionLevel),
	)
}

type ServiceDependencySuite struct {
	testContext *resource.TestContext
	suite.Suite
}

func Test_RoleDependency(t *testing.T) {
	suite.Run(t, new(ServiceDependencySuite))
}

var _ suite.SetupAllSuite = &ServiceDependencySuite{}
var _ suite.TearDownTestSuite = &ServiceDependencySuite{}

func (s *ServiceDependencySuite) TearDownTest() {
	// Clear any messages received in fake central during the test run
	s.testContext.GetFakeCentral().ClearReceivedBuffer()
}

func (s *ServiceDependencySuite) SetupSuite() {
	if testContext, err := resource.NewContext(s.T()); err != nil {
		s.Fail("failed to setup test context: %s", err)
	} else {
		s.testContext = testContext
	}
}

func (s *ServiceDependencySuite) Test_PermutationTest() {
	s.testContext.RunWithResourcesPermutation(
		[]resource.YamlTestFile{
			NginxDeployment,
			NginxRole,
			NginxRoleBinding,
		}, "Role Dependency", func(t *testing.T, testC *resource.TestContext, _ map[string]k8s.Object) {
			// Test context already takes care of creating and destroying resources
			time.Sleep(2 * time.Second)
			assertLastDeploymentHasPermissionLevel(
				t,
				testC.GetFakeCentral().GetAllMessages(),
				storage.PermissionLevel_ELEVATED_IN_NAMESPACE,
			)
			testC.GetFakeCentral().ClearReceivedBuffer()
		},
	)
}
