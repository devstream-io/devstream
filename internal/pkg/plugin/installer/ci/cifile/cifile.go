package cifile

import (
	"errors"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/template"
)

type CIFileConfigMap map[string]string
type CIFileVarsMap map[string]interface{}

type CIFileConfig struct {
	Type server.CIServerType `validate:"oneof=jenkins github gitlab" mapstructure:"type"`
	// ConfigLocation represent location of ci config, it can be a remote location or local location
	ConfigLocation downloader.ResourceLocation `validate:"required_without=ConfigContentMap" mapstructure:"configLocation"`
	// Contents respent map of ci fileName to fileContent
	ConfigContentMap CIFileConfigMap `validate:"required_without=ConfigLocation" mapstructure:"configContents"`
	Vars             CIFileVarsMap   `mapstructure:"vars"`
}

// SetContent is used to config ConfigContentMap for ci
func (c *CIFileConfig) SetContent(content string) {
	ciFileName := c.newCIServerClient().CIFilePath()
	if c.ConfigContentMap == nil {
		c.ConfigContentMap = CIFileConfigMap{}
	}
	c.ConfigContentMap[ciFileName] = content
}

func (c *CIFileConfig) SetContentMap(contentMap map[string]string) {
	c.ConfigContentMap = contentMap
}

func (c *CIFileConfig) getGitfileMap() (gitFileMap git.GitFileContentMap, err error) {
	if len(c.ConfigContentMap) == 0 {
		// 1. if ConfigContentMap is empty, get GitFileContentMap from ConfigLocation
		gitFileMap, err = c.getConfigContentFromLocation()
	} else {
		// 2. else render CIFileConfig.ConfigContentMap values
		gitFileMap = make(git.GitFileContentMap)
		ciServerClient := c.newCIServerClient()
		for filePath, content := range c.ConfigContentMap {
			scmFilePath := ciServerClient.GetGitNameFunc()(filePath, "")
			scmFileContent, err := c.renderContent(content)
			if err != nil {
				log.Debugf("ci render git files failed: %+v", err)
				return nil, err
			}
			gitFileMap[scmFilePath] = []byte(scmFileContent)
		}
	}
	if len(gitFileMap) == 0 {
		return nil, errors.New("ci can't get valid ci files")
	}
	return gitFileMap, err
}

func (c *CIFileConfig) newCIServerClient() (ciClient server.CIServerOptions) {
	return server.NewCIServer(c.Type)
}

func (c *CIFileConfig) renderContent(ciFileContent string) (string, error) {
	needRenderContent := len(c.Vars) > 0
	if needRenderContent {
		return template.New().FromContent(ciFileContent).SetDefaultRender(ciTemplateName, c.Vars).Render()
	}
	return ciFileContent, nil

}
func (c *CIFileConfig) getConfigContentFromLocation() (git.GitFileContentMap, error) {
	// 1. get resource
	log.Debugf("ci start to get config files [%s]...", c.ConfigLocation)
	CIFileConfigPath, err := c.ConfigLocation.Download()
	if err != nil {
		return nil, fmt.Errorf("ci get files by %s failed: %w", c.ConfigLocation, err)
	}
	// 2. get ci content map from CIFileConfigPath
	ciClient := c.newCIServerClient()
	return file.GetFileMap(
		CIFileConfigPath, ciClient.FilterCIFilesFunc(),
		ciClient.GetGitNameFunc(), processCIFilesFunc(c.Vars),
	)
}

func (m CIFileVarsMap) Set(k, v string) {
	m[k] = v
}
