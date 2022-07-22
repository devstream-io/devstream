package reposcaffolding

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/util"
	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

// DstRepo is the destination repo to push scaffolding project
type DstRepo struct {
	Owner             string `validate:"required_without=Org"`
	Org               string `validate:"required_without=Owner"`
	Repo              string `validate:"required"`
	Branch            string `validate:"required"`
	PathWithNamespace string
	// TODO: used for gitlab
	BaseURL    string `validate:"omitempty,url"`
	Visibility string `validate:"omitempty,oneof=public private internal"`
}

func (d *DstRepo) createLocalRepoPath(workpath string) (string, error) {
	localPath := filepath.Join(workpath, d.Repo)
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		return "", err
	}
	return localPath, nil
}

// this method generate a walker func to render and copy files from srcRepoPath to dstRepoPath
func (d *DstRepo) generateRenderWalker(
	srcRepoPath, dstRepoPath string, renderConfig map[string]interface{},
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
		outputWithRepoName, err := replaceAppNameInPathStr(relativePath, d.Repo)
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
			err = template.RenderForFile(
				"repo-scaffolding", path, outputPath, renderConfig,
			)
			if err != nil {
				return err
			}
		} else {
			if err = util.CopyFile(path, outputPath); err != nil {
				return err
			}
		}
		return nil
	}
}

func (d *DstRepo) createRepoRenderConfig() map[string]interface{} {
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

func (d *DstRepo) createGithubClient(needAuth bool) (*github.Client, error) {
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

// TODO: add gitlab support, temporary comment code
// func (d *DstRepo) createGitlabClient() (*gitlab.Client, error) {
// return gitlab.NewClient(gitlab.WithBaseURL(d.BaseURL))
// }

// func (d *DstRepo) creategitlabOpts() *gitlab.CreateProjectOptions {
// return &gitlab.CreateProjectOptions{
// Name:       d.Repo,
// Branch:     d.Branch,
// Namespace:  d.Org,
// Visibility: d.Visibility,
// }
// }
