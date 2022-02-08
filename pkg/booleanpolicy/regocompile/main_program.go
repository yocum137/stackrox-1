package regocompile

var (
	mainProgramTemplate = `
package policy.main

{{- range .Functions }}
{{.}}
{{- end }}

{{- range .Ors }}
violations[result] {
	{{- range .SomeStatements }}
	some {{.}}
	{{- end }}
	{{- range .Fields }}
	{{.FuncName}}Result := FuncName(.JSONPath) 
	{{.FuncName}}Result["match"]
	{{- end }}
	results := {
		{{-range $index, $field := .Fields }}
			{{- if $index }},{{end }} 
			"$field.Name": {{.FuncName}}Result["values"],
		{{-end}}
	}
}
{{- end }}
`
)
