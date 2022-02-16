package golang

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/github"
)

const (
	DefaultWorkPath      = ".github-repo-scaffolding-golang"
	DefaultTemplateRepo  = "dtm-scaffolding-golang"
	DefaultTemplateOwner = "merico-dev"
	TransitBranch        = "init-with-devstream"
	MainBranch           = "main"
)

type Config struct {
	AppName   string
	ImageRepo string
	Repo      Repo
}

type Repo struct {
	Name  string
	Owner string
}

func InitRepoLocalAndPushToRemote(repoPath string, param *Param, ghClient *github.Client) error {
	if err := ghClient.CreateRepo(); err != nil {
		log.Infof("Failed to create repo: %s", err)
		return err
	}
	log.Info("Repo created.")

	if err := WalkLocalRepoPath(repoPath, param, ghClient); err != nil {
		return err
	}

	return MergeCommits(ghClient)
}

func WalkLocalRepoPath(repoPath string, param *Param, ghClient *github.Client) error {
	appName := param.Repo
	if err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Found dir: %s", path)
			return nil
		}
		log.Debugf("Found file: %s", path)

		if strings.Contains(path, ".git/") {
			log.Debugf("Ignore this file -> %s", "./git/xxx")
			return nil
		}

		if strings.HasSuffix(path, "README.md") {
			log.Debugf("Ignore this file -> %s", "README.md")
			return nil
		}

		pathForGithub, err := genPathForGithub(path)
		if err != nil {
			return err
		}

		var content []byte
		if strings.Contains(path, "tpl") {
			content, err = Render(path, param)
			if err != nil {
				return err
			}
		} else {
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		}
		log.Debugf("Content size: %d", len(content))

		if newPathForGithub, err := replaceAppNameInPathStr(pathForGithub, appName); err != nil {
			return err
		} else {
			// the main branch needs a initial commit
			if strings.Contains(newPathForGithub, "gitignore") {
				err := ghClient.CreateFile(content, strings.TrimSuffix(newPathForGithub, ".tpl"), MainBranch)
				if err != nil {
					log.Debugf("Failed to add the .gitignore file.")
					return err
				}
				log.Debugf("Added the .gitignore file.")
				return ghClient.NewBranch(MainBranch, TransitBranch)
			}
			return ghClient.CreateFile(content, strings.TrimSuffix(newPathForGithub, ".tpl"), TransitBranch)
		}
	}); err != nil {
		return err
	}

	return nil
}

func MergeCommits(ghClient *github.Client) error {
	number, err := ghClient.NewPullRequest(TransitBranch, MainBranch)
	if err != nil {
		return err
	}

	return ghClient.MergePullRequest(number, github.MergeMethodSquash)
}

func Render(filePath string, param *Param) ([]byte, error) {
	config := Config{
		AppName:   param.Repo,
		ImageRepo: param.ImageRepo,
		Repo: Repo{
			Name:  param.Repo,
			Owner: param.Owner,
		},
	}
	log.Debugf("filePath: %s", filePath)
	log.Debugf("Config %v", config)

	textBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	textStr := string(textBytes)

	tpl := template.New("github-repo-scaffolding-golang").Delims("[[", "]]")
	parsed, err := tpl.Parse(textStr)
	if err != nil {
		log.Debugf("template parse file failed: %s", err)
		return nil, err
	}

	var buf bytes.Buffer
	if err = parsed.Execute(&buf, config); err != nil {
		log.Debugf("template execute failed: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func genPathForGithub(filePath string) (string, error) {
	splitStrs := strings.SplitN(filePath, "/", 3)
	if len(splitStrs) != 3 {
		return "", fmt.Errorf("unknown format: %s", filePath)
	}
	retStr := splitStrs[2]
	log.Debugf("Path for github: %s", retStr)
	return retStr, nil
}

func replaceAppNameInPathStr(filePath, appName string) (string, error) {
	log.Debugf("Got filePath %s", filePath)

	pet := "_app_name_"
	reg, err := regexp.Compile(pet)
	if err != nil {
		return "", err
	}
	newFilePath := reg.ReplaceAllString(filePath, appName)

	log.Debugf("New filePath: %s", newFilePath)

	return newFilePath, nil
}

func buildState(param *Param) map[string]interface{} {
	res := make(map[string]interface{})
	res["owner"] = param.Owner
	res["repoName"] = param.Repo
	return res
}
