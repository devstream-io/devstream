package cifile

import (
	"errors"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	createCommitMsg = "update ci config"
	deleteCommitMsg = "delete ci files"
	// this variable is only used for github to fork a branch and create pr
	defaultBranch = "feat-repo-ci-update"
)

func PushCIFiles(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	// 1. get git content by config
	gitMap, err := opts.CIFileConfig.getGitfileMap()
	if err != nil {
		return err
	}
	//2. init git client
	gitClient, err := scm.NewClientWithAuth(opts.ProjectRepo)
	if err != nil {
		return err
	}
	//3. push ci files to git repo
	_, err = gitClient.PushFiles(&git.CommitInfo{
		CommitMsg:    createCommitMsg,
		GitFileMap:   gitMap,
		CommitBranch: defaultBranch,
	}, true)
	return err
}

func DeleteCIFiles(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	// 1. get git content by config
	gitMap, err := opts.CIFileConfig.getGitfileMap()
	if err != nil {
		return err
	}
	if len(gitMap) == 0 {
		return errors.New("can't get valid Jenkinsfile, please check your config")
	}
	//2. init git client
	gitClient, err := scm.NewClientWithAuth(opts.ProjectRepo)
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
