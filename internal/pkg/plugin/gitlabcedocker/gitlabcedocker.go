package gitlabcedocker

// TODO(dtm): Add your logic here.

const gitlabImageName = "gitlab/gitlab-ce:rc"

var (
	gitlabContainerName = "gitlab"
)

// dockerOperator is an interface for docker operations
// It is implemented by sshDockerOperator
// in the future, we can add other implementations such as sshDockerOperator
type dockerOperator interface {
	IfImageExists(imageName string) bool
	PullImage(image string) error
	RemoveImage(image string) error

	IfContainerExists(container string) bool
	IfContainerRunning(container string) bool
	RunContainer(options Options) error
	StopContainer(container string) error
	RemoveContainer(container string) error

	ListContainerMounts(container string) ([]string, error)
}

func getDockerOperator(_ Options) dockerOperator {
	// just return a sshDockerOperator for now
	return &sshDockerOperator{}
}

func buildState(containerRunning bool, volumes []string) map[string]interface{} {
	return map[string]interface{}{
		"containerRunning": containerRunning,
		"volumes":          volumes,
	}
}
