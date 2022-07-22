package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

func GetStaticState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	dstRepo := opts.DestinationRepo
	res := make(map[string]interface{})
	res["owner"] = dstRepo.Owner
	res["org"] = dstRepo.Org
	res["repoName"] = dstRepo.Repo

	outputs := make(map[string]interface{})
	outputs["owner"] = dstRepo.Owner
	outputs["org"] = dstRepo.Org
	outputs["repo"] = dstRepo.Repo
	if dstRepo.Owner != "" {
		outputs["repoURL"] = fmt.Sprintf("https://github.com/%s/%s.git", dstRepo.Owner, dstRepo.Repo)
	} else {
		outputs["repoURL"] = fmt.Sprintf("https://github.com/%s/%s.git", dstRepo.Org, dstRepo.Repo)
	}
	res["outputs"] = outputs
	return res, nil
}

func GetDynamicState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	switch opts.RepoType {
	case "github":
		return getGithubStatus(&opts.DestinationRepo)
	case "gitlab":
		// TODO: add gitlab status logic
		return nil, nil
	default:
		return nil, fmt.Errorf("read state not support repo type: %s", opts.RepoType)
	}

}

func getGithubStatus(dstRepo *DstRepo) (map[string]interface{}, error) {
	ghClient, err := dstRepo.createGithubClient(true)
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
