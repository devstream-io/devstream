package java

import (
	"bytes"
	"html/template"
)

var gitlab_ci_template = `
image: docker:stable

stages:
  - package
  - docker_build
  - k8s_deploy

{{if .Package.Enable}}
mvn_package_job:
  image: {{.Package.BaseOption.Image}}
  stage: package
  tags:
    - {{.Package.BaseOption.Tags}}
  script: 
    {{range .Package.ScriptCommand}}
    {{.}}
    {{end}} 
  artifacts:
    paths:
      - target/*.jar
{{end}}

{{if .Build.Enable}}
docker_build_job:
  image: {{.Build.BaseOption.Image}}
  stage: docker_build
  tags: 
    - {{.Build.BaseOption.Tags}}
  script:
    {{range .Build.ScriptCommand}}
    {{.}}
    {{end}} 
{{end}}

{{if .Deploy.Enable}}
k8s_deploy_job:
  image: 
    name: {{.Deploy.BaseOption.Image}}
    entrypoint: [""]
  stage: k8s_deploy
  tags: 
    - {{.Deploy.BaseOption.Tags}}
  script:
    {{range .Deploy.ScriptCommand}}
    {{.}}
    {{end}} 
{{end}}
`

// Render gitlab-ci.yml template with Options
func renderTmpl(Opts *Options) (string, error) {
	t, err := template.New("gitlabci-java").Option("missingkey=error").Parse(gitlab_ci_template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, Opts)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
