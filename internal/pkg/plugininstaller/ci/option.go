package ci

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/template"
)

type CIConfig struct {
	Type      ciRepoType             `validate:"oneof=jenkins github gitlab" mapstructure:"type"`
	LocalPath string                 `validate:"required_without=RemoteURL" mapstructure:"localPath"`
	RemoteURL string                 `validate:"required_without=LocalPath" mapstructure:"remoteURL"`
	Vars      map[string]interface{} `mapstructure:"vars"`
}

type Options struct {
	CIConfig    *CIConfig    `mapstructure:"ci" validate:"required"`
	ProjectRepo *common.Repo `mapstructure:"projectRepo" validate:"required"`
}

func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

// getCIFile will generate ci files by config
func (opt *Options) buildGitMap() (gitMap git.GitFileContentMap, err error) {
	if opt.CIConfig.LocalPath != "" {
		// get content from local if localPath is set
		gitMap, err = opt.CIConfig.getFromLocation(opt.ProjectRepo.Repo)
	} else {
		// get remote content if location is not set
		gitMap, err = opt.CIConfig.getFromURL(opt.ProjectRepo.Repo)
	}
	if len(gitMap) == 0 {
		return nil, errors.New("can't get valid ci files, please check your config")
	}
	return gitMap, err
}

func (c *CIConfig) getFromURL(appName string) (git.GitFileContentMap, error) {
	gitFileMap := make(git.GitFileContentMap)
	content, err := template.New().FromURL(c.RemoteURL).SetDefaultRender(appName, c.Vars).Render()
	if err != nil {
		return nil, err
	}
	fileName := getGitNameFunc(c.Type)("", path.Base(c.RemoteURL))
	gitFileMap[fileName] = []byte(content)
	return gitFileMap, nil
}

func (c *CIConfig) getFromLocation(appName string) (git.GitFileContentMap, error) {
	gitFileMap := make(git.GitFileContentMap)
	info, err := os.Stat(c.LocalPath)
	if err != nil {
		return nil, err
	}
	// process dir
	if info.IsDir() {
		return file.WalkDir(
			c.LocalPath, filterCIFilesFunc(c.Type),
			getGitNameFunc(c.Type), processCIFilesFunc(appName, c.Vars),
		)
	}
	// process file
	gitFilePath := getCIFilePath(c.Type)
	if c.Type == ciGitHubWorkConfigLocation {
		gitFilePath = filepath.Join(gitFilePath, filepath.Base(c.LocalPath))
	}
	content, err := os.ReadFile(c.LocalPath)
	if err != nil {
		return nil, err
	}
	gitFileMap[gitFilePath] = content
	return gitFileMap, nil
}
