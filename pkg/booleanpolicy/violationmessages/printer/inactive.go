package printer

import (
	"strconv"

	"github.com/stackrox/rox/pkg/booleanpolicy/augmentedobjs"
	"github.com/stackrox/rox/pkg/search"
)

const (
	inactiveTemplate = `Deployment must have inactive set to {{ .Inactive }}`
)

func inactivePrinter(fieldMap map[string][]string) ([]string, error) {
	type resultFields struct {
		ContainerName string
		Inactive      bool
	}

	r := resultFields{}
	var err error
	r.ContainerName = maybeGetSingleValueFromFieldMap(augmentedobjs.ContainerNameCustomTag, fieldMap)
	inactive, err := getSingleValueFromFieldMap(search.Inactive.String(), fieldMap)
	if err != nil {
		return nil, err
	}
	if r.Inactive, err = strconv.ParseBool(inactive); err != nil {
		return nil, err
	}
	return executeTemplate(inactiveTemplate, r)
}
