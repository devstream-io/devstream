package reposcaffolding

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

// installRepo will install repo by opts config
func installRepo(rawOptions configmanager.RawOptions) error {
	const (
		defaultCommitMsg    = "init by devstream"
		defaultCommitBranch = "init-by-devstream"
	)
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
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
	repoNameWithBranch := fmt.Sprintf("%s-%s", opts.SourceRepo.Repo, opts.SourceRepo.Branch)
	gitMap, err := file.GetFileMapByWalkDir(
		repoDir, filterGitFiles,
		getRepoFileNameFunc(appName, repoNameWithBranch),
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
		CommitMsg:    defaultCommitMsg,
		CommitBranch: defaultCommitBranch,
		GitFileMap:   gitMap,
	})
}

// deleteRepo will delete repo by options
func deleteRepo(rawOptions configmanager.RawOptions) error {
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
		return err
	}

	client, err := scm.NewClientWithAuth(opts.DestinationRepo)
	if err != nil {
		return err
	}
	return client.DeleteRepo()
}
