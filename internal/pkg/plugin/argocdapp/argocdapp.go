package argocdapp

import (
	_ "embed"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

//go:embed tpl/helm.tpl.yaml
var helmApplicationConfig string

func pushArgocdConfigFiles(rawOptions configmanager.RawOptions) error {
	const createCommitMsg = "create argocdapp config files"
	const commitBranch = "feat-argocdapp-configs"
	opts, err := newOptions(rawOptions)
	if err != nil {
		return err
	}

	// 1. init scm client for check argocdapp config exists and push argocdapp files
	repoInfo := &git.RepoInfo{CloneURL: git.ScmURL(opts.Source.RepoURL)}
	scmClient, err := scm.NewClientWithAuth(repoInfo)
	if err != nil {
		return err
	}

	// 2. check argocdapp config existed in project repo
	configFileExist, err := opts.Source.checkPathExist(scmClient)
	if err != nil {
		return fmt.Errorf("argocdapp scm client get config path info failed: %w", err)
	}
	if configFileExist {
		return nil
	}

	// 3. get argocd configFiles from remote
	const configLocation = downloader.ResourceLocation("https://github.com/devstream-io/dtm-pipeline-templates.git//argocdapp/helm")
	gitFiles, err := opts.getArgocdDefaultConfigFiles(configLocation)
	if err != nil {
		return err
	}

	// 4. push git files to ProjectRepo
	_, err = scmClient.PushFiles(&git.CommitInfo{
		CommitMsg:    createCommitMsg,
		GitFileMap:   gitFiles,
		CommitBranch: commitBranch,
	}, true)
	return err
}
