package golang

import (
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	gitHubClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	path, err := gitHubClient.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	if path == "" {
		// file not found
		return nil, nil
	}

	log.Debugf("Language is: %s.", ga.GetLanguage(opt.Language))

	return ga.BuildReadState(path), nil
}
