package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

// InstallRepo will install repo by opts config
func InstallRepo(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Create and render repo get from given url
	dirPath, err := opts.CreateAndRenderLocalRepo()
	if err != nil {
		return err
	}

	// 2. Push local repo to remote
	switch opts.DestinationRepo.RepoType {
	case "github":
		err = opts.PushToRemoteGithub(dirPath)
	case "gitlab":
		err = opts.PushToRemoteGitlab(dirPath)
	default:
		err = fmt.Errorf("scaffolding not support repo destination: %s", opts.DestinationRepo.RepoType)
	}
	if err != nil {
		return err
	}
	return nil
}

// DeleteRepo will delete repo by options
func DeleteRepo(options plugininstaller.RawOptions) error {
	var err error
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	dstRepo := opts.DestinationRepo
	switch dstRepo.RepoType {
	case "github":
		// 1. create ghClient
		ghClient, err := dstRepo.CreateGithubClient(true)
		if err != nil {
			return err
		}
		// 2. delete github repo
		return ghClient.DeleteRepo()
	case "gitlab":
		gLclient, err := dstRepo.CreateGitlabClient()
		if err != nil {
			return err
		}
		return gLclient.DeleteProject(dstRepo.PathWithNamespace)
	}
	return fmt.Errorf("scaffolding not support repo destination: %s", dstRepo.RepoType)
}
