package java

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed gitlabci-template.yml
var gitlabCITemplate string

// Render gitlab-ci.yml template with Options
func renderTmpl(Opts *Options) (string, error) {
	t := template.Must(template.New("gitlabci-java").Option("missingkey=error").Parse(gitlabCITemplate))

	var buf bytes.Buffer
	err := t.Execute(&buf, Opts)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
