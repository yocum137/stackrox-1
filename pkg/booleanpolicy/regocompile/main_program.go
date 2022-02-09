package regocompile

import (
	"strings"
	"text/template"
)

var (
	mainProgramTemplate = template.Must(template.New("").Parse(`
package policy.main

# Custom utility functions.
and(args) = result {
	result := sum([argAsInt | arg := args[_]; argAsInt := to_number(arg)]) == count(args)
}


{{ $root := . }}

{{- range .Functions }}
{{.}}
{{- end }}

{{- range .Conditions }}
violations[result] {
	{{- range $root.IndexesToDeclare }}
	some idx{{.}}
	{{- end }}
	{{- range $field := .Fields }}
	{{- range .FuncNames }}
	{{.}}Result := {{ .}}(input.{{ $field.JSONPath }}) 
	{{.}}Result["match"]{{- if $field.Negate }} == false {{- end }}
	{{- end }}
	{{- end }}
	result := {
		{{- range $fieldIndex, $field := .Fields }}
			{{- if $fieldIndex}},{{- end }} 
			"{{ $field.Name }}": 
				{{- range $funcNameIndex, $funcName := $field.FuncNames }}
					{{- if $funcNameIndex}} &{{ end }} {{ $funcName }}Result["values"]
				{{- end }}
		{{- end }}
	}
}
{{- end }}
`))
)

type fieldInCondition struct {
	Name      string
	JSONPath  string
	Negate    bool
	FuncNames []string
}

type condition struct {
	Fields []fieldInCondition
}

type mainProgramArgs struct {
	IndexesToDeclare []int
	Functions        []string
	Conditions       []condition
}

func generateMainProgram(args *mainProgramArgs) (string, error) {
	var sb strings.Builder
	if err := mainProgramTemplate.Execute(&sb, args); err != nil {
		return "", err
	}
	return sb.String(), nil
}
