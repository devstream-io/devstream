package jenkins

import (
	"errors"
	"net/url"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

const (
	defaultNameSpace               = "jenkins"
	defaultAdminSecretName         = "jenkins"
	defautlAdminSecretUserName     = "jenkins-admin-user"
	defautlAdminSecretUserPassword = "jenkins-admin-password"
)

func SetJobDefaultConfig(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}
	// config projectRepo
	projectRepo, err := common.NewRepoFromURL(opts.ProjectURL, opts.ProjectBranch)
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo
	if opts.JobName == "" {
		opts.JobName = projectRepo.Repo
	}

	ciConfig := &ci.CIConfig{
		Type: "jenkins",
	}
	// config CIConfig
	_, err = url.ParseRequestURI(opts.JenkinsfilePath)
	// if path is url, download from remote
	if err == nil {
		ciConfig.RemoteURL = opts.JenkinsfilePath
	} else {
		ciConfig.LocalPath = opts.JenkinsfilePath
	}
	opts.CIConfig = ciConfig
	err = setAdminToken(opts)
	if err != nil {
		return nil, err
	}
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
	if token := opts.ProjectRepo.GetRepoToken(); token == "" {
		return nil, errors.New("git repo token is required")
	}
	if opts.ProjectRepo.RepoType == "github" {
		return nil, errors.New("jenkins job not support github for now")
	}
	return options, nil
}

func setAdminToken(opts *JobOptions) error {
	// 1. check username is set and has env password
	userName := opts.JenkinsUser
	jenkinsPassword := os.Getenv("JENKINS_PASSWORD")
	if userName != "" && jenkinsPassword != "" {
		opts.BasicAuth = &jenkins.BasicAuth{
			Username: userName,
			Password: jenkinsPassword,
		}
		return nil
	}
	// 2. if not set, get user and password from secret
	secretAuth := getAuthFromSecret()
	if secretAuth != nil && secretAuth.IsNameMatch(userName) {
		opts.BasicAuth = secretAuth
		return nil
	}
	return errors.New("jenkins uesrname and password is required")
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
