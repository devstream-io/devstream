package ci

type ciRepoType string

const (
	defaultBranch                         = "feat-repo-ci-update"
	defaultCommitMsg                      = "update ci config"
	ciJenkinsConfigLocation               = "Jenkinsfile"
	ciGitHubWorkConfigLocation            = ".github/workflows"
	ciGitLabConfigLocation                = ".gitlab-ci.yml"
	ciJenkinsType              ciRepoType = "jenkins"
	ciGitLabType               ciRepoType = "gitlab"
	ciGitHubType               ciRepoType = "github"
	deleteCommitMsg                       = "delete ci files"
)
