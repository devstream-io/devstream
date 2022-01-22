package golang

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/util/github"
)

const (
	DefaultWorkPath = ".github-repo-scaffolding-golang"
	// TODO(daniel-hutao): Make it configurable
	DefaultTemplateTag   = "v0.0.1"
	DefaultAssetName     = "dtm-scaffolding-golang-v0.0.1.tar.gz"
	DefaultTemplateRepo  = "dtm-scaffolding-golang"
	DefaultTemplateOwner = "merico-dev"
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
		log.Errorf("Failed to create repo: %s", err)
		return err
	}

	err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Found dir: %s", path)
			return nil
		}

		if strings.Contains(path, ".git") {
			return nil
		}

		if strings.HasSuffix(path, "README.md") {
			return nil
		}

		log.Debugf("Found file: %s", path)

		newPath, err := genPathForGithub(path)
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

		return ghClient.CreateFile(content, strings.TrimSuffix(newPath, ".tpl"))
	})

	if err != nil {
		return err
	}
	return nil
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
	splitStrs := strings.SplitN(filePath, "/", 2)
	if len(splitStrs) != 2 {
		return "", fmt.Errorf("unknown format: %s", filePath)
	}
	return splitStrs[1], nil
}
