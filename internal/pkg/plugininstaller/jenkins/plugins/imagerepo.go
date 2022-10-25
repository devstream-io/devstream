package plugins

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/base"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	imageRepoSecretName = "repo-auth"
	defaultImageProject = "library"
)

type ImageRepoJenkinsConfig struct {
	base.ImageRepoStepConfig `mapstructure:",squash"`
}

func (g *ImageRepoJenkinsConfig) getDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{}
}

// imageRepo config will create kubernetes secret for docker auth
// pipeline in jenkins will mount this secret to login image repo
func (g *ImageRepoJenkinsConfig) config(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	log.Info("jenkins plugin imageRepo start config...")
	secretData, err := g.generateDockerAuthSecretData()
	if err != nil {
		return nil, err
	}
	// use k8s client to create secret
	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}
	namespace := jenkinsClient.GetBasicInfo().Namespace
	_, err = client.ApplySecret(imageRepoSecretName, namespace, secretData, nil)
	log.Debugf("jenkins imageRepo secret %s/%s create status: %+v", namespace, imageRepoSecretName, err)
	return nil, err
}

func (g *ImageRepoJenkinsConfig) generateDockerAuthSecretData() (map[string][]byte, error) {
	imageRepoPasswd := os.Getenv("IMAGE_REPO_PASSWORD")
	if imageRepoPasswd == "" {
		return nil, fmt.Errorf("the environment variable IMAGE_REPO_PASSWORD is not set")
	}
	tmpStr := fmt.Sprintf("%s:%s", g.User, imageRepoPasswd)
	authStr := base64.StdEncoding.EncodeToString([]byte(tmpStr))
	log.Debugf("jenkins plugin imageRepo auth string: %s.", authStr)

	configJsonStrTpl := `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
	configJsonStr := fmt.Sprintf(configJsonStrTpl, g.URL, authStr)
	log.Debugf("config.json: %s.", configJsonStr)

	return map[string][]byte{
		"config.json": []byte(configJsonStr),
	}, nil
}

func (p *ImageRepoJenkinsConfig) setRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
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
