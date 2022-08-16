package quay

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/stringutils"
	"github.com/stretchr/testify/assert"
)

const (
	// This is a robot token that can only pull from quay.io/integration/nginx
	testOauthToken = "0j9dhT9jCNFpsVAzwLavnyeEy2HWnrfTQnbJgQF8"
)

func TestQuay(t *testing.T) {
	integration := &storage.ImageIntegration{
		IntegrationConfig: &storage.ImageIntegration_Quay{
			Quay: &storage.QuayConfig{
				Endpoint: "https://quay.io",
			},
		},
	}

	q, err := newRegistry(integration)
	assert.NoError(t, err)
	assert.NoError(t, filterOkErrors(q.Test()))

	data, err := ioutil.ReadFile("/Users/connorgorman/go/src/github.com/stackrox/stackrox/top50taggedimages")
	if err != nil {
		panic(err)
	}
	file, err := os.Create("/Users/connorgorman/go/src/github.com/stackrox/stackrox/quay-output")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	digestSet := set.NewStringSet()
	for _, image := range strings.Split(string(data), "\n") {
		remote, tag := stringutils.Split2(image, ":")
		img := &storage.Image{
			Name: &storage.ImageName{
				Registry: "quay.io",
				Remote:   remote,
				Tag:      tag,
			},
		}
		meta, err := q.Metadata(img)
		if err != nil {
			panic(err)
		}

		digest := meta.GetV2().GetDigest()
		if digest == "" {
			digest = meta.V1.GetDigest()
			if digest == "" {
				fmt.Println("skipping", image)
				continue
			}
		}
		if !digestSet.Add(digest) {
			continue
		}
		fmt.Fprintln(file, "quay.io/"+image+"@"+digest)
	}
}

func filterOkErrors(err error) error {
	if err != nil &&
		(strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "status=502")) {
		// Ignore failures that can indicate quay.io outage
		return nil
	}
	return err
}
