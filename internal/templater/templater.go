package templater

import (
	"bytes"
	"text/template"
)

// Render renders a template string.
func Render(tmpl string) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		return "", err
	}

	return buf.String(), nil
}
