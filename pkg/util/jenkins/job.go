package jenkins

import (
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

func (jenkins *jenkins) CreateOrUpdateJob(config, jobName string) (job *gojenkins.Job, created bool, err error) {
	// create or update
	job, err = jenkins.GetJob(jenkins.ctx, jobName)
	if isNotFoundError(err) {
		job, err = jenkins.CreateJob(jenkins.ctx, config, jobName)
		created = true
		return job, true, errors.WithStack(err)
	} else if err != nil {
		return job, false, errors.WithStack(err)
	}

	err = job.UpdateConfig(jenkins.ctx, config)
	return job, false, errors.WithStack(err)
}

func BuildRenderedScript(vars any) (string, error) {
	return template.Render("jenkins-script-template", jobGroovyScript, vars)
}

func isNotFoundError(err error) bool {
	if err != nil {
		return err.Error() == errorNotFound.Error()
	}
	return false
}
