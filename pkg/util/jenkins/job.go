package jenkins

import (
	"context"
	_ "embed"

	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"

	"github.com/devstream-io/devstream/pkg/util/template"
)

var (
	errorNotFound = errors.New("404")
)

//go:embed tpl/seedjob.tpl.groovy
var jobGroovyScript string

func (jenkins *jenkins) GetFolderJob(jobName string, jobFolder string) (*gojenkins.Job, error) {
	if jobFolder != "" {
		return jenkins.GetJob(context.Background(), jobName, jobFolder)
	}
	return jenkins.GetJob(context.Background(), jobName)
}

func BuildRenderedScript(vars any) (string, error) {
	return template.Render("jenkins-script-template", jobGroovyScript, vars)
}

func IsNotFoundError(err error) bool {
	if err != nil {
		return err.Error() == errorNotFound.Error()
	}
	return false
}
