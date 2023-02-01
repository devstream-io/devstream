package param

import "github.com/devstream-io/devstream/pkg/util/log"

type Param struct {
	GithubUsername     string
	GitHubRepo         string
	GithubToken        string
	DockerhubUsername  string
	DockerhubToken     string
	Language           string
	Framework          string
	RepoScaffoldingURL string
}

func GetParams() (*Param, error) {
	lang, frame, url, err := selectRepoScaffolding()
	if err != nil {
		return nil, err
	}

	githubUsername, err := getGitHubUsername()
	if err != nil {
		return nil, err
	}

	githubRepo, err := getGitHubRepo(lang, frame)
	if err != nil {
		return nil, err
	}

	githubToken, err := getGitHubToken()
	if err != nil {
		return nil, err
	}

	dockerhubUsername, err := getDockerHubUsername()
	if err != nil {
		return nil, err
	}

	dockerhubToken, err := getDockerHubToken()
	if err != nil {
		return nil, err
	}

	param := &Param{
		Language:           lang,
		Framework:          frame,
		RepoScaffoldingURL: url,
		GithubUsername:     githubUsername,
		GitHubRepo:         githubRepo,
		GithubToken:        githubToken,
		DockerhubUsername:  dockerhubUsername,
		DockerhubToken:     dockerhubToken,
	}
	// TODO: change to debug level
	log.Infof("param: %+v", param)
	return param, nil
}
