package jenkins

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"

	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

const (
	domain          = "_"
	credentialScope = "GLOBAL"
)

type jenkins struct {
	gojenkins.Jenkins
	ctx       context.Context
	BasicInfo *JenkinsConfigOption
}

type JenkinsConfigOption struct {
	URL           string
	Namespace     string
	EnableRestart bool
	Offline       bool
	BasicAuth     *BasicAuth
}

type JenkinsAPI interface {
	ExecuteScript(script string) (string, error)
	GetFolderJob(jobName, jobFolder string) (*gojenkins.Job, error)
	DeleteJob(ctx context.Context, name string) (bool, error)
	InstallPluginsIfNotExists(plugin []*JenkinsPlugin) error
	CreateGiltabCredential(id, token string) error
	CreateSSHKeyCredential(id, userName, privateKey string) error
	CreateSecretCredential(id, secretText string) error
	CreatePasswordCredential(id, userName, password string) error
	ConfigCascForRepo(repoCascConfig *RepoCascConfig) error
	ApplyDingTalkBot(config dingtalk.BotConfig) error
	GetBasicInfo() *JenkinsConfigOption
}

func NewClient(configOption *JenkinsConfigOption) (JenkinsAPI, error) {
	url := strings.TrimSuffix(configOption.URL, "/")
	jenkinsClient := &jenkins{}
	jenkinsClient.Server = url

	var basicAuth *gojenkins.BasicAuth
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create a cookie jar")
	}

	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	basicAuthInfo := configOption.BasicAuth
	if basicAuthInfo.usePassWordAuth() {
		basicAuth = &gojenkins.BasicAuth{
			Username: basicAuthInfo.Username, Password: basicAuthInfo.Password,
		}
	} else {
		httpClient.Transport = &setBearerToken{token: basicAuthInfo.Token, rt: httpClient.Transport}
	}

	jenkinsClient.Requester = &gojenkins.Requester{
		Base:      url,
		SslVerify: true,
		Client:    httpClient,
		BasicAuth: basicAuth,
	}
	if _, err := jenkinsClient.Init(context.TODO()); err != nil {
		return nil, errors.Wrap(err, "couldn't init Jenkins API client")
	}

	status, err := jenkinsClient.Poll(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "couldn't poll data from Jenkins API")
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("couldn't poll data from Jenkins API, invalid status code returned: %d", status)
	}
	jenkinsClient.ctx = context.TODO()
	jenkinsClient.BasicInfo = configOption
	return jenkinsClient, nil
}

func (j *jenkins) GetBasicInfo() *JenkinsConfigOption {
	return j.BasicInfo
}

func (o *JenkinsConfigOption) IsOffline() bool {
	return o.Offline
}
