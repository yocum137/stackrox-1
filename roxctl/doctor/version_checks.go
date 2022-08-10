package doctor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/version"
	"github.com/stackrox/rox/roxctl/common/environment"
)

const (
	versionsFile = "versions.json"
	centralFile = "kubernetes/_central-cluster/stackrox/central/deployment-central.yaml"
	admissionControllerFile = "kubernetes/_central-cluster/stackrox/admission-control/deployment-admission-control.yaml"

	sensorFileP = "kubernetes/%s/stackrox/sensor/deployment-sensor.yaml"
	collectorFileP = "kubernetes/%s/stackrox/collector/collector-sensor.yaml"

	yamlVersionKey = "app.kubernetes.io/version"

	scannerVersionByGithubP   = "https://raw.githubusercontent.com/stackrox/stackrox/%s/SCANNER_VERSION"
	collectorVersionByGithubP = "https://raw.githubusercontent.com/stackrox/stackrox/%s/COLLECTOR_VERSION"
)

type scannerVersionCheck struct { }
var _ diagBundleCheck = (*scannerVersionCheck)(nil)
func (vc scannerVersionCheck) Name() string { return "scanner version" }

type collectorVersionCheck struct { }
var _ diagBundleCheck = (*collectorVersionCheck)(nil)
func (vc collectorVersionCheck) Name() string { return "collector version" }

type centralVersionCheck struct { }
var _ diagBundleCheck = (*centralVersionCheck)(nil)
func (vc centralVersionCheck) Name() string { return "central services versions" }

type sensorVersionCheck struct { }
var _ diagBundleCheck = (*sensorVersionCheck)(nil)
func (vc sensorVersionCheck) Name() string { return "secured cluster services versions" }

func (vc scannerVersionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := OK
	msgs := make([]string, 0)

	expectedVersions, err := parseVersionsFile(extractedBundlePath)
	if err != nil {
		return checkError(err)
	}

	// Check that scanner version correspond to the GitHub one for this commit.
	scannerGithubURL := fmt.Sprintf(scannerVersionByGithubP, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", scannerGithubURL)
	scannerGithubVersion, err := getXVersionFromGithub(scannerGithubURL, cachedHttpClient)
	if err != nil {
		return checkError(err)
	}
	if scannerGithubVersion != expectedVersions.ScannerVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Scanner version (GitHub): %s, Scanner version (%q): %s", scannerGithubVersion, versionsFile, expectedVersions.ScannerVersion))
	}

	return s, msgs, nil
}

func (vc collectorVersionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := OK
	msgs := make([]string, 0)

	expectedVersions, err := parseVersionsFile(extractedBundlePath)
	if err != nil {
		return checkError(err)
	}

	// Check that collector version correspond to the GitHub one for this commit.
	collectorGithubURL := fmt.Sprintf(collectorVersionByGithubP, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", collectorGithubURL)
	collectorGithubVersion, err := getXVersionFromGithub(collectorGithubURL, cachedHttpClient)
	if err != nil {
		return checkError(err)
	}

	if collectorGithubVersion != expectedVersions.CollectorVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Collector version (GitHub): %s, Collector version (%q): %s", collectorGithubVersion, versionsFile, expectedVersions.CollectorVersion))
	}

	return s, msgs, nil
}

func (vc centralVersionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := OK
	msgs := make([]string, 0)

	expectedVersions, err := parseVersionsFile(extractedBundlePath)
	if err != nil {
		return checkError(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// Extract Central actual version.
	c, err := os.ReadFile(filepath.Join(extractedBundlePath, centralFile))
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot read %q file", centralFile))
	}

	centralActualVersion := getYamlValueByKey(string(c), yamlVersionKey)
	cliEnvironment.Logger().InfofLn("Extracted Central version: %q", centralActualVersion)
	if centralActualVersion == "" {
		return checkError(errors.Errorf("could not determine Central actual version: key %q not found in %q", yamlVersionKey, centralFile))
	}

	// Check that extracted version corresponds to the declared one.
	if centralActualVersion != expectedVersions.MainVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Central actual version: %s, Central version (%q): %s", centralActualVersion, versionsFile, expectedVersions.MainVersion))
	}

	////////////////////////////////////////////////////////////////////////////
	// Extract Admission Controller actual version.
	c, err = os.ReadFile(filepath.Join(extractedBundlePath, admissionControllerFile))
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot read %q file", admissionControllerFile))
	}

	admissionControllerActualVersion := getYamlValueByKey(string(c), yamlVersionKey)
	cliEnvironment.Logger().InfofLn("Extracted Admission Controller version: %q", admissionControllerActualVersion)
	if admissionControllerActualVersion == "" {
		return checkError(errors.Errorf("could not determine Admission Controller actual version: key %q not found in %q", yamlVersionKey, centralFile))
	}

	// Check that extracted version corresponds to the declared one.
	if admissionControllerActualVersion != expectedVersions.MainVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Admission Controller actual version: %s, Admission Controller version (%q): %s", admissionControllerActualVersion, versionsFile, expectedVersions.MainVersion))
	}

	return s, msgs, nil
}

func (vc sensorVersionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := OK
	msgs := make([]string, 0)

	expectedVersions, err := parseVersionsFile(extractedBundlePath)
	if err != nil {
		return checkError(err)
	}

	// Extract Central actual version. We need this to check if it is smaller
	// than Sensor version.
	c, err := os.ReadFile(filepath.Join(extractedBundlePath, centralFile))
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot read %q file", centralFile))
	}

	centralActualVersion := getYamlValueByKey(string(c), yamlVersionKey)
	cliEnvironment.Logger().InfofLn("Extracted Central version: %q", centralActualVersion)
	if centralActualVersion == "" {
		return checkError(errors.Errorf("could not determine Central actual version: key %q not found in %q", yamlVersionKey, centralFile))
	}

	clusterNames, _ := getAllSecuredClusters(extractedBundlePath)

	// Iterate all secured clusters and check every sensor and collector.
	for _, cluster := range clusterNames {
		////////////////////////////////////////////////////////////////////////
		// Sensor.
		sensorFile := fmt.Sprintf(sensorFileP, cluster)
		c, err := os.ReadFile(filepath.Join(extractedBundlePath, sensorFile))
		if err != nil {
			return checkError(errors.Wrapf(err, "cannot read %q file", sensorFile))
		}

		sensorActualVersion := getYamlValueByKey(string(c), yamlVersionKey)
		cliEnvironment.Logger().InfofLn("Extracted %q's Sensor version: %q", cluster, sensorActualVersion)
		if sensorActualVersion == "" {
			return checkError(errors.Errorf("could not determine %q's Sensor actual version: key %q not found in %q", cluster, yamlVersionKey, sensorFile))
		}

		// Check that extracted version corresponds to the declared one.
		if sensorActualVersion != expectedVersions.MainVersion {
			s.AtLeast(Warning)
			msgs = append(msgs, fmt.Sprintf("%q's Sensor actual version: %s, Sensor version (%q): %s", cluster, sensorActualVersion, versionsFile, expectedVersions.MainVersion))
		}

		// Check that Sensor version is not greater than Central version.
		if sensorActualVersion > centralActualVersion {
			s.AtLeast(Problem)
			msgs = append(msgs, fmt.Sprintf("%q's Sensor actual version %s is ahead of Central version %s; this can trigger incorrect or unexpected behaviour", cluster, sensorActualVersion, centralActualVersion))
		}

		////////////////////////////////////////////////////////////////////////
		// Collector.
		// TODO(alexr)
	}

	return s, msgs, nil
}

////////////////////////////////////////////////////////////////////////////////
// Helpers

func parseVersionsFile(extractedBundlePath string) (vs version.Versions, err error) {
	content, err := os.ReadFile(filepath.Join(extractedBundlePath, versionsFile))
	if err != nil {
		err = errors.Wrapf(err, "cannot read %q file", versionsFile)
		return
	}

	err = json.Unmarshal(content, &vs)
	if err != nil {
		err = errors.Wrapf(err, "cannot parse %q into expected JSON", versionsFile)
	}

	return
}

func getXVersionFromGithub(url string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating a request to GitHub")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "querying %q", url)
	}
	defer utils.IgnoreError(resp.Body.Close)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("received %s: %s", resp.Status, body)
	}

	return strings.TrimSpace(string(body)), nil
}

func getYamlValueByKey(content string, yamlKey string) string {
	var value string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, ":", 2)
		if split[0] == yamlVersionKey {
			value = strings.TrimSpace(split[1])
			break
		}
	}
	return value
}

func getAllSecuredClusters(extractedBundlePath string) ([]string, error) {
	dir := filepath.Join(extractedBundlePath, "kubernetes")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "iterate %q directory", dir)
	}

	clusterNames := make([]string, 0)
	for _, d := range files {
		if !d.IsDir() {
			continue
		}
		if strings.Contains(d.Name(), "central-cluster") {
			continue
		}
		clusterNames = append(clusterNames, d.Name())
	}

	return clusterNames, nil
}
