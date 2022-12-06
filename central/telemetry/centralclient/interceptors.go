package centralclient

import (
	"strings"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/telemetry/phonehome"
)

// Adds Path, Code and User-Agent properties to API Call events for the API
// paths which start from the prefixes specified in the
// rhacs.redhat.com/telemetry-apipaths central deployment annotation
// ("*" value enables all paths) and are not in the ignoredPath list.
func apiCall(rp *phonehome.RequestParams, props map[string]any) bool {
	for _, ip := range ignoredPaths {
		if strings.HasPrefix(rp.Path, ip) {
			return false
		}
	}
	if trackedPaths.Contains("*") || trackedPaths.Contains(rp.Path) {
		props["Path"] = rp.Path
		props["Code"] = rp.Code
		props["User-Agent"] = rp.UserAgent
		return true
	}
	return false
}

// Adds Post Cluster call specific properties to the Post Cluster event.
func postCluster(rp *phonehome.RequestParams, props map[string]any) bool {
	if rp.Path != "/v1.ClustersService/PostCluster" {
		return false
	}
	props["Code"] = rp.Code
	if req, ok := rp.GrpcReq.(*storage.Cluster); ok {
		props["Cluster Type"] = req.GetType()
		props["Cluster ID"] = req.GetId()
	}
	return true
}

// Adds properties to the roxctl event.
func roxctl(rp *phonehome.RequestParams, props map[string]any) bool {
	if !strings.Contains(rp.UserAgent, "roxctl") {
		return false
	}
	props["Path"] = rp.Path
	props["Code"] = rp.Code
	props["User-Agent"] = rp.UserAgent
	if rp.HttpReq != nil {
		props["Protocol"] = "HTTP"
	} else {
		props["Protocol"] = "GRPC"
	}
	return true
}
