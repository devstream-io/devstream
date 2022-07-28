package jenkinsgithub

import "github.com/devstream-io/devstream/pkg/util/jenkins"

// TODO(dtm): Add your logic here.

const (
	jenkinsCredentialID       = "credential-by-devstream-jenkins-github-integ"
	jenkinsCredentialDesc     = "Jenkins Pipeline secret, created by devstream/jenkins-github-integ"
	jenkinsCredentialUsername = "github-by-devstream-jenkins-github-integ"
	githubAuthID              = "3a3b9ece-ad38-4209-8808-a37fbe74cc95"
)

// NewJenkinsFromOptions creates a Jenkins client from the given options and test the connection.
func NewJenkinsFromOptions(opts *Options) (*jenkins.Jenkins, error) {
	return jenkins.NewJenkins(opts.J.URL, opts.J.User, opts.J.Password)
}

type resource struct {
	CredentialsCreated bool
	JobCreated         bool
}

func buildResource(res *resource, pluginMap map[string]bool) map[string]interface{} {
	m := map[string]interface{}{
		"credentialsCreated": res.CredentialsCreated,
		"jobCreated":         res.JobCreated,
	}
	for k, v := range pluginMap {
		m[k] = v
	}

	return m
}
