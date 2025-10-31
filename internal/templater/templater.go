package templater

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// Render renders a template string.
func Render(tmpl string) (string, error) {
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}
