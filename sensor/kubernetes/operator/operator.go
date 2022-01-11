// Package operator contains "operational logic" that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.
package operator

import (
	"context"
	"math/rand"
	"time"

	"github.com/cloudflare/cfssl/helpers"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	grpcUtil "github.com/stackrox/rox/pkg/grpc/util"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/sensor/common/clusterid"
	"github.com/stackrox/rox/sensor/kubernetes/sensor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	issueCertificatesTimeout            = 2 * time.Minute
	fetchSecretsTimeout                 = 2 * time.Minute
	updateSecretsTimeout                = 2 * time.Minute
	refreshSecretsMaxNumAttempts        = uint(5)
	refreshSecretAttemptWaitTime        = 5 * time.Minute
	refreshSecretAllAttemptsFailedWaitTime = 2 * time.Hour
	localScannerCredentialsSecretName   = "scanner-local-tls"
	localScannerDBCredentialsSecretName = "scanner-db-local-tls"
)

var (
	log             = logging.LoggerForModule()
	grpcCallOptions = []grpc.CallOption{grpc.UseCompressor(gzip.Name)}
)

type operatorImpl struct {
	secretsClient                        corev1.SecretInterface
	centralConnection                    *grpcUtil.LazyClientConn
	ctx                                  context.Context
	localScannerServiceClient            central.LocalScannerServiceClient
	numLocalScannerSecretRefreshAttempts uint
	refreshTimer                         *time.Timer
}

// Operator performs some operational logic that provides Sensor with some self-operation capabilities
// irrespective of the way it was deployed.
type Operator interface {
	Start(ctx context.Context) error
	Stop() error
}

// New creates a new operator
func New(k8sClient kubernetes.Interface, centralConnection *grpcUtil.LazyClientConn) Operator {
	return &operatorImpl{
		secretsClient: k8sClient.CoreV1().Secrets(sensor.GetSensorNamespace()),
		// Operations on this connection will block until centralConnection.Set is
		// called during Sensor.Start().
		centralConnection: centralConnection,
	}
}

// Start launches the processes that implement the "operational logic" of Sensor.
func (o *operatorImpl) Start(ctx context.Context) error {
	log.Info("Starting embedded operator.")

	o.ctx = ctx
	// FIXME: put all LocalScannerSecrets functionality in its own file and struct
	o.localScannerServiceClient = central.NewLocalScannerServiceClient(o.centralConnection)
	if err := o.scheduleLocalScannerSecretsRefresh(); err != nil {
		return errors.Wrapf(err, "failure scheduling local scanner secrets refresh")
	}

	log.Info("Embedded operator started.")

	return nil
}

func (o *operatorImpl) Stop() error {
	if o.refreshTimer != nil {
		o.refreshTimer.Stop()
	}
	return nil
}

func (o *operatorImpl) scheduleLocalScannerSecretsRefresh() error {
	localScannerCredsSecret, localScannerDBCredsSecret, fetchErr := o.fetchLocalScannerSecrets()
	if k8sErrors.IsNotFound(fetchErr) {
		log.Warnf("Some local scanner secret is missing, "+
			"operator will not maintain any local scanner secret fresh : %v", fetchErr)
		return nil
	}
	if fetchErr != nil {
		// FIXME wrap
		return fetchErr
	}

	// If certificates are already expired this refreshes immediately.
	o.doScheduleLocalScannerSecretsRefresh(getScannerSecretsDuration(localScannerCredsSecret, localScannerDBCredsSecret))
	return nil
}

func (o *operatorImpl) doScheduleLocalScannerSecretsRefresh(timeToRefresh time.Duration) {
	o.refreshTimer = time.AfterFunc(timeToRefresh, func() {
		nextTimeToRefresh, err := o.refreshLocalScannerSecrets()
		if err == nil {
			log.Infof("Successfully refreshed local Scanner credential secrets %s and %s, " +
				"will refresh again in %s",
				localScannerCredentialsSecretName, localScannerDBCredentialsSecretName, nextTimeToRefresh)
			o.numLocalScannerSecretRefreshAttempts = 0
			o.doScheduleLocalScannerSecretsRefresh(nextTimeToRefresh)
		} else {
			log.Errorf("Attempt %d to refresh local Scanner credential secrets, will retry in %s",
				o.numLocalScannerSecretRefreshAttempts, refreshSecretAttemptWaitTime)
			o.numLocalScannerSecretRefreshAttempts++
			if o.numLocalScannerSecretRefreshAttempts < refreshSecretsMaxNumAttempts {
				o.doScheduleLocalScannerSecretsRefresh(refreshSecretAttemptWaitTime)
			} else {
				log.Errorf("Failed to refresh local Scanner credential secrets after %d attempts, " +
					"will wait %s and restart the retry cycle",
					refreshSecretsMaxNumAttempts, refreshSecretAllAttemptsFailedWaitTime)
				o.numLocalScannerSecretRefreshAttempts = 0
				o.doScheduleLocalScannerSecretsRefresh(refreshSecretAllAttemptsFailedWaitTime)
			}
		}
	})
}

func getScannerSecretsDuration(localScannerCredsSecret, localScannerDBCredsSecret *v1.Secret) time.Duration {
	scannerDuration := getScannerSecretDuration(localScannerCredsSecret)
	scannerDBDuration := getScannerSecretDuration(localScannerDBCredsSecret)
	if scannerDuration > scannerDBDuration {
		return scannerDBDuration
	}
	return scannerDuration
}

func getScannerSecretDuration(scannerSecret *v1.Secret) time.Duration {
	scannerCertsData := scannerSecret.Data
	scannerCertBytes := scannerCertsData[mtls.ServiceCertFileName]
	scannerCert, err := helpers.ParseCertificatePEM(scannerCertBytes)
	if err != nil {
		// Note this also covers a secret with no certificates stored, which should be refreshed immediately.
		return 0
	}

	certValidityDurationSecs := scannerCert.NotAfter.Sub(scannerCert.NotBefore).Seconds()
	durationBeforeRenewalAttempt :=
		time.Duration(certValidityDurationSecs/2) - time.Duration(rand.Intn(int(certValidityDurationSecs/10)))
	certRenewalTime := scannerCert.NotBefore.Add(durationBeforeRenewalAttempt)
	timeToRefresh := time.Until(certRenewalTime)
	if timeToRefresh.Seconds() <= 0 {
		// Certificate is already expired.
		return 0
	}
	return timeToRefresh
}

func (o *operatorImpl) issueScannerCertificates() (*central.IssueLocalScannerCertsResponse, error) {
	// Blocks until Hello protocol is completed, when CentralCommunication sets the cluster id
	// after receiving it from centralHello.
	clusterid.Get()
	// We only support local Scanner running on the same namespace as Sensor.
	localScannerNamespace := sensor.GetSensorNamespace()

	ctx, cancel := context.WithTimeout(o.ctx, issueCertificatesTimeout)
	defer cancel()
	request := central.IssueLocalScannerCertsRequest{
		Namespace: localScannerNamespace,
	}

	return o.localScannerServiceClient.IssueLocalScannerCerts(ctx, &request, grpcCallOptions...)
}

func (o *operatorImpl) fetchLocalScannerSecrets() (*v1.Secret, *v1.Secret, error) {
	ctx, cancel := context.WithTimeout(o.ctx, fetchSecretsTimeout)
	defer cancel()

	// FIXME multierror
	localScannerCredsSecret, err := o.secretsClient.Get(ctx, localScannerCredentialsSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "for secret %s", localScannerCredentialsSecretName)
	}
	localScannerDBCredsSecret, err := o.secretsClient.Get(ctx, localScannerDBCredentialsSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "for secret %s", localScannerDBCredentialsSecretName)
	}

	return localScannerCredsSecret, localScannerDBCredsSecret, nil
}

func updateLocalScannerSecret(secret *v1.Secret, certificates *central.LocalScannerCertificates) {
	secret.Data = map[string][]byte{
		mtls.ServiceCertFileName: certificates.Cert,
		mtls.CACertFileName:      certificates.Ca,
		mtls.ServiceKeyFileName:  certificates.Key,
	}
}

// When any of the secrets is missing this returns and err such that k8sErrors.IsNotFound(err) is true
// On success it returns the duration after which the secrets should be refreshed
func (o *operatorImpl) refreshLocalScannerSecrets() (time.Duration, error) {
	// TODO: get both secrets just in case, only update that required by checking getNextSecretRefresh is 0
	localScannerCredsSecret, localScannerDBCredsSecret, err := o.fetchLocalScannerSecrets()
	if err != nil {
		// FIXME wrap
		return 0, err
	}
	certificates, err := o.issueScannerCertificates()
	if err != nil {
		// FIXME wrap
		return 0, err
	}
	updateLocalScannerSecret(localScannerCredsSecret, certificates.ScannerCerts)
	updateLocalScannerSecret(localScannerDBCredsSecret, certificates.ScannerDbCerts)

	ctx, cancel := context.WithTimeout(o.ctx, updateSecretsTimeout)
	defer cancel()
	// FIXME do a loop, and apply pattern elsewhere
	localScannerCredsSecret, err = o.secretsClient.Update(ctx, localScannerCredsSecret, metav1.UpdateOptions{})
	if err != nil {
		// FIXME wrap
		return 0, err
	}
	localScannerDBCredsSecret, err = o.secretsClient.Update(ctx, localScannerDBCredsSecret, metav1.UpdateOptions{})
	if err != nil {
		// FIXME wrap
		return 0, err
	}

	timeToRefresh := getScannerSecretsDuration(localScannerCredsSecret, localScannerDBCredsSecret)
	return timeToRefresh, nil
}
