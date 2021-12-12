package deploytime

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/defaults/policies"
	"github.com/stackrox/rox/pkg/detection"
	"github.com/stackrox/rox/pkg/images/types"
	pol "github.com/stackrox/rox/pkg/policies"
)

// NOTE: This benchmark tests assumes some things that require manual interaction:
//	 		- You need to upload two images to a temporary registry (ttl.sh) and add signatures for it via cosign.
//		 	- You MUST use image sha256:4424e31f2c366108433ecca7890ad527b243361577180dfd9a5bb36e828abf47 in the temporary registry (this will otherwise not get meaningful results for your tests)
//			- Replace the values for signedContainerImage, unsignedContainerImage and wronglySignedContainerImage
//			- Within the matcher impl (pkg/booleanpolicy/matcher.go#l237) use the option signature.WithCustomKeyChain(authn.DefaultKeychain) (otherwise you can only run the benchmark in K8S envs)
//			- Within the image signature policy file under testdata/, specify your cosign public key as base64 encoded string within the value of the policy
var signedContainerImage = &storage.ContainerImage{
	Id: "ttl.sh/7f06ce8a-21b2-44f2-a096-d9416e0dace4",
	Name: &storage.ImageName{
		Registry: "ttl.sh",
		Remote:   "7f06ce8a-21b2-44f2-a096-d9416e0dace4",
		Tag:      "72h",
		FullName: "ttl.sh/7f06ce8a-21b2-44f2-a096-d9416e0dace4:72h",
	},
	NotPullable: false,
}

var unsignedContainerImage = &storage.ContainerImage{
	Id:                   "ttl.sh/967f77f5-c474-49d2-bd74-58bb7e856959",
	Name:                 &storage.ImageName{
		Registry:             "ttl.sh",
		Remote:               "967f77f5-c474-49d2-bd74-58bb7e856959",
		Tag:                  "72h",
		FullName:             "ttl.sh/967f77f5-c474-49d2-bd74-58bb7e856959:72h",
	},
	NotPullable:          false,
}

var wronglySignedContainerImage = &storage.ContainerImage{
	Id:                   "ttl.sh/6a565840-48eb-449c-8f4b-285a734ee418",
	Name:                 &storage.ImageName{
		Registry:             "ttl.sh",
		Remote:               "6a565840-48eb-449c-8f4b-285a734ee418",
		Tag:                  "72h",
		FullName:             "ttl.sh/6a565840-48eb-449c-8f4b-285a734ee418:72h",
	},
	NotPullable:          false,
}

func generateDeployments(amount int) map[*storage.Deployment][]*storage.Image {
	deploymentsAndImages := make(map[*storage.Deployment][]*storage.Image, amount)
	for i := 0; i < amount; i++ {
		var dep *storage.Deployment
		var imgs []*storage.Image
		if i % 3 == 0 {
			dep, imgs = deploymentWithUnsignedImages(i)
		} else if i % 3 == 1 {
			dep, imgs = deploymentWithWronglySignedImages(i)
		} else if i % 3 == 2 {
			dep, imgs = deploymentWithSignedImages(i)
		}
		deploymentsAndImages[dep] = imgs
	}
	return deploymentsAndImages
}

func deploymentWithUnsignedImages(i int) (*storage.Deployment, []*storage.Image) {
	d := &storage.Deployment{
		Id:          fmt.Sprintf("dep-%d", i),
		Name:        fmt.Sprintf("dep-%d", i),
		Replicas:    1,
		ClusterId:   "cluster",
		ClusterName: "cluster",
		Containers: []*storage.Container{{
			Id:    "c1",
			Image: unsignedContainerImage,
			Name:  "c1",
		}},
	}
	img := types.ToImage(unsignedContainerImage)
	return d, []*storage.Image{img}
}

func deploymentWithWronglySignedImages(i int) (*storage.Deployment, []*storage.Image) {
	d := &storage.Deployment{
		Id:          fmt.Sprintf("dep-%d", i),
		Name:        fmt.Sprintf("dep-%d", i),
		Replicas:    1,
		ClusterId:   "cluster",
		ClusterName: "cluster",
		Containers: []*storage.Container{{
			Id:    "c1",
			Image: wronglySignedContainerImage,
			Name:  "c1",
		}},
	}
	img := types.ToImage(wronglySignedContainerImage)
	return d, []*storage.Image{img}
}

func deploymentWithSignedImages(i int) (*storage.Deployment, []*storage.Image) {
	d := &storage.Deployment{
		Id:          fmt.Sprintf("dep-%d", i),
		Name:        fmt.Sprintf("dep-%d", i),
		Replicas:    1,
		ClusterId:   "cluster",
		ClusterName: "cluster",
		Containers: []*storage.Container{{
			Id:    "c1",
			Image: signedContainerImage,
			Name:  "c1",
		}},
	}
	img := types.ToImage(signedContainerImage)
	return d, []*storage.Image{img}
}

func createDefaultPoliciesSet() (detection.PolicySet, error) {
	policies, err := policies.DefaultPolicies()
	if err != nil {
		return nil, err
	}

	set := detection.NewPolicySet()

	for _, policy := range policies {
		if err := set.UpsertPolicy(policy); err != nil {
			return nil, err
		}
	}
	return set, nil
}

func createDefaultPoliciesAndImageSignaturePoliciesSet() (detection.PolicySet, error) {
	set, err := createDefaultPoliciesSet()
	if err != nil {
		return nil, err
	}

	p, err := createImageSignaturePolicy()
	if err != nil {
		return nil, err
	}

	err = set.UpsertPolicy(p)
	if err != nil {
		return nil, err
	}

	return set, nil
}

func createImageSignaturePolicy() (*storage.Policy, error) {
	contents, err := os.ReadFile("testdata/image-signature-policy.json")
	if err != nil {
		return nil, err
	}

	var policy storage.Policy
	err = jsonpb.Unmarshal(bytes.NewReader(contents), &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func benchmarkPolicyEvaluation(d Detector, depsAndImgs map[*storage.Deployment][]*storage.Image, b *testing.B) {
	for n := 0; n < b.N; n++ {
		for dep, images := range depsAndImgs {
			_, err := d.Detect(DetectionContext{EnforcementOnly: false}, dep, images, func(policy *storage.Policy) bool {
				return pol.AppliesAtDeployTime(policy)
			})
			if err != nil {
				b.Logf("found an unexpected error: %v", err)
			}
		}
	}
}

func BenchmarkImageSignaturePolicyAndDefaultPolicies(b *testing.B) {
	depsAndImgs := generateDeployments(500)
	set, err := createDefaultPoliciesAndImageSignaturePoliciesSet()
	if err != nil {
		b.Fatalf("Found an unexpected error during setup: %v", err)
	}
	detector := NewDetector(set)
	b.ResetTimer()
	benchmarkPolicyEvaluation(detector, depsAndImgs, b)
}

func BenchmarkDefaultPolicies(b *testing.B) {
	depsAndImgs := generateDeployments(500)
	set, err := createDefaultPoliciesSet()
	if err != nil {
		b.Fatalf("Found an unexpected error during setup: %v", err)
	}
	detector := NewDetector(set)
	b.ResetTimer()
	benchmarkPolicyEvaluation(detector, depsAndImgs, b)
}
