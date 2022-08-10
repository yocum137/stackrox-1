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

	yamlVersionKey = "app.kubernetes.io/version"

	scannerVersionByGithub = "https://raw.githubusercontent.com/stackrox/stackrox/%s/SCANNER_VERSION"
	collectorVersionByGithub = "https://raw.githubusercontent.com/stackrox/stackrox/%s/COLLECTOR_VERSION"
)

type scannerVersionCheck struct { }
var _ diagBundleCheck = (*scannerVersionCheck)(nil)
func (vc scannerVersionCheck) Name() string { return "scanner versions match" }

type collectorVersionCheck struct { }
var _ diagBundleCheck = (*collectorVersionCheck)(nil)
func (vc collectorVersionCheck) Name() string { return "collector versions match" }

type centralVersionCheck struct { }
var _ diagBundleCheck = (*centralVersionCheck)(nil)
func (vc centralVersionCheck) Name() string { return "central versions match" }

func (vc scannerVersionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := OK
	msgs := make([]string, 0)

	expectedVersions, err := parseVersionsFile(extractedBundlePath)
	if err != nil {
		return checkError(err)
	}

	// Check that scanner version correspond to the GitHub one for this commit.
	scannerGithubURL := fmt.Sprintf(scannerVersionByGithub, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", scannerGithubURL)
	scannerGithubVersion, err := getXVersionFromGithub(cliEnvironment, scannerGithubURL, cachedHttpClient)
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
	collectorGithubURL := fmt.Sprintf(collectorVersionByGithub, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", collectorGithubURL)
	collectorGithubVersion, err := getXVersionFromGithub(cliEnvironment, collectorGithubURL, cachedHttpClient)
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

	// Extract Central's actual version.
	centralActualVersion := ""
	c, err := os.ReadFile(filepath.Join(extractedBundlePath, centralFile))
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot read %q file", centralFile))
	}
	for _, line := range strings.Split(string(c), "\n") {
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, ":", 2)
		if split[0] == yamlVersionKey {
			centralActualVersion = strings.TrimSpace(split[1])
			break
		}
	}

	cliEnvironment.Logger().InfofLn("Extracted Central version: %q", centralActualVersion)
	if centralActualVersion == "" {
		return checkError(errors.Errorf("could not determine Central actual version: key %q not found in %q", yamlVersionKey, centralFile))
	}

	// Check that extracted version corresponds to the declared one.
	if centralActualVersion != expectedVersions.MainVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Central actual version: %s, Central version (%q): %s", centralActualVersion, versionsFile, expectedVersions.MainVersion))
	}

	return s, msgs, nil
}

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

func getXVersionFromGithub(cliEnvironment environment.Environment, url string, client *http.Client) (string, error) {
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
