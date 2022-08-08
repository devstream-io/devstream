package reposcaffolding

import (
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	transitBranch    = "init-with-devstream"
	defaultCommitMsg = "init with devstream"
)

type Options struct {
	SourceRepo      *SrcRepo     `validate:"required" mapstructure:"source_repo"`
	DestinationRepo *common.Repo `validate:"required" mapstructure:"destination_repo"`
	Vars            map[string]interface{}
}

func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (opts *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

// CreateAndRenderLocalRepo will download repo from source repo and render it locally
func (opts *Options) CreateAndRenderLocalRepo() (string, error) {
	// 1. get download url
	githubCodeZipDownloadURL, err := opts.SourceRepo.getDownloadURL()
	if err != nil {
		log.Debugf("reposcaffolding get download url failed: %s", err)
		return "", err
	}
	// 2. download zip file and unzip this file then render folders
	projectDir, err := file.NewTemplate().FromRemote(githubCodeZipDownloadURL).UnzipFile().RenderRepoDIr(
		opts.DestinationRepo.Repo, opts.renderTplConfig(),
	).Run()
	if err != nil {
		log.Debugf("reposcaffolding process file error: %s", err)
		return "", err
	}
	// 3. join download path and repo name to get repo path
	repoDirName := opts.SourceRepo.getRepoName()
	return filepath.Join(projectDir, repoDirName), nil
}

// PushToRemoteGitLab push local repo to remote gitlab repo
func (opts *Options) PushToRemoteGitlab(repoPath string) error {
	dstRepo := opts.DestinationRepo
	// 1. init gitlab client
	c, err := dstRepo.CreateGitlabClient()
	if err != nil {
		log.Debugf("Gitlab push: init gitlab client failed %s", err)
		return err
	}

	// 2. create the project
	if err := c.CreateProject(dstRepo.BuildgitlabOpts()); err != nil {
		log.Errorf("Failed to create repo: %s.", err)
		return err
	}

	// if encounter error, delete repo
	var needRollBack bool
	defer func() {
		if !needRollBack {
			return
		}
		// need to clean the repo created when retErr != nil
		if err := c.DeleteProject(dstRepo.PathWithNamespace); err != nil {
			log.Errorf("Failed to delete the repo %s: %s.", dstRepo.PathWithNamespace, err)
		}
	}()

	needRollBack, err = c.PushLocalPathToBranch(
		repoPath, dstRepo.Branch, dstRepo.PathWithNamespace, defaultCommitMsg,
	)
	if err != nil {
		log.Errorf("Failed to push to remote: %s.", err)
		return err
	}
	return nil
}

// PushToRemoteGithub push local repo to remote github repo
func (opts *Options) PushToRemoteGithub(repoPath string) error {
	dstRepo := opts.DestinationRepo
	// 1. init github client
	ghClient, err := dstRepo.CreateGithubClient(true)
	if err != nil {
		log.Debugf("Github push: init github client failed %s", err)
		return err
	}

	// 2. init repo
	if err := ghClient.InitRepo(dstRepo.Branch); err != nil {
		return err
	}

	// if encounter rollout error, delete repo
	var needRollBack bool
	defer func() {
		if !needRollBack {
			return
		}
		// need to clean the repo created when retErr != nil
		if err := ghClient.DeleteRepo(); err != nil {
			log.Errorf("Failed to delete the repo %s: %s.", dstRepo.Repo, err)
		}
	}()

	// 3. push local path to repo
	needRollBack, err = ghClient.PushLocalPathToBranch(transitBranch, dstRepo.Branch, repoPath)
	if err != nil {
		return err
	}
	return nil
}

func (opts *Options) renderTplConfig() map[string]interface{} {
	renderConfig := opts.DestinationRepo.CreateRepoRenderConfig()
	for k, v := range opts.Vars {
		renderConfig[k] = v
	}
	return renderConfig
}
