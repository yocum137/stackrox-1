package lint

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/golang/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/protowalk"
	"github.com/stackrox/rox/pkg/sliceutils"
	scannerV1 "github.com/stackrox/scanner/generated/scanner/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	// Ensure proto files are added to the proto registry
	_ "github.com/stackrox/rox/generated/api/v1"
)

var (
	// This regex matches all words or word fragments that suggest a field refers to an endpoint.
	endpointKeywordRegex = regexp.MustCompile(`addr|host|endpoint|\burls?\b|\bips?\b|registry|server`)

	// This is an allowlist of fields for which no annotation is required.
	// If you add an entry to a list, please explain why validation can be skipped for this field in a comment.
	noAnnotationAllowedFields = map[interface{}][]string{
		// UI endpoints are used by the client only, login_url set by server only.
		(*storage.AuthProvider)(nil): {"ui_endpoint", "extra_ui_endpoints", "login_url"},
		(*storage.Notifier)(nil):     {"ui_endpoint"},
		// Not used as connection targets by Central.
		(*storage.Cluster)(nil):              {"central_api_endpoint"},
		(*storage.DynamicClusterConfig)(nil): {"registry_override"},
		(*storage.StaticClusterConfig)(nil):  {"central_api_endpoint"},
		// Not contacted by neither Central nor Scanner.
		(*scannerV1.PythonComponent)(nil): {"download_url"},
		// ID only, endpoint stored in separate 'endpoint' field.
		(*storage.ECRConfig)(nil): {"registry_id"},
		// False positive matches.
		(*storage.SecurityContext_SeccompProfile)(nil): {"localhost_profile"},
		(*storage.PortConfig_ExposureInfo)(nil):        {"external_hostnames", "external_ips", "service_cluster_ip"},
		(*storage.Deployment)(nil):                     {"host_network", "host_pid", "host_ipc"},
		// Matching only.
		(*storage.VulnerabilityRequest_Scope_Image)(nil): {"registry"},
		// TODO: revisit
		(*storage.ImageName)(nil): {"registry"},
	}
)

func TestEndpointFieldsInServiceInputsHaveValidateTag(t *testing.T) {
	serviceInputTypes := getServiceInputTypes()
	for inputTy, svcs := range serviceInputTypes {
		var offendingFields []string
		protowalk.WalkProto(inputTy, func(fieldPath protowalk.FieldPath) bool {
			f := fieldPath.Field()
			if !f.IsLeaf() {
				return true
			}

			fieldName := f.ProtoName()
			if strings.HasPrefix(fieldName, "SORT_") {
				// Ignore sort fields
				return false
			}

			// Convert the field name to snake case, and then replace all underscores with spaces (as _ is a word
			// character from a regex perspective and thus isn't matched by \b).
			fieldName = snakeCaseify(fieldName)
			if !endpointKeywordRegex.MatchString(strings.ReplaceAll(fieldName, "_", " ")) {
				return false
			}

			// Check if field is allowlisted
			if allowedFields := noAnnotationAllowedFields[reflect.Zero(f.ContainingType).Interface()]; sliceutils.StringFind(allowedFields, fieldName) != -1 {
				return false
			}

			validateTag := fieldPath.Field().Tag.Get("validate")
			if validateTag == "" {
				offendingFields = append(offendingFields, fmt.Sprintf(" - field %s of type %s (used as %s in input)", f.ProtoName(), f.ContainingType, fieldPath.ProtoPath()))
			}
			return true
		})
		if len(offendingFields) > 0 {
			assert.Failf(t, "fields with missing validate tag found",
				"Type %s used as input for method %s (and %d others)\ncontains fields that seem to refer to endpoints/URLs, but are missing 'validate' tags:\n"+
					"%s\n"+
					"For proper endpoint validation, all such fields must be annotated with `validate:\"nolocalendpoints\"`.\n"+
					"If you believe this to be a false positive, add a `validate:\"-\"` tag, or add the field to the allowlist\n"+
					"at the top of this test file.",
				inputTy, svcs[0].FullName(), len(svcs)-1, strings.Join(offendingFields, "\n"))
		}
	}
}

/////////////////////////////////////////////////////
// Helpers                                         //
/////////////////////////////////////////////////////

// getServiceInputTypes returns a map containing as keys all types that are used as inputs by proto methods, and
// mapping them to the methods that use them.
func getServiceInputTypes() map[reflect.Type][]protoreflect.MethodDescriptor {
	result := make(map[reflect.Type][]protoreflect.MethodDescriptor)
	protoregistry.GlobalFiles.RangeFilesByPackage("v1", func(fileDesc protoreflect.FileDescriptor) bool {
		for i := 0; i < fileDesc.Services().Len(); i++ {
			svc := fileDesc.Services().Get(i)
			for j := 0; j < svc.Methods().Len(); j++ {
				method := svc.Methods().Get(j)
				ty := proto.MessageType(string(method.Input().FullName()))
				if ty == nil {
					panic(fmt.Errorf("service method %s has unknown input type %s", method.FullName(), method.Input().FullName()))
				}
				result[ty] = append(result[ty], method)
			}
		}
		return true
	})
	return result
}

// snakeCaseify attempts to translate the input word (which can be in either snake or camel case) to snake case.
func snakeCaseify(input string) string {
	var words []string
	var currWord strings.Builder
	var maybeNextWord rune

	for _, c := range input {
		lc := unicode.ToLower(c)
		if c == lc && (!unicode.IsDigit(c) || maybeNextWord == 0) {
			if maybeNextWord != 0 {
				if currWord.Len() > 0 {
					words = append(words, currWord.String())
					currWord.Reset()
				}
				currWord.WriteRune(maybeNextWord)
				maybeNextWord = 0
			}
			currWord.WriteRune(lc)
		} else {
			if maybeNextWord != 0 {
				currWord.WriteRune(maybeNextWord)
			} else {
				if currWord.Len() > 0 {
					words = append(words, currWord.String())
					currWord.Reset()
				}
			}
			maybeNextWord = lc
		}
	}
	if maybeNextWord != 0 {
		currWord.WriteRune(maybeNextWord)
	}
	if currWord.Len() > 0 {
		words = append(words, currWord.String())
	}
	return strings.Join(words, "_")
}
