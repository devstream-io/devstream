package jenkins

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

const (
	defaultNameSpace               = "jenkins"
	defaultAdminSecretName         = "jenkins"
	defautlAdminSecretUserName     = "jenkins-admin-user"
	defautlAdminSecretUserPassword = "jenkins-admin-password"
)

// SetJobDefaultConfig config default fields for usage
func SetJobDefaultConfig(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}
	// config default values
	projectRepo, err := common.NewRepoFromURL(opts.ProjectURL, opts.ProjectBranch)
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo
	if opts.JobName == "" {
		opts.JobName = projectRepo.Repo
	}
	opts.CIConfig = buildCIConfig(opts.JenkinsfilePath)
	basicAuth, err := buildAdminToken(opts.JenkinsUser)
	if err != nil {
		return nil, err
	}
	opts.BasicAuth = basicAuth
	opts.SecretToken = generateRandomSecretToken()
	return opts.encode()
}

func ValidateJobConfig(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}
	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}
	if opts.ProjectRepo.RepoType == "github" {
		return nil, errors.New("jenkins job not support github for now")
	}
	return options, nil
}

func buildAdminToken(userName string) (*jenkins.BasicAuth, error) {
	// 1. check username is set and has env password
	jenkinsPassword := os.Getenv("JENKINS_PASSWORD")
	if userName != "" && jenkinsPassword != "" {
		return &jenkins.BasicAuth{
			Username: userName,
			Password: jenkinsPassword,
		}, nil
	}
	// 2. if not set, get user and password from secret
	secretAuth := getAuthFromSecret()
	if secretAuth != nil && secretAuth.IsNameMatch(userName) {
		return secretAuth, nil
	}
	return nil, errors.New("jenkins uesrname and password is required")
}

func getAuthFromSecret() *jenkins.BasicAuth {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return nil
	}
	secret, err := k8sClient.GetSecret(defaultNameSpace, defaultAdminSecretName)
	if err != nil {
		return nil
	}
	user, ok := secret[defautlAdminSecretUserName]
	if !ok {
		return nil
	}
	password, ok := secret[defautlAdminSecretUserPassword]
	if !ok {
		return nil
	}
	return &jenkins.BasicAuth{
		Username: user,
		Password: password,
	}
}

func generateRandomSecretToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:32]
}

func buildCIConfig(jenkinsFilePath string) *ci.CIConfig {
	ciConfig := &ci.CIConfig{
		Type: "jenkins",
	}
	// config CIConfig
	url, err := url.ParseRequestURI(jenkinsFilePath)
	// if path is url, download from remote
	if err != nil || url.Host == "" {
		ciConfig.LocalPath = jenkinsFilePath
	} else {
		ciConfig.RemoteURL = jenkinsFilePath
	}
	return ciConfig
}

func SetHarborAuth(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}
	namespace := opts.JenkinsNamespace
	harborURL := opts.HarborURL
	harborUser := opts.HarborUser

	harborPasswd := os.Getenv("HARBOR_PASSWORD")
	if harborPasswd == "" {
		return nil, fmt.Errorf("the environment variable HARBOR_PASSWORD is not set")
	}

	return options, createDockerSecret(namespace, harborURL, harborUser, harborPasswd)
}

func createDockerSecret(namespace, url, username, password string) error {
	tmpStr := fmt.Sprintf("%s:%s", username, password)
	authStr := base64.StdEncoding.EncodeToString([]byte(tmpStr))
	log.Debugf("Auth string: %s.", authStr)

	configJsonStrTpl := `{
"auths": {
   "%s": {
     "auth": "%s"
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
	_, err = client.ApplySecret("docker-config", namespace, data, nil)
	if err != nil {
		return err
	}
	log.Infof("Secret %s/%s has been created.", namespace, "docker-config")
	return nil
}
