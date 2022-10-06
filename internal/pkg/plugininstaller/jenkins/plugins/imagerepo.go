package plugins

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	defaultImageProject = "library"
	imageRepoSecretName = "docker-config"
)

type ImageRepoJenkinsConfig struct {
	AuthNamespace string `mapstructure:"secretNamespace"`
	URL           string `mapstructure:"url" validate:"url"`
	User          string `mapstructure:"user"`
}

func (g *ImageRepoJenkinsConfig) GetDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{}
}

func (g *ImageRepoJenkinsConfig) PreConfig(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	log.Info("jenkins plugin imageRepo start config...")
	imageRepoPasswd := os.Getenv("IMAGE_REPO_PASSWORD")
	if imageRepoPasswd == "" {
		return nil, fmt.Errorf("the environment variable IMAGE_REPO_PASSWORD is not set")
	}

	return nil, g.createDockerSecret(g.AuthNamespace, g.URL, g.User, imageRepoPasswd)
}

func (g *ImageRepoJenkinsConfig) createDockerSecret(namespace, url, username, password string) error {
	tmpStr := fmt.Sprintf("%s:%s", username, password)
	authStr := base64.StdEncoding.EncodeToString([]byte(tmpStr))
	log.Debugf("Auth string: %s.", authStr)

	configJsonStrTpl := `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
	configJsonStr := fmt.Sprintf(configJsonStrTpl, url, authStr)
	log.Debugf("config.json: %s.", configJsonStr)

	// create secret in k8s
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}

	data := map[string][]byte{
		"config.json": []byte(configJsonStr),
	}
	_, err = client.ApplySecret(imageRepoSecretName, namespace, data, nil)
	if err != nil {
		return err
	}
	log.Infof("Secret %s/%s has been created.", namespace, imageRepoSecretName)
	return nil
}

func (p *ImageRepoJenkinsConfig) UpdateJenkinsFileRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
	var host string
	imageURL, err := url.ParseRequestURI(p.URL)
	if err != nil {
		host = p.URL
	} else {
		host = imageURL.Host
	}
	repositoryURL := path.Join(host, defaultImageProject)
	vars.ImageRepositoryURL = repositoryURL
	vars.ImageAuthSecretName = imageRepoSecretName
}
