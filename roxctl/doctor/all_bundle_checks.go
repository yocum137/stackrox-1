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

		var caption string
		switch {
		case err != nil:
			numError++
			caption = "ERROR running check"
			msg = append([]string{fmt.Sprintf("Error: %v\n", err)}, msg...)
		case s == OK:
			numOK++
			caption = "OK"
		case s == Warning:
			numWarn++
			caption = "WARNING"
		case s == Problem:
			numProblem++
			caption = "PROBLEM"
		default:
			numError++
			caption = "?"
		}

		_, _ = fmt.Fprintf(cliEnvironment.ColorWriter(), "[%s] %s\n", c.Name(), caption)
		for _, m := range msg {
			_, _ = fmt.Fprintf(cliEnvironment.ColorWriter(), "\t> %s\n", m)
		}
	}

	_, _ = fmt.Fprintf(cliEnvironment.ColorWriter(), "\n%d check(s) run out of %d, %d problem(s) and %d warning(s) found\n",
		len(registeredChecks) - numError,
		len(registeredChecks),
		numProblem,
		numWarn,
	)

	_, _ = fmt.Fprint(cliEnvironment.ColorWriter(), "\nLegend:\n" +
		"\tWARNING: depending on other factors, this might or might not be a problem\n" +
		"\tPROBLEM: this is likely to cause issues\n" +
		"\tERROR running check: the check could not run => no result\n")

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
