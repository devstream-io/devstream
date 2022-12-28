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

// JobScriptRenderInfo is used to render jenkins job groovy script
type JobScriptRenderInfo struct {
	// jenkins related info
	FolderName string
	JobName    string
	// repo related info
	RepoCredentialsId string
	Branch            string
	RepoType          string
	RepoURL           string
	RepoName          string
	RepoOwner         string
	RepositoryURL     string
	SecretToken       string
	GitlabConnection  string
}

// JenkinsFileRenderInfo is used to render Jenkinsfile
type JenkinsFileRenderInfo struct {
	AppName string `mapstructure:"AppName"`
	// imageRepo variables
	ImageRepositoryURL  string `mapstructure:"ImageRepositoryURL"`
	ImageAuthSecretName string `mapstructure:"ImageAuthSecretName"`
	// dingtalk variables
	DingtalkRobotID string `mapstructure:"DingtalkRobotID"`
	DingtalkAtUser  string `mapstructure:"DingtalkAtUser"`
	// sonarqube variables
	SonarqubeEnable bool `mapstructure:"SonarqubeEnable"`
	// custom variables
	Custom map[string]interface{} `mapstructure:"Custom"`
}

func (jenkins *jenkins) GetFolderJob(jobName string, jobFolder string) (*gojenkins.Job, error) {
	if jobFolder != "" {
		return jenkins.GetJob(context.Background(), jobName, jobFolder)
	}
	return jenkins.GetJob(context.Background(), jobName)
}

func BuildRenderedScript(vars any) (string, error) {
	return template.NewRenderClient(
		&template.TemplateOption{Name: "jenkins-script-template"}, template.ContentGetter,
	).Render(jobGroovyScript, vars)
}

func IsNotFoundError(err error) bool {
	if err != nil {
		return err.Error() == errorNotFound.Error()
	}
	return false
}
