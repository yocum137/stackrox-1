package tester

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	centralDebug "github.com/stackrox/rox/sensor/debugger/central"
	"gopkg.in/yaml.v3"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

type TestFile struct {
	TestCase []TestCase `yaml:"testCase"`
}

type TestCase struct {
	Name       string       `yaml:"name"`
	Config     Config       `yaml:"config"`
	InputFiles []InputFile  `yaml:"inputFiles"`
	Assertions []Assertions `yaml:"assertions"`
}

type InputFile struct {
	Kind string `yaml:"kind"`
	File string `yaml:"file"`
}

type Config struct {
	RunAllCombinations bool          `yaml:"runAllCombinations"`
	ResourceBackoff    time.Duration `yaml:"resourceBackoff"`
	ExtraWaitTimer     time.Duration `yaml:"extraWaitTimer"`
	Namespace          string        `yaml:"namespace"`
}

type Assertions struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Assertion   string `yaml:"assertion"`
	Value       string `yaml:"value"`
}

func runPermutation(files []InputFile, i int, cb func(r []InputFile)) {
	if i > len(files) {
		cb(files)
		return
	}
	runPermutation(files, i+1, cb)
	for j := i + 1; j < len(files); j++ {
		files[i], files[j] = files[j], files[i]
		runPermutation(files, i+1, cb)
		files[i], files[j] = files[j], files[i]
	}
}

func UnmarshalTestFile(testCaseFilePath string) (*TestFile, error) {
	content, err := os.ReadFile(testCaseFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading test case file")
	}
	var testCaseFile TestFile
	if err := yaml.Unmarshal(content, &testCaseFile); err != nil {
		return nil, errors.Wrap(err, "unmarshaling yaml content")
	} else {
		return &testCaseFile, nil
	}
}

func createTestNs(ctx context.Context, r *resources.Resources, name string) (*v1.Namespace, error) {
	nsObj := v1.Namespace{}
	nsObj.Name = name
	if err := r.Create(ctx, &nsObj); err != nil {
		return nil, err
	}
	return &nsObj, nil
}

func objByKind(kind string) k8s.Object {
	switch kind {
	case "Deployment":
		return &appsV1.Deployment{}
	case "Role":
		return &v12.Role{}
	case "Binding":
		return &v12.RoleBinding{}
	default:
		log.Fatalf("unrecognized resource kind %s\n", kind)
		return nil
	}
}

func applyFile(ctx context.Context, ns string, r *resources.Resources, file InputFile) error {
	d := os.DirFS("tools/local-sensor/cases/resources")
	obj := objByKind(file.Kind)
	if err := decoder.DecodeFile(
		d,
		file.File,
		obj,
		decoder.MutateNamespace(ns),
	); err != nil {
		return err
	}

	if err := r.Create(ctx, obj); err != nil {
		return err
	}

	return nil
}

func runAssertions(messages []*central.MsgFromSensor, assertions []Assertions) bool {
	allPassed := true
	for _, assertion := range assertions {
		lastMessage := getLastMessage(messages, filterKind(assertion.Kind), filterDeploymentName("nginx-deployment"))
		if lastMessage == nil {
			fmt.Printf("No messages found for resources of kind %s\n", assertion.Kind)
			continue
		}
		event := lastMessage.GetEvent()
		//fmt.Printf("LAST EVENT: %+v\n", event.GetDeployment())
		if passed, err := CheckFields(event, assertion.Kind, assertion.Assertion, assertion.Value); err != nil {
			log.Fatalf("error on fields dynamically: %s\n", err)
		} else {
			if passed {
				fmt.Printf("  [SUCCESS] %s\n", assertion.Description)
			} else {
				allPassed = false
				fmt.Printf("  [FAILED] %s\n", assertion.Description)
			}
		}
	}
	return allPassed
}

func RunTestCase(testCase TestCase, fakeCentral *centralDebug.FakeService) error {
	var result [][]InputFile
	runPermutation(testCase.InputFiles, 0, func(perm []InputFile) {
		newPerm := make([]InputFile, len(perm))
		copy(newPerm, perm)
		result = append(result, newPerm)
	})
	mainConfig := envconf.New().WithKubeconfigFile(conf.ResolveKubeConfigFile())
	// Create namespace for each feature
	r, err := resources.New(mainConfig.Client().RESTConfig())
	if err != nil {
		return err
	}

	for i, perm := range result {
		fmt.Printf(" Running test permutation (%d/%d): %v\n", i+1, len(result), perm)
		if err := runTest(r, perm, testCase, fakeCentral); err != nil {
			return err
		}
	}

	return nil
}

func runTest(r *resources.Resources, files []InputFile, testCase TestCase, fakeCentral *centralDebug.FakeService) error {
	var namespaceObj *v1.Namespace
	defer func() {
		if namespaceObj == nil {
			return
		}

		// delete namespace
		err := r.Delete(context.Background(), namespaceObj)
		if err != nil {
			fmt.Printf("failed to cleanup namespace: %s\n", namespaceObj.Name)
		}

		// wait for deletion to be finished
		if err := wait.For(conditions.New(r).ResourceDeleted(namespaceObj)); err != nil {
			fmt.Printf("failed to wait for namespace deletion")
		}
	}()
	// create namespace
	ctx := context.TODO()
	namespaceName := "sensor-integration"
	var err error
	namespaceObj, err = createTestNs(ctx, r, namespaceName)
	if err != nil {
		return err
	}

	// create all resources with backoff
	for _, resFile := range files {
		if err := applyFile(ctx, namespaceName, r, resFile); err != nil {
			return errors.Wrapf(err, "applying resource file %s", resFile)
		}
		time.Sleep(testCase.Config.ResourceBackoff)
	}

	//fmt.Printf("Waiting for %v before running assertions\n", testCase.Config.ExtraWaitTimer)
	time.Sleep(testCase.Config.ExtraWaitTimer)

	// get all central messages
	messages := fakeCentral.GetAllMessages()
	if allPassed := runAssertions(messages, testCase.Assertions); !allPassed {
		return errors.Errorf("some assertions for test %s failed", testCase.Name)
	}

	return nil
}
