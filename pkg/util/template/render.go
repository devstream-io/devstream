package template

import (
	"bytes"
	"html/template"
)

func Render(name, templateStr string, variable any) (string, error) {
	t, err := template.New(name).Delims("[[", "]]").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err = t.Execute(&buff, variable); err != nil {
		return "", err
	}
	return buff.String(), nil
}
