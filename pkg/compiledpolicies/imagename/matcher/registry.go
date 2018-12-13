package matcher

import (
	"fmt"
	"regexp"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/compiledpolicies/utils"
)

func init() {
	compilers = append(compilers, newRegistryMatcher)
}

func newRegistryMatcher(policy *storage.Policy) (Matcher, error) {
	registryPolicy := policy.GetFields().GetImageName().GetRegistry()
	if registryPolicy == "" {
		return nil, nil
	}

	registryRegex, err := utils.CompileStringRegex(registryPolicy)
	if err != nil {
		return nil, err
	}
	matcher := &registryMatcherImpl{registryRegex}
	return matcher.match, nil
}

type registryMatcherImpl struct {
	registryRegex *regexp.Regexp
}

func (p *registryMatcherImpl) match(name *storage.ImageName) []*storage.Alert_Violation {
	var violations []*storage.Alert_Violation
	if name.GetRegistry() != "" && p.registryRegex.MatchString(name.GetRegistry()) {
		v := &storage.Alert_Violation{
			Message: fmt.Sprintf("Image registry matched: %s", p.registryRegex),
		}
		violations = append(violations, v)
	}
	return violations
}
