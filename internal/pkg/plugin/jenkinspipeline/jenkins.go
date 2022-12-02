package jenkinspipeline

import (
	"errors"
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	defaultAdminSecretName         = "jenkins"
	defaultAdminSecretUserName     = "jenkins-admin-user"
	defaultAdminSecretUserPassword = "jenkins-admin-password"
	jenkinsPasswordEnvKey          = "JENKINS_PASSWORD"
)

type jenkinsOption struct {
	URL           string `mapstructure:"url" validate:"required,url"`
	User          string `mapstructure:"user"`
	Namespace     string `mapstructure:"namespace"`
	EnableRestart bool   `mapstructure:"enableRestart"`
	Offline       bool   `mapstructure:"offline"`
}

func (j *jenkinsOption) newClient() (jenkins.JenkinsAPI, error) {
	auth, err := j.getBasicAuth()
	if err != nil {
		return nil, errors.New("jenkins uesrname and password is required")
	}
	jenkinsConfig := &jenkins.JenkinsConfigOption{
		BasicAuth:     auth,
		URL:           j.URL,
		Namespace:     j.Namespace,
		EnableRestart: j.EnableRestart,
		Offline:       j.Offline,
	}
	return jenkins.NewClient(jenkinsConfig)
}

func (j *jenkinsOption) getBasicAuth() (*jenkins.BasicAuth, error) {
	jenkinsPassword := os.Getenv(jenkinsPasswordEnvKey)
	// 1. check username is set and has env password
	if j.User != "" && jenkinsPassword != "" {
		log.Debugf("jenkins get auth token from env")
		return &jenkins.BasicAuth{
			Username: j.User,
			Password: jenkinsPassword,
		}, nil
	}
	// 2. if not set, get user and password from secret
	k8sClient, err := k8s.NewClient()
	if err != nil {
		secretAuth := getAuthFromSecret(k8sClient, j.Namespace)
		if secretAuth != nil && secretAuth.CheckNameMatch(j.User) {
			log.Debugf("jenkins get auth token from secret")
			return secretAuth, nil
		}
	}
	return nil, errors.New("jenkins uesrname and password is required")
}

func getAuthFromSecret(k8sClient k8s.K8sAPI, namespace string) *jenkins.BasicAuth {
	secret, err := k8sClient.GetSecret(namespace, defaultAdminSecretName)
	if err != nil {
		log.Warnf("jenkins get auth from k8s failed: %+v", err)
		return nil
	}
	user, ok := secret[defaultAdminSecretUserName]
	if !ok {
		return nil
	}
	password, ok := secret[defaultAdminSecretUserPassword]
	if !ok {
		return nil
	}
	return &jenkins.BasicAuth{
		Username: user,
		Password: password,
	}
}
