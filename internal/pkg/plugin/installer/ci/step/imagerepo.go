package step

import (
	"encoding/base64"
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	// imageRepoDockerSecretName is used for creating k8s secret
	// and it will be used by jenkins for mount
	imageRepoDockerSecretName = "image-repo-auth"
	// imageRepoSecretName is used for github action secret
	imageRepoSecretName = "IMAGE_REPO_SECRET"
	imageRepoUserName   = "IMAGE_REPO_USER"
)

type ImageRepoStepConfig struct {
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func (g *ImageRepoStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{}
}

// imageRepo config will create kubernetes secret for docker auth
// pipeline in jenkins will mount this secret to login image repo
func (g *ImageRepoStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
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
	_, err = client.ApplySecret(imageRepoDockerSecretName, namespace, secretData, nil)
	log.Debugf("jenkins imageRepo secret %s/%s create status: %+v", namespace, imageRepoSecretName, err)
	return nil, err
}

func (g *ImageRepoStepConfig) ConfigSCM(client scm.ClientOperation) error {
	if g.Password == "" {
		return fmt.Errorf("config field password is not set")
	}

	if err := client.AddRepoSecret(imageRepoUserName, g.User); err != nil {
		return err
	}
	return client.AddRepoSecret(imageRepoSecretName, g.Password)
}

func (g *ImageRepoStepConfig) generateDockerAuthSecretData() (map[string][]byte, error) {
	if g.Password == "" {
		return nil, fmt.Errorf("config field password is not set")
	}
	tmpStr := fmt.Sprintf("%s:%s", g.User, g.Password)
	authStr := base64.StdEncoding.EncodeToString([]byte(tmpStr))
	authURL := g.GetImageRepoURL()
	log.Debugf("jenkins plugin imageRepo auth string: %s.", authStr)

	configJsonStrTpl := `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
	configJsonStr := fmt.Sprintf(configJsonStrTpl, authURL, authStr)
	log.Debugf("config.json: %s.", configJsonStr)

	return map[string][]byte{
		"config.json": []byte(configJsonStr),
	}, nil
}

func (g *ImageRepoStepConfig) GetImageRepoURL() string {
	if g.URL == "" {
		// default use docker image repo
		return "https://index.docker.io/v1/"
	}
	return g.URL
}
