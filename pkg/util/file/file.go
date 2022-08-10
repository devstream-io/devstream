package file

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultTempName    = "pkg-util-file-create_"
	appNamePlaceHolder = "_app_name_"
)

type fileProcesser func(string) (string, error)
type renderProcesser func(string, string, map[string]interface{}) (string, error)

// TemplateConfig contains all file template getter, process and render info
type TemplateConfig struct {
	info         string
	templateName string
	vars         map[string]interface{}
	getter       fileProcesser
	processor    fileProcesser
	render       renderProcesser
	processErr   error
}

func NewTemplate() *TemplateConfig {
	return &TemplateConfig{}
}

// FromLocal will check if local file exist
func (c *TemplateConfig) FromLocal(path string) *TemplateConfig {
	ex, err := os.Executable()
	if err != nil {
		c.processErr = err
	} else {
		exPath := filepath.Dir(ex)
		c.info = filepath.Join(exPath, path)
	}
	c.getter = getFileFromLocal
	return c
}

// FromRemote will create a file from remote url
func (c *TemplateConfig) FromRemote(url string) *TemplateConfig {
	c.getter = getFileFromURL
	c.info = url
	return c
}

// FromContent will create a file for content
func (c *TemplateConfig) FromContent(content string) *TemplateConfig {
	c.getter = getFileFromContent
	c.info = content
	return c
}

func (c *TemplateConfig) UnzipFile() *TemplateConfig {
	c.processor = unZipFileProcesser
	return c
}

func (c *TemplateConfig) RenderRepoDIr(templateName string, vars map[string]interface{}) *TemplateConfig {
	c.render = renderGitRepoDir
	c.templateName = templateName
	c.vars = vars
	return c
}

// RenderFile will create render config
func (c *TemplateConfig) RenderFile(templateName string, vars map[string]interface{}) *TemplateConfig {
	c.render = renderFile
	c.templateName = templateName
	c.vars = vars
	return c
}

func (c *TemplateConfig) Run() (string, error) {
	// check if has error before
	if c.processErr != nil {
		return "", c.processErr
	}

	// check if info is empty
	if c.info == "" {
		return "", fmt.Errorf("file util: content is not setted")
	}
	var (
		outPutName string
		err        error
	)
	// 1. run Getter func to get file
	outPutName, err = c.getter(c.info)
	if err != nil {
		return "", err
	}
	// 2. if need file process, run processer
	if c.processor != nil {
		outPutName, err = c.processor(outPutName)
		if err != nil {
			return "", err
		}
	}
	// 3. if need render, render func
	if c.render != nil {
		return c.render(c.templateName, outPutName, c.vars)
	}
	return outPutName, nil
}

// CopyFile will copy file content from src to dst
func CopyFile(src, dest string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytesRead, 0644)
}
