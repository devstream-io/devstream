package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Read check the health for github-repo-scaffolding-golang with provided param.
func Read(param map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(param, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("options are illegal")
	}

	return buildReadState(&opts)
}

func buildReadState(opts *Options) (map[string]interface{}, error) {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Repo:     opts.Repo,
		NeedAuth: true,
	}

	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	repo, err := ghClient.GetRepoDescription()
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, nil
	}

	res := make(map[string]interface{})
	res["owner"] = *repo.Owner.Login
	res["repoName"] = *repo.Name

	outputs := make(map[string]interface{})
	outputs["owner"] = opts.Owner
	outputs["repo"] = opts.Repo
	outputs["repoURL"] = fmt.Sprintf("https://github.com/%s/%s.git", opts.Owner, opts.Repo)

	res["outputs"] = outputs

	return res, nil
}
