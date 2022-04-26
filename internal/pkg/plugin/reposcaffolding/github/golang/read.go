package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	rs "github.com/devstream-io/devstream/internal/pkg/plugin/common/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Read check the health for github-repo-scaffolding-golang with provided param.
func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts rs.Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := rs.Validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	return buildReadState(&opts)
}

func buildReadState(opts *rs.Options) (map[string]interface{}, error) {
	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Org:      opts.Org,
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
	res["owner"] = opts.Owner
	res["org"] = opts.Org
	res["repoName"] = *repo.Name

	outputs := make(map[string]interface{})

	if opts.Owner == "" {
		outputs["owner"] = opts.Owner
	} else {
		outputs["owner"] = *repo.Owner.Login
	}
	if opts.Org == "" {
		outputs["org"] = opts.Org
	} else {
		outputs["org"] = *repo.Organization.Login
	}
	outputs["repo"] = opts.Repo
	outputs["repoURL"] = *repo.CloneURL
	res["outputs"] = outputs

	return res, nil
}
