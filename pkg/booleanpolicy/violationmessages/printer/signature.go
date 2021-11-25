package printer

const (
	signatureTemplate = `Images within deployment are not signed.`
)

func imageSignaturePrinter(fieldMap map[string][]string) ([]string, error) {
	type resultFields struct {
	}

	r := resultFields{}

	return executeTemplate(signatureTemplate, r)
}
