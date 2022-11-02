package replay

import (
	"testing"
	"time"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	centralDebug "github.com/stackrox/rox/sensor/debugger/central"
	"github.com/stackrox/rox/sensor/debugger/k8s"
	"github.com/stackrox/rox/sensor/debugger/message"
	"github.com/stackrox/rox/sensor/testutils"
)

var (
	testCase replayTestCase

	setupTestCase sync.Once
)

func Benchmark_ReplayEvents(b *testing.B) {
	b.StopTimer()

	setupTestCase.Do(func() {
		testCase.fakeClient = k8s.MakeFakeClient()
		envIsolator := envisolator.NewEnvIsolator(b)
		envIsolator.Setenv("ROX_MTLS_CERT_FILE", "../../../tools/local-sensor/certs/cert.pem")
		envIsolator.Setenv("ROX_MTLS_KEY_FILE", "../../../tools/local-sensor/certs/key.pem")
		envIsolator.Setenv("ROX_MTLS_CA_FILE", "../../../tools/local-sensor/certs/caCert.pem")
		envIsolator.Setenv("ROX_MTLS_CA_KEY_FILE", "../../../tools/local-sensor/certs/caKey.pem")

		policies, err := testutils.GetPoliciesFromFile("data/policies.json")
		if err != nil {
			panic(err)
		}
		testCase.fakeCentral = centralDebug.MakeFakeCentralWithInitialMessages(
			message.SensorHello("00000000-0000-4000-A000-000000000000"),
			message.ClusterConfig(),
			message.PolicySync(policies),
			message.BaselineSync([]*storage.ProcessBaseline{}))

		testCase.resyncTime = 1 * time.Second
		testCase.ackChannel = make(chan *central.SensorEvent)
		testCase.writer = setupSensor(testCase.fakeCentral, testCase.fakeClient, testCase.resyncTime, testCase.ackChannel)
		//defer writer.close()
	})

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		testCase.fakeCentral.ClearReceivedBuffer()
		runTestCase(b, replayTestCase{
			k8sEventsFile:    "data/safety-net-alerts-k8s-trace.jsonl",
			sensorOutputFile: "data/safety-net-alerts-central-out.bin",
			ackChannel:       testCase.ackChannel,
			resyncTime:       testCase.resyncTime,
			fakeCentral:      testCase.fakeCentral,
			fakeClient:       testCase.fakeClient,
			writer:           testCase.writer,
		})
		time.Sleep(2 * testCase.resyncTime)
	}
}
