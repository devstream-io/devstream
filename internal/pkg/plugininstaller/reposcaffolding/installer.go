package reposcaffolding

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// InstallRepo will install repo by opts config
func InstallRepo(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Create temp dir
	dirName, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(dirName); err != nil {
			log.Errorf("Failed to clear workpath %s: %s.", dirName, err)
		}
	}()

	// 2. Create and render repo get from given url
	err = opts.CreateAndRenderLocalRepo(dirName)
	if err != nil {
		return err
	}

	// 2. Push local repo to remote
	repoLoc := filepath.Join(dirName, opts.DestinationRepo.Repo)
	switch opts.RepoType {
	case "github":
		err = opts.PushToRemoteGithub(repoLoc)
	case "gitlab":
		err = opts.PushToRemoteGitlab(repoLoc)
	default:
		err = fmt.Errorf("scaffolding not support repo destination: %s", opts.RepoType)
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

	switch opts.RepoType {
	case "github":
		// 1. create ghClient
		ghClient, err := opts.DestinationRepo.createGithubClient(true)
		if err != nil {
			return err
		}
		// 2. delete github repo
		return ghClient.DeleteRepo()
	case "gitlab":
		dstRepo := opts.DestinationRepo
		gLclient, err := dstRepo.createGitlabClient()
		if err != nil {
			return err
		}
		return gLclient.DeleteProject(dstRepo.PathWithNamespace)
	}
	return fmt.Errorf("scaffolding not support repo destination: %s", opts.RepoType)
}
