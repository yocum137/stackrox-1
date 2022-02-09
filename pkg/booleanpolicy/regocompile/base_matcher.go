package regocompile

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/pkg/booleanpolicy/query"
	"github.com/stackrox/rox/pkg/parse"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/predicate/basematchers"
)

var (
	simpleMatchFuncTemplate = template.Must(template.New("").Parse(`
{{.Name}}(val) = result {
	result := { "match": {{ .MatchCode }}, "values": {val} }
}
`))

	timestampPtrType = reflect.TypeOf((*types.Timestamp)(nil))
)

type simpleMatchFuncGenerator struct {
	Name      string
	MatchCode string
}

var (
	invalidRegoFuncNameChars = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

var (
	// ErrRegoNotYetSupported is an error that indicates that a certain query is not yet supported by rego.
	// It will be removed once rego is supported for all queries.
	ErrRegoNotYetSupported = errors.New("as-yet unsupported rego path")
)

func sanitizeFuncName(name string) string {
	return invalidRegoFuncNameChars.ReplaceAllString(name, "_")
}

// getRegoFunctionName returns a rego function name for matching the field to the given value.
// The idx is also required, and is used to ensure the function name is unique.
func getRegoFunctionName(field, value string, idx int) string {
	return sanitizeFuncName(fmt.Sprintf("match%sTo%d%s", field, idx, value))
}

func (s *simpleMatchFuncGenerator) GenerateRego() (string, error) {
	var sb strings.Builder
	err := simpleMatchFuncTemplate.Execute(&sb, s)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (s *simpleMatchFuncGenerator) FuncName() string {
	return s.Name
}

type regoMatchFuncGenerator interface {
	GenerateRego() (string, error)
	FuncName() string
}

func generateStringMatchCode(value string) (string, error) {
	negated := strings.HasPrefix(value, search.NegationPrefix)
	if negated {
		value = strings.TrimPrefix(value, search.NegationPrefix)

	}
	var matchCode string
	if strings.HasPrefix(value, search.RegexPrefix) {
		matchCode = fmt.Sprintf("regex.match(`^(?i:%s)$`, val)", strings.TrimPrefix(value, search.RegexPrefix))
	} else if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) && len(value) > 1 {
		matchCode = fmt.Sprintf(`val == "%s"`, value[1:len(value)-1])
	} else {
		matchCode = fmt.Sprintf(`startswith(val, "%s")`, value)
	}
	if negated {
		matchCode = fmt.Sprintf(`(%s) == false`, matchCode)
	}
	return matchCode, nil
}

func generateBoolMatchCode(value string) (string, error) {
	boolValue, err := parse.FriendlyParseBool(value)
	if err != nil {
		return "", err
	}
	if boolValue {
		return "val", nil
	}
	return "val == false", nil
}

func getSimpleMatchFuncGeneratorFromCode(query *query.FieldQuery, valueIndex int, matchCode string) regoMatchFuncGenerator {
	return &simpleMatchFuncGenerator{
		Name:      getRegoFunctionName(query.Field, query.Values[valueIndex], valueIndex),
		MatchCode: matchCode,
	}
}

func getSimpleMatchFuncGenerator(query *query.FieldQuery, valueIndex int, matchCodeGenerator func(string) (string, error)) (regoMatchFuncGenerator, error) {
	value := query.Values[valueIndex]
	matchCode, err := matchCodeGenerator(value)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate match code for val %s in field %s: %w", value, query.Field, err)
	}
	return getSimpleMatchFuncGeneratorFromCode(query, valueIndex, matchCode), nil
}

func getStringMatchFuncGenerator(query *query.FieldQuery, valueIndex int) (regoMatchFuncGenerator, error) {
	return getSimpleMatchFuncGenerator(query, valueIndex, generateStringMatchCode)
}

func getBoolMatchFuncGenerators(query *query.FieldQuery, valueIndex int) (regoMatchFuncGenerator, error) {
	return getSimpleMatchFuncGenerator(query, valueIndex, generateBoolMatchCode)
}

func invertCmpStr(cmpStr string) string {
	switch cmpStr {
	case "<=":
		return ">"
	case "<":
		return ">="
	case ">":
		return "<="
	case ">=":
		return "<"
	}
	return cmpStr
}

func getTimestampMatchCode(value string) (string, error) {
	if value == search.NullString {
		return `val == 0`, nil
	}

	cmpStr, value := basematchers.ParseNumericPrefix(value)
	if cmpStr == "" {
		cmpStr = "=="
	}

	timestampValue, durationValue, err := basematchers.ParseTimestampQuery(value)
	if err != nil {
		return "", err
	}
	if timestampValue != nil {
		goTime, err := types.TimestampFromProto(timestampValue)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`and([val != 0, val %s %d])`, cmpStr, goTime.UnixNano()), nil
	}

	// If we're using a duration value, we need to invert the query.
	// This is because, for example, >90d means more than 90 days ago,
	// which means <=(ts of time.Now().Add(-90days).
	durationNanos := durationValue.Nanoseconds()
	cmpStr = invertCmpStr(cmpStr)
	return fmt.Sprintf(`and([val != 0, val %s (time.now_ns() - %d)])`, cmpStr, durationNanos), nil
}

func getTimestampMatchFuncGenerator(query *query.FieldQuery, valueIndex int) (regoMatchFuncGenerator, error) {
	return getSimpleMatchFuncGenerator(query, valueIndex, getTimestampMatchCode)
}

func getPtrMatchFuncGenerator(query *query.FieldQuery, valueIndex int, typ reflect.Type) (regoMatchFuncGenerator, error) {
	// Special case for pointer to timestamp.
	if typ == timestampPtrType {
		return getTimestampMatchFuncGenerator(query, valueIndex)
	}
	value := query.Values[valueIndex]
	if value == search.NullString {
		return getSimpleMatchFuncGeneratorFromCode(query, valueIndex, "val == null"), nil
	}
	return getBaseMatchFuncGenerator(query, valueIndex, typ.Elem())
}

func getBaseMatchFuncGenerator(query *query.FieldQuery, valueIndex int, typ reflect.Type) (regoMatchFuncGenerator, error) {
	switch kind := typ.Kind(); kind {
	case reflect.String:
		return getStringMatchFuncGenerator(query, valueIndex)
	case reflect.Ptr:
		return getPtrMatchFuncGenerator(query, valueIndex, typ)
	case reflect.Array, reflect.Slice:
		// return generateSliceMatcher
	case reflect.Map:
		// return generateMapMatcher
	case reflect.Bool:
		return getBoolMatchFuncGenerators(query, valueIndex)
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		// return generateIntMatcher
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		// return generateUintMatcher
	case reflect.Float64, reflect.Float32:
		// return generateFloatMatcher
	default:
		return nil, fmt.Errorf("invalid kind for base query: %s", kind)
	}
	return nil, ErrRegoNotYetSupported
}

func getMatchFuncGenerators(query *query.FieldQuery, typ reflect.Type) ([]regoMatchFuncGenerator, error) {
	if query.MatchAll {
		return []regoMatchFuncGenerator{
			&simpleMatchFuncGenerator{Name: sanitizeFuncName(fmt.Sprintf("matchAll%s", query.Field)), MatchCode: "true"},
		}, nil
	}
	var generators []regoMatchFuncGenerator
	for i := range query.Values {
		generator, err := getBaseMatchFuncGenerator(query, i, typ)
		if err != nil {
			return nil, err
		}
		generators = append(generators, generator)
	}
	return generators, nil
}

type regoMatchFunc struct {
	functionCode string
	functionName string
}

func generateMatchersForField(fieldQuery *query.FieldQuery, typ reflect.Type) ([]regoMatchFunc, error) {
	if (fieldQuery.MatchAll && len(fieldQuery.Values) > 0) || (!fieldQuery.MatchAll && len(fieldQuery.Values) == 0) {
		return nil, errors.New("invalid number of values")
	}
	if len(fieldQuery.Values) > 1 && fieldQuery.Operator != query.Or && fieldQuery.Operator != query.And {
		return nil, fmt.Errorf("invalid operator: %s", fieldQuery.Operator)
	}

	generators, err := getMatchFuncGenerators(fieldQuery, typ)
	if err != nil {
		return nil, err
	}
	if len(generators) == 0 {
		return nil, fmt.Errorf("got no generators for fieldQuery %+v", fieldQuery)
	}
	var funcs []regoMatchFunc
	for _, gen := range generators {
		code, err := gen.GenerateRego()
		if err != nil {
			return nil, err
		}
		funcs = append(funcs, regoMatchFunc{functionCode: code, functionName: gen.FuncName()})
	}
	return funcs, nil
}
