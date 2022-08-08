package common

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

// Repo is the repo info of github or gitlab
type Repo struct {
	Owner             string `validate:"required_without=Org"`
	Org               string `validate:"required_without=Owner"`
	Repo              string `validate:"required"`
	Branch            string `validate:"required"`
	PathWithNamespace string
	RepoType          string `validate:"oneof=gitlab github" mapstructure:"repo_type"`
	// This is config for gitlab
	BaseURL    string `validate:"omitempty,url"`
	Visibility string `validate:"omitempty,oneof=public private internal"`
}

// CreateLocalRepoPath create local path for repo
func (d *Repo) CreateLocalRepoPath(workpath string) (string, error) {
	localPath := filepath.Join(workpath, d.Repo)
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		return "", err
	}
	return localPath, nil
}

// Generate is a walker func to render and copy files from srcRepoPath to dstRepoPath
func (d *Repo) GenerateRenderWalker(
	srcRepoPath, dstRepoPath, appNamePlaceHolder string, renderConfig map[string]interface{},
) func(path string, info fs.FileInfo, err error) error {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		relativePath := strings.Replace(path, srcRepoPath, "", 1)
		if strings.Contains(path, ".git/") {
			log.Debugf("Walk: ignore file %s.", "./git/xxx")
			return nil
		}

		if strings.HasSuffix(path, "README.md") {
			log.Debugf("Walk: ignore file %s.", "README.md")
			return nil
		}

		// replace template with appName
		outputWithRepoName, err := replaceAppNameInPathStr(relativePath, appNamePlaceHolder, d.Repo)
		if err != nil {
			log.Debugf("Walk: Replace file name failed %s.", path)
			return err
		}
		outputPath := filepath.Join(dstRepoPath, outputWithRepoName)

		if info.IsDir() {
			log.Debugf("Walk: found dir: %s.", path)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
				return err
			}
			log.Debugf("Walk: new output dir created: %s.", outputPath)
			return nil
		}

		log.Debugf("Walk: found file: %s.", path)

		// if file endswith tpl, render this file, else copy this file directly
		if strings.Contains(path, "tpl") {
			outputPath = strings.TrimSuffix(outputPath, ".tpl")
			return template.RenderForFile("repo-scaffolding", path, outputPath, renderConfig)
		}
		return file.CopyFile(path, outputPath)
	}
}

// CreateRepoRenderConfig will generate template render variables
func (d *Repo) CreateRepoRenderConfig() map[string]interface{} {
	var owner = d.Owner
	if d.Org != "" {
		owner = d.Org
	}

	renderConfigMap := map[string]interface{}{
		"AppName": d.Repo,
		"Repo": map[string]string{
			"Name":  d.Repo,
			"Owner": owner,
		},
	}
	return renderConfigMap
}

// CreateGithubClient build github client connection info
func (d *Repo) CreateGithubClient(needAuth bool) (*github.Client, error) {
	ghOptions := &github.Option{
		Owner:    d.Owner,
		Org:      d.Org,
		Repo:     d.Repo,
		NeedAuth: needAuth,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}
	return ghClient, nil
}

// CreateGitlabClient build gitlab connection info
func (d *Repo) CreateGitlabClient() (*gitlab.Client, error) {
	return gitlab.NewClient(gitlab.WithBaseURL(d.BaseURL))
}

// BuildgitlabOpts build gitlab connection options
func (d *Repo) BuildgitlabOpts() *gitlab.CreateProjectOptions {
	return &gitlab.CreateProjectOptions{
		Name:       d.Repo,
		Branch:     d.Branch,
		Namespace:  d.Org,
		Visibility: d.Visibility,
	}
}

func replaceAppNameInPathStr(filePath, appNamePlaceHolder, appName string) (string, error) {
	if !strings.Contains(filePath, appNamePlaceHolder) {
		return filePath, nil
	}
	newFilePath := regexp.MustCompile(appNamePlaceHolder).ReplaceAllString(filePath, appName)
	log.Debugf("Replace file path place holder. Before: %s, after: %s.", filePath, newFilePath)
	return newFilePath, nil
}
