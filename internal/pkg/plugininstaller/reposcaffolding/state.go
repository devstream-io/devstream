package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

func GetDynamicState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
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

	resState := make(statemanager.ResourceState)
	resState["owner"] = dstRepo.Owner
	resState["org"] = dstRepo.Org
	resState["repoName"] = *repo.Name

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
	resState.SetOutputs(outputs)

	return resState, nil
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

	resState := make(statemanager.ResourceState)
	outputs := make(map[string]interface{})

	log.Debugf("GitLab Project is: %#v\n", project)

	if project.Owner != nil {
		log.Debugf("GitLab project owner is: %#v.\n", project.Owner)
		resState["owner"] = project.Owner.Username
		resState["org"] = project.Owner.Organization
		outputs["owner"] = project.Owner.Username
		outputs["org"] = project.Owner.Organization
	} else {
		resState["owner"] = dstRepo.Owner
		resState["org"] = dstRepo.Org
		outputs["owner"] = dstRepo.Owner
		outputs["org"] = dstRepo.Org
	}
	resState["repoName"] = project.Name
	outputs["repo"] = project.Name
	outputs["repoURL"] = project.HTTPURLToRepo
	resState.SetOutputs(outputs)

	return resState, nil
}
