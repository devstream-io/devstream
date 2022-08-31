package jenkins

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"
)

type jenkins struct {
	gojenkins.Jenkins
	ctx context.Context
}

type JenkinsAPI interface {
	CreateOrUpdateJob(config, jobName string) (job *gojenkins.Job, created bool, err error)
	ExecuteScript(script string) (string, error)
	GetJob(ctx context.Context, id string, parentIDs ...string) (*gojenkins.Job, error)
	DeleteJob(ctx context.Context, name string) (bool, error)
	InstallPluginsIfNotExists(plugin []string, enableRestart bool) error
	CreateGiltabCredential(id, token string) error
	ConfigCasc(cascScript string) error
}

type setBearerToken struct {
	rt    http.RoundTripper
	token string
}

func (t *setBearerToken) transport() http.RoundTripper {
	if t.rt != nil {
		return t.rt
	}
	return http.DefaultTransport
}

func (t *setBearerToken) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.transport().RoundTrip(r)
}

func NewClient(url string, basicAuthInfo *BasicAuth) (JenkinsAPI, error) {
	url = strings.TrimSuffix(url, "/")
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

	if basicAuthInfo.UsePassWordAuth() {
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
	return jenkinsClient, nil
}
