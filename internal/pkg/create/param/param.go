package param

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

	githubRepo, err := getGitHubRepo()
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
	return param, nil
}
