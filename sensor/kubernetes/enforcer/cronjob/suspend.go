package cronjob

import (
	"context"
	"fmt"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/branding"
	kubernetesPkg "github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/retry"
	"github.com/stackrox/rox/sensor/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	batchV1      = "batch/v1"
	batchV1beta1 = "batch/v1beta1"
)

// Suspend suspends the cron job
func Suspend(ctx context.Context, client kubernetes.Interface, deploymentInfo *central.DeploymentEnforcement) (err error) {
	depName := deploymentInfo.GetDeploymentName()
	depType := deploymentInfo.GetDeploymentType()
	forcePatch := true

	if ok, apiErr := utils.HasAPI(client, batchV1, kubernetesPkg.CronJob); ok && apiErr == nil {
		patch := fmt.Sprintf("{\"metadata\": {\"name\": %q}, \"kind\": %q, \"apiVersion\": %s, \"spec\": {\"suspend\": true}}", depName, depType, batchV1)

		_, err = client.BatchV1().CronJobs(deploymentInfo.GetNamespace()).Patch(ctx, depName, types.ApplyPatchType,
			[]byte(patch),
			metav1.PatchOptions{
				TypeMeta: metav1.TypeMeta{
					Kind:       depType,
					APIVersion: batchV1,
				},
				FieldManager: branding.GetProductName(),
				Force:        &forcePatch,
			})
		if err != nil {
			return retry.MakeRetryable(err)
		}
	} else {
		patch := fmt.Sprintf("{\"metadata\": {\"name\": %s}, \"kind\": %s, \"apiVersion\": %s, \"spec\": {\"suspend\": true}}", depName, depType, batchV1beta1)

		_, err = client.BatchV1beta1().CronJobs(deploymentInfo.GetNamespace()).Patch(ctx, depName, types.ApplyPatchType,
			[]byte(patch),
			metav1.PatchOptions{
				TypeMeta: metav1.TypeMeta{
					Kind:       depType,
					APIVersion: batchV1beta1,
				},
				FieldManager: branding.GetProductName(),
				Force:        &forcePatch,
			})
		if err != nil {
			return retry.MakeRetryable(err)
		}
	}
	return nil
}
