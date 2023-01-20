package fixtures

import (
	"strconv"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/fixtures/fixtureconsts"
)

// GetProcessIndicator returns a mock ProcessIndicator.
func GetProcessIndicator() *storage.ProcessIndicator {
	return &storage.ProcessIndicator{
		Id:           "b3523d84-ac1a-4daa-a908-62d196c5a741",
		DeploymentId: GetDeployment().GetId(),
		Signal: &storage.ProcessSignal{
			ContainerId:  "containerid",
			Name:         "apt-get",
			Args:         "install nmap",
			ExecFilePath: "bin",
			LineageInfo: []*storage.ProcessSignal_LineageInfo{
				{
					ParentUid:          22,
					ParentExecFilePath: "/bin/bash",
				},
				{
					ParentUid:          28,
					ParentExecFilePath: "/bin/curl",
				},
			},
		},
	}
}

// GetScopedProcessIndicator returns a mock ProcessIndicator belonging to the input scope.
func GetScopedProcessIndicator(ID string, clusterID string, namespace string) *storage.ProcessIndicator {
	return &storage.ProcessIndicator{
		Id:        ID,
		ClusterId: clusterID,
		Namespace: namespace,
	}
}

// GetProcessIndicatorForProcessListeningOnPorts returns a mock ProcessIndicator which should be matched with 
// a listening endpoint
func GetProcessIndicatorForListeningEndpoint(jobID int, port int) *storage.ProcessIndicator {
	return &storage.ProcessIndicator{
		Id:           "b3523d84-ac1a-4daa-a908-62d196c5a741",
		PodId:		fixtureconsts.PodUID1,
		DeploymentId: GetDeployment().GetId(),
		Signal: &storage.ProcessSignal{
			ContainerId:  "containerid",
			Name:         "socat",
			Args:         "TCP-LISTEN:" + strconv.Itoa(port) + ",fork STDOUT",
			ExecFilePath: "socat" + strconv.Itoa(jobID),
			LineageInfo: []*storage.ProcessSignal_LineageInfo{
				{
					ParentUid:          22,
					ParentExecFilePath: "/bin/bash",
				},
			},
		},
	}
}
