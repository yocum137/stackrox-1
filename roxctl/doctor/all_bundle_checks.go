package doctor

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/roxctl/common/environment"
)

func runAllChecks(cliEnvironment environment.Environment, extractedBundlePath string) error {
	registeredChecks := []diagBundleCheck{
		collectorVersionCheck{},
		scannerVersionCheck{},
		centralVersionCheck{},
	}

	var numOK, numWarn, numProblem, numError int
	for _, c := range registeredChecks {
		s, msg, err := c.Run(cliEnvironment, extractedBundlePath)

		var caption, details string
		switch {
		case err != nil:
			numError++
			caption = "error running check => results unavailable"
			details = fmt.Sprintf("\t> Error: %v\n", err) + joinCheckMessage(msg)
		case s == OK:
			numOK++
			caption = "OK"
			details = joinCheckMessage(msg)
		case s == Warning:
			numWarn++
			caption = "WARNING: depending on other factors, this might or might not be a problem"
			details = joinCheckMessage(msg)
		case s == Problem:
			numProblem++
			caption = "PROBLEM: this is likely to cause issues"
			details = joinCheckMessage(msg)
		default:
			numError++
			caption = "?"
		}

		cliEnvironment.Logger().PrintfLn("\n[%s] %s", c.Name(), caption)
		if details != "" {
			cliEnvironment.Logger().PrintfLn("%s", details)
		}

	}

	cliEnvironment.Logger().PrintfLn("\n%d check(s) run out of %d, %d problem(s) and %d warning(s) found",
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

func joinCheckMessage(msg []string) string {
	if len(msg) == 0 {
		return ""
	}

	return "\t> " + strings.Join(msg, "\n\t> ")
}
