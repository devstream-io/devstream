package java

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

var (
	defaultTags = "gitlab-java"

	defaultMVNPackageJobImg    = "maven:3.6.2-jdk-14"
	defaultMVNPackageJobScript = "mvn clean package -B"

	defaultDockerBuildJobImg = "docker:latest"

	defaultK8sDeployJobImg = "bitnami/kubectl:latest"
)

func buildState(opts *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opts.PathWithNamespace,
		"branch":            opts.Branch,
	}
}
