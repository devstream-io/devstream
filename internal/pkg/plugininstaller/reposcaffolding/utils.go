package reposcaffolding

import (
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func downloadGithubRepo(org, repo, workpath string) error {
	ghOption := &github.Option{
		Org:      org,
		Repo:     repo,
		NeedAuth: false,
		WorkPath: workpath,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return err
	}

	if err = ghClient.DownloadLatestCodeAsZipFile(); err != nil {
		return err
	}

	return nil
}

func walkLocalRepoPath(workpath string, opts *Options) error {
	// 1. get src path and dst path
	srcRepoPath := opts.SourceRepo.getLocalRepoPath(workpath)
	dstOpts := opts.DestinationRepo
	dstRepoPath, err := dstOpts.CreateLocalRepoPath(workpath)
	if err != nil {
		log.Debugf("Walk: create output dir failed: %s", err)
		return err
	}

	// 2. config template render config
	renderConfig := opts.renderTplConfig()

	// 3. create walk func
	walkFunc := dstOpts.GenerateRenderWalker(srcRepoPath, dstRepoPath, appNamePlaceHolder, renderConfig)

	// 4. walk iter srcRepoPath to execuate walk func logic
	if err := filepath.Walk(srcRepoPath, walkFunc); err != nil {
		return err
	}
	return nil
}
