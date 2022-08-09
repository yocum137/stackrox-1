package doctor

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/version"
	"github.com/stackrox/rox/roxctl/common/environment"
)

type versionCheck struct { }

var _ diagBundleCheck = (*versionCheck)(nil)

const (
	versionsFile = "versions.json"
	centralFile = "kubernetes/_central-cluster/stackrox/central/deployment-central.yaml"

	yamlVersionKey = "app.kubernetes.io/version"

	scannerVersionByGithub = "https://raw.githubusercontent.com/stackrox/stackrox/%s/SCANNER_VERSION"
	collectorVersionByGithub = "https://raw.githubusercontent.com/stackrox/stackrox/%s/COLLECTOR_VERSION"
)

func (vc versionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error) {
	s := Undefined
	msgs := make([]string, 0)

	v, err := os.ReadFile(filepath.Join(extractedBundlePath, versionsFile))
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot read %q file", versionsFile))
	}

	var expectedVersions version.Versions
	err = json.Unmarshal(v, &expectedVersions)
	if err != nil {
		return checkError(errors.Wrapf(err, "cannot parse %q into expected JSON", versionsFile))
	}

	// TODO(alexr): cache HTTP client.
	httpClient := newHTTPClient()

	// Check that scanner version correspond to the GitHub one for this commit.
	scannerGithubURL := fmt.Sprintf(scannerVersionByGithub, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", scannerGithubURL)
	scannerGithubVersion, err := getXVersionFromGithub(cliEnvironment, scannerGithubURL, httpClient)
	if err != nil {
		return checkError(err)
	}
	if scannerGithubVersion != expectedVersions.ScannerVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Scanner version (GitHub): %s, Scanner version (%q): %s", scannerGithubVersion, versionsFile, expectedVersions.ScannerVersion))
	}

	// Check that collector version correspond to the GitHub one for this commit.
	collectorGithubURL := fmt.Sprintf(collectorVersionByGithub, expectedVersions.GitCommit)
	cliEnvironment.Logger().InfofLn("Querying %q", collectorGithubURL)
	collectorGithubVersion, err := getXVersionFromGithub(cliEnvironment, collectorGithubURL, httpClient)
	if err != nil {
		return checkError(err)
	}

	if collectorGithubVersion != expectedVersions.CollectorVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Collector version (GitHub): %s, Collector version (%q): %s", collectorGithubVersion, versionsFile, expectedVersions.CollectorVersion))
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
	if centralActualVersion != expectedVersions.MainVersion {
		s.AtLeast(Warning)
		msgs = append(msgs, fmt.Sprintf("Central actual version: %s, Central version (%q): %s", centralActualVersion, versionsFile, expectedVersions.MainVersion))
	}




	return s, msgs, nil
}

func (vc versionCheck) Name() string {
	return "Image Versions Match"
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

// newHTTPClient returns a new HTTP client.
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
