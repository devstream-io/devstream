package reposcaffolding

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

// InstallRepo will install repo by opts config
func InstallRepo(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Download and render repo by SourceRepo
	sourceClient, err := scm.NewClient(opts.SourceRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	gitMap, err := opts.downloadAndRenderScmRepo(sourceClient)
	if err != nil {
		return err
	}

	// 2. push repo to DestinationRepo
	dstClient, err := scm.NewClient(opts.DestinationRepo.BuildRepoInfo())
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

	client, err := scm.NewClient(opts.DestinationRepo.BuildRepoInfo())
	if err != nil {
		return err
	}
	return client.DeleteRepo()
}
