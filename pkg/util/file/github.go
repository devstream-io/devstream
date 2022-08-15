package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

// DownloadAndUnzipAndRenderGitHubRepo will download/unzip/render git repo files and return unzip files dir path
func DownloadAndUnzipAndRenderGitHubRepo(url, templateName string, vars map[string]interface{}) (string, error) {
	zipFile, err := getFileFromURL(url)
	if err != nil {
		return "", err
	}

	unZipDir, err := unzip(zipFile)
	if err != nil {
		return "", err
	}

	dstPath, err := renderGitRepoDir(templateName, unZipDir, vars)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

// renderGitRepoDir will render files in srcPath to dstPath, It will render two things
// 1. replace filename with appNamePlaceHolder to templateName
// 2. render files in srcPath and output to dstPath
func renderGitRepoDir(templateName, srcPath string, vars map[string]interface{}) (string, error) {
	// 1. create temp dir for destination
	dstPath, err := os.MkdirTemp("", defaultTempName)
	if err != nil {
		return "", err
	}
	if err := filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		relativePath := strings.Replace(path, srcPath, "", 1)
		if strings.Contains(relativePath, ".git/") {
			log.Debugf("Walk: ignore file %s.", "./git/xxx")
			return nil
		}

		if strings.HasSuffix(relativePath, "README.md") {
			log.Debugf("Walk: ignore file %s.", "README.md")
			return nil
		}

		// replace template with appName
		outputWithRepoName, err := replaceAppNameInPathStr(relativePath, appNamePlaceHolder, templateName)
		if err != nil {
			log.Debugf("Walk: Replace file name failed %s.", path)
			return err
		}
		outputPath := filepath.Join(dstPath, outputWithRepoName)

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
			return template.RenderForFile(templateName, path, outputPath, vars)
		}
		return CopyFile(path, outputPath)
	}); err != nil {
		return "", err
	}
	return dstPath, nil
}
