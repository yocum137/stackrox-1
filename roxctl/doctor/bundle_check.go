package doctor

import (
	"github.com/stackrox/rox/roxctl/common/environment"
)

type CheckStatus int
const (
	Undefined CheckStatus = iota
	OK
	Warning
	Problem
)

type diagBundleCheck interface {
	// TODO(alexr): make error a slice of errors
	Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, []string, error)
	Name() string
}

func checkError(err error) (CheckStatus, []string, error) {
	return Undefined, nil, err
}

func (before *CheckStatus) AtLeast(suggested CheckStatus) {
	if suggested > *before {
		*before = suggested
	}
}
