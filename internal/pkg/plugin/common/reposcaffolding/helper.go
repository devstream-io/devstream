package reposcaffolding

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/devstream-io/devstream/pkg/util/github"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func walkLocalRepoPath(workpath string, opts *Options, ghClient *github.Client) error {
	repoPath := filepath.Join(workpath, fmt.Sprintf("%s-%s", TemplateRepo, MainBranch))
	appName := opts.Repo
	outputDir := fmt.Sprintf("%s/%s", workpath, opts.Repo)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	if err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Walk: found dir: %s.", path)
			oldDir, err := replaceAppNameInPathStr(path, appName)
			if err != nil {
				return err
			}
			oldDirSections := strings.Split(oldDir, "/")
			newDir := fmt.Sprintf("%s/%s", outputDir, strings.Join(oldDirSections[2:], "/"))
			if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
				return err
			}
			log.Debugf("Walk: new output dir created: %s.", newDir)
			return nil
		}

		if strings.Contains(path, ".git/") {
			log.Debugf("Walk: ignore file %s.", "./git/xxx")
			return nil
		}

		if strings.HasSuffix(path, "README.md") {
			log.Debugf("Walk: ignore file %s.", "README.md")
			return nil
		}

		log.Debugf("Walk: found file: %s.", path)

		outputPath := fmt.Sprintf("%s/%s", outputDir, strings.Join(strings.Split(path, "/")[2:], "/"))
		if outputPath, err = replaceAppNameInPathStr(outputPath, appName); err != nil {
			return err
		}
		outputPath = strings.TrimSuffix(outputPath, ".tpl")

		if strings.Contains(path, "tpl") {
			err = render(path, outputPath, opts)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func render(filePath, output string, opts *Options) error {
	var owner = opts.Owner
	if opts.Org != "" {
		owner = opts.Org
	}

	config := Config{
		AppName:   opts.Repo,
		ImageRepo: opts.ImageRepo,
		Repo: Repo{
			Name:  opts.Repo,
			Owner: owner,
		},
	}

	log.Debugf("Render filePath: %s.", filePath)
	log.Debugf("Render config: %v.", config)
	log.Debugf("Render output: %s.", output)

	textBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	textStr := string(textBytes)

	tpl := template.New("github-repo-scaffolding-golang").Delims("[[", "]]")
	parsed, err := tpl.Parse(textStr)
	if err != nil {
		log.Debugf("Template parse file failed: %s.", err)
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		log.Error("create file: ", err)
		return err
	}
	defer f.Close()
	if err = parsed.Execute(f, config); err != nil {
		log.Debugf("Template execution failed: %s.", err)
		return err
	}

	return nil
}

func replaceAppNameInPathStr(filePath, appName string) (string, error) {
	if !strings.Contains(filePath, AppNamePlaceHolder) {
		return filePath, nil
	}

	reg, err := regexp.Compile(AppNamePlaceHolder)
	if err != nil {
		return "", err
	}
	newFilePath := reg.ReplaceAllString(filePath, appName)

	log.Debugf("Replace file path place holder. Before: %s, after: %s.", filePath, newFilePath)

	return newFilePath, nil
}
