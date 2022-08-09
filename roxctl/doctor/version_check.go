package doctor

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/roxctl/common/environment"
)

type versionCheck struct {
	typ int
}

var _ diagBundleCheck = (*versionCheck)(nil)

func (vc versionCheck) Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, string, error) {
	switch vc.typ {
	case 0: return OK, "", nil
	case 1: return Warning, "warning msg", nil
	case 2: return Problem, "PROBLEM msg", nil
	default:
		return Undefined, "", errors.New("unknown typ")
	}
}

func (vc versionCheck) Name() string {
	return fmt.Sprintf("Version check - %v", vc.typ)
}
