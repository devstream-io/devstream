package reposcaffolding

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	transitBranch      = "init-with-devstream"
	appNamePlaceHolder = "_app_name_"
)

type Options struct {
	RepoType        string  `validate:"oneof=gitlab github" mapstructure:"repo_type"`
	SourceRepo      SrcRepo `validate:"required" mapstructure:"source_repo"`
	DestinationRepo DstRepo `validate:"required" mapstructure:"destination_repo"`
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

func (opts *Options) CreateAndRenderLocalRepo(workpath string) error {
	// 1. download template scaffolding repo
	err := opts.SourceRepo.DownloadRepo(workpath)
	if err != nil {
		return err
	}
	// 2. walk iter repo files to render template
	if err := walkLocalRepoPath(workpath, opts); err != nil {
		log.Debugf("create local repo failed walk: %s.", err)
		return err
	}
	return nil
}

// PushToRemoteGitLab push local repo to remote gitlab repo
func (opts *Options) PushToRemoteGitlab(repoPath string) error {
	// TODO: add gitlab push func
	return nil
}

// PushToRemoteGithub push local repo to remote github repo
func (opts *Options) PushToRemoteGithub(repoPath string) error {
	dstRepo := &opts.DestinationRepo
	// 1. init github client
	ghClient, err := dstRepo.createGithubClient(true)
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
	renderConfig := opts.DestinationRepo.createRepoRenderConfig()
	for k, v := range opts.Vars {
		renderConfig[k] = v
	}
	return renderConfig
}
