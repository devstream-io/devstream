package jenkins

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

const (
	defaultNameSpace               = "jenkins"
	defaultAdminSecretName         = "jenkins"
	defautlAdminSecretUserName     = "jenkins-admin-user"
	defautlAdminSecretUserPassword = "jenkins-admin-password"
	defaultImageProject            = "library"
)

// SetJobDefaultConfig config default fields for usage
func SetJobDefaultConfig(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	// config scm and projectRepo values
	projectRepo, err := common.NewRepoFromURL(opts.SCM.Type, opts.SCM.APIURL, opts.SCM.CloneURL, opts.SCM.Branch)
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo

	if opts.SCM.Branch == "" {
		opts.SCM.Branch = projectRepo.Branch
	}
	switch projectRepo.RepoType {
	case "gitlab":
		// set secretToken for gitlab webhook
		sshKey := os.Getenv("GITLAB_SSHKEY")
		if sshKey != "" && opts.SCM.SSHprivateKey == "" {
			opts.SCM.SSHprivateKey = sshKey
		}
	case "github":
		log.Debug("jenkins pipeline repo type is github")
	}
	opts.SecretToken = generateRandomSecretToken()

	// config jenkins and job
	if opts.Jenkins.Namespace == "" {
		opts.Jenkins.Namespace = "jenkins"
	}

	if opts.Pipeline.JobName == "" {
		opts.Pipeline.JobName = projectRepo.Repo
	}

	// config ci related values
	opts.buildCIConfig()

	// config jenkins connection info
	basicAuth, err := buildAdminToken(opts.Jenkins.User)
	if err != nil {
		return nil, err
	}
	opts.BasicAuth = basicAuth
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

	switch opts.ProjectRepo.RepoType {
	case "gitlab":
		if os.Getenv(gitlab.TokenEnvKey) == "" {
			return nil, fmt.Errorf("jenkins-pipeline github should set env %s", gitlab.TokenEnvKey)
		}
	case "github":
		if os.Getenv(github.TokenEnvKey) == "" {
			return nil, fmt.Errorf("jenkins-pipeline github should set env %s", github.TokenEnvKey)
		}
	default:
		return nil, fmt.Errorf("jenkins-pipeline doesn't support repo type %s", opts.ProjectRepo.RepoType)
	}

	// check jenkins job name
	if strings.Contains(opts.Pipeline.JobName, "/") {
		strs := strings.Split(opts.Pipeline.JobName, "/")
		if len(strs) != 2 || len(strs[0]) == 0 || len(strs[1]) == 0 {
			return nil, fmt.Errorf("jobName illegal: %s", opts.Pipeline.JobName)
		}
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

func SetHarborAuth(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	namespace := opts.Jenkins.Namespace
	imageRepoUrl := opts.Pipeline.getImageHost()
	imageRepoUser := opts.Pipeline.ImageRepo.User

	imageRepoPasswd := os.Getenv("IMAGE_REPO_PASSWORD")
	if imageRepoPasswd == "" {
		return nil, fmt.Errorf("the environment variable IMAGE_REPO_PASSWORD is not set")
	}

	return options, createDockerSecret(namespace, imageRepoUrl, imageRepoUser, imageRepoPasswd)
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
