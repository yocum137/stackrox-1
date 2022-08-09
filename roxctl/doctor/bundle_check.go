package doctor

import "github.com/stackrox/rox/roxctl/common/environment"

type CheckStatus int
const (
	Undefined CheckStatus = iota
	OK
	Warning
	Problem
)

type diagBundleCheck interface {
	Run(cliEnvironment environment.Environment, extractedBundlePath string) (CheckStatus, string, error)
	Name() string
}
