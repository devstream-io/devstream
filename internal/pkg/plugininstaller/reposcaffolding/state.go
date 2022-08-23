package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func GetDynamicState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	dstRepo := opts.DestinationRepo
	switch dstRepo.RepoType {
	case "github":
		return getGithubStatus(dstRepo)
	case "gitlab":
		return getGitlabStatus(dstRepo)
	default:
		return nil, fmt.Errorf("read state not support repo type: %s", dstRepo.RepoType)
	}

}

func getGithubStatus(dstRepo *common.Repo) (map[string]interface{}, error) {
	repoInfo := dstRepo.BuildRepoInfo()
	ghClient, err := github.NewClient(repoInfo)
	if err != nil {
		return nil, err
	}

	repo, err := ghClient.DescribeRepo()
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, nil
	}

	res := make(map[string]interface{})
	res["owner"] = dstRepo.Owner
	res["org"] = dstRepo.Org
	res["repoName"] = *repo.Name

	outputs := make(map[string]interface{})

	if dstRepo.Owner == "" {
		outputs["owner"] = dstRepo.Owner
	} else {
		outputs["owner"] = *repo.Owner.Login
	}
	if dstRepo.Org == "" {
		outputs["org"] = dstRepo.Org
	} else {
		outputs["org"] = *repo.Organization.Login
	}
	outputs["repo"] = dstRepo.Repo
	outputs["repoURL"] = *repo.CloneURL
	res["outputs"] = outputs

	return res, nil

}

func getGitlabStatus(dstRepo *common.Repo) (map[string]interface{}, error) {
	c, err := gitlab.NewClient(dstRepo.BuildRepoInfo())
	if err != nil {
		return nil, err
	}

	project, err := c.DescribeRepo()
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	res := make(map[string]interface{})
	outputs := make(map[string]interface{})

	log.Debugf("GitLab Project is: %#v\n", project)

	if project.Owner != nil {
		log.Debugf("GitLab project owner is: %#v.\n", project.Owner)
		res["owner"] = project.Owner.Username
		res["org"] = project.Owner.Organization
		outputs["owner"] = project.Owner.Username
		outputs["org"] = project.Owner.Organization
	} else {
		res["owner"] = dstRepo.Owner
		res["org"] = dstRepo.Org
		outputs["owner"] = dstRepo.Owner
		outputs["org"] = dstRepo.Org
	}
	res["repoName"] = project.Name
	outputs["repo"] = project.Name
	outputs["repoURL"] = project.HTTPURLToRepo
	res["outputs"] = outputs

	return res, nil
}
