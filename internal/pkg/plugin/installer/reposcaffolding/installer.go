package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

// InstallRepo will install repo by opts config
func InstallRepo(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	// 1. Download repo by SourceRepo
	sourceClient, err := scm.NewClient(opts.SourceRepo)
	if err != nil {
		return err
	}
	repoDir, err := sourceClient.DownloadRepo()
	if err != nil {
		log.Debugf("reposcaffolding process files error: %s", err)
		return err
	}

	// 2. render repo with variables
	appName := opts.DestinationRepo.Repo
	gitMap, err := file.GetFileMapByWalkDir(
		repoDir, filterGitFiles,
		getRepoFileNameFunc(appName, opts.SourceRepo.GetRepoNameWithBranch()),
		processRepoFileFunc(appName, opts.renderTplConfig()),
	)
	if err != nil {
		return fmt.Errorf("render RepoTemplate files failed with error: %w", err)
	}

	// 3. push repo to DestinationRepo
	dstClient, err := scm.NewClientWithAuth(opts.DestinationRepo)
	if err != nil {
		return err
	}
	return scm.PushInitRepo(dstClient, &git.CommitInfo{
		CommitMsg:    scm.DefaultCommitMsg,
		CommitBranch: scm.TransitBranch,
		GitFileMap:   gitMap,
	})
}

// DeleteRepo will delete repo by options
func DeleteRepo(options configmanager.RawOptions) error {
	var err error
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	client, err := scm.NewClientWithAuth(opts.DestinationRepo)
	if err != nil {
		return err
	}
	return client.DeleteRepo()
}
