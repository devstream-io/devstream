package java

import "github.com/devstream-io/devstream/internal/pkg/statemanager"

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

func buildStatus(opts *Options) statemanager.ResourceStatus {
	resStatus := make(statemanager.ResourceStatus)
	resStatus["pathWithNamespace"] = opts.PathWithNamespace
	resStatus["branch"] = opts.Branch
	return resStatus
}
