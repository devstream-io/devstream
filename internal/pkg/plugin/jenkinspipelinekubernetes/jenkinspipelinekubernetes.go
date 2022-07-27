package jenkinspipelinekubernetes

import (
	_ "embed"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

const (
	jenkinsCredentialID       = "credential-by-devstream-jenkins-pipeline-kubernetes"
	jenkinsCredentialDesc     = "Jenkins Pipeline secret, created by devstream/jenkins-pipeline-kubernetes"
	jenkinsCredentialUsername = "github-by-devstream-jenkins-pipeline-kubernetes"
)

//go:embed job-template.xml
var jobTemplate string

// NewJenkinsFromOptions creates a Jenkins client from the given options and test the connection.
func NewJenkinsFromOptions(opts *Options) (*jenkins.Jenkins, error) {
	return jenkins.NewJenkins(opts.J.URL, opts.J.User, opts.J.Password)
}

// TODO(aFlyBird0): enhance the resource fields here to be read and the way to read it, such as:
// plugins install info(GitHub Pull Request Builder Plugin and OWASP Markup Formatter must be installed)
// should we keep an eye on job configuration && status changes? maybe not.
type resource struct {
	CredentialsCreated bool
	JobCreated         bool
}

func (res *resource) toMap() map[string]interface{} {
	return map[string]interface{}{
		"credentialsCreated": res.CredentialsCreated,
		"jobCreated":         res.JobCreated,
	}
}
