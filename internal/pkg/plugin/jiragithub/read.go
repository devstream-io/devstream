package jiragithub

import (
	"github.com/merico-dev/stream/pkg/util/github"
)

// Read get jira-github-integ workflows.
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
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	path, err := ghClient.GetWorkflowPath()
	if err != nil {
		return nil, err
	}

	return BuildReadState(path), nil
}

func BuildReadState(path string) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = path
	return res
}
