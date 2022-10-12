package ci

import (
	"errors"

	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

const (
	createCommitMsg = "update ci config"
	deleteCommitMsg = "delete ci files"
	// this variable is only used for github to fork a branch and create pr
	defaultBranch = "feat-repo-ci-update"
)

func PushCIFiles(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	// 1. get git content by config
	gitMap, err := opts.buildGitMap()
	if err != nil {
		return err
	}
	//3. init git client
	gitClient, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	//4. push ci files to git repo
	_, err = gitClient.PushFiles(&git.CommitInfo{
		CommitMsg:    createCommitMsg,
		GitFileMap:   gitMap,
		CommitBranch: defaultBranch,
	}, true)
	return err
}

func DeleteCIFiles(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	// 1. get git content by config
	gitMap, err := opts.buildGitMap()
	if err != nil {
		return err
	}
	if len(gitMap) == 0 {
		return errors.New("can't get valid Jenkinsfile, please check your config")
	}
	//2. init git client
	gitClient, err := scm.NewClient(opts.ProjectRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	//3. delete ci files from git repo
	commitInfo := &git.CommitInfo{
		CommitMsg:    deleteCommitMsg,
		GitFileMap:   gitMap,
		CommitBranch: opts.ProjectRepo.Branch,
	}
	return gitClient.DeleteFiles(commitInfo)
}
