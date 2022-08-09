package doctor

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/roxctl/common/environment"
)

func runAllChecks(cliEnvironment environment.Environment, extractedBundlePath string) error {
	registeredChecks := []diagBundleCheck{
		versionCheck{0},
		versionCheck{1},
		versionCheck{2},
		versionCheck{42},
	}

	var numOK, numWarn, numProblem, numError int
	for _, c := range registeredChecks {
		s, msg, err := c.Run(cliEnvironment, extractedBundlePath)

		var caption string
		var details interface{}
		switch {
		case err != nil:
			numError++
			caption = "error running check => results unavailable"
			details = err
		case s == OK:
			numOK++
			caption = "OK"
		case s == Warning:
			numWarn++
			caption = "WARNING: depending on other factors, this might or might not be a problem"
			details = msg
		case s == Problem:
			numProblem++
			caption = "PROBLEM: this is likely to cause issues"
			details = msg
		}

		cliEnvironment.Logger().PrintfLn("\n[%s] %s", c.Name(), caption)
		if details != nil {
			cliEnvironment.Logger().PrintfLn("\t> %v", details)
		}

	}

	cliEnvironment.Logger().PrintfLn("\n%d checks run out of %d, %d problem(s) and %d warning(s) found",
		len(registeredChecks) - numError,
		len(registeredChecks),
		numProblem,
		numWarn,
	)

	if numProblem != 0 {
		return errors.Errorf("%d problem(s) found", numProblem)
	}
	return nil
}
