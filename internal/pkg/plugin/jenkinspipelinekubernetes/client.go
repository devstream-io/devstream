package jenkinspipelinekubernetes

import (
	_ "embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/parnurzeal/gorequest"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// ClientInf represents the client abstraction for jenkins
type ClientInf interface {
	CreateCredentialSecretText() error
	CreateCredentialUsernamePassword() error
	GetCrumb() (string, error)
	GetCrumbHeader() (headerKey, headerValue string, err error)
	CreateItem(jobXmlContent string) error
}

// Client is the client for jenkins
type Client struct {
	Opts *Options
}

func NewClient(options *Options) *Client {
	return &Client{
		Opts: options,
	}
}

//func (c *Client) ifCredentialExists() bool {
//	// todo
//	return false
//}

// CreateCredentialSecretText creates a credential in the type of "Secret text"
func (c *Client) CreateCredentialSecretText() error {

	accessURL := c.Opts.GetJenkinsAccessURL()
	crumb, err := c.GetCrumb()
	if err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	// TODO(aFlyBird0): use gorequest to do the request
	cmdCreateCredential := fmt.Sprintf(`
curl -H %s -X POST '%s/credentials/store/system/domain/_/createCredentials' \
--data-urlencode 'json={
  "": "0",
   "credentials": {
	   "scope": "GLOBAL",
	   "id": "%s",
	   "secret": "%s",
	   "description": "%s",
	   "$class": "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
  }
}'`, crumb, accessURL, jenkinsCredentialID, c.Opts.GitHubToken, jenkinsCredentialDesc)

	cmd := exec.Command("sh", "-c", cmdCreateCredential)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	return nil
}

// CreateCredentialUsernamePassword creates a credential in the type of "Username with password"
func (c *Client) CreateCredentialUsernamePassword() error {

	accessURL := c.Opts.GetJenkinsAccessURL()
	crumb, err := c.GetCrumb()
	if err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	// TODO(aFlyBird0): use gorequest to do the request
	cmdCreateCredential := fmt.Sprintf(`
curl -H %s -X POST '%s/credentials/store/system/domain/_/createCredentials' \
--data-urlencode 'json={
	"": "0",
	"credentials": {
		"scope": "GLOBAL",
		"id": "%s",
		"username": "foo-useless-username",
		"password": "%s",
		"description": "%s",
		"$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
  }
}'`, crumb, accessURL, jenkinsCredentialID, c.Opts.GitHubToken, jenkinsCredentialDesc)

	cmd := exec.Command("sh", "-c", cmdCreateCredential)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	return nil
}

// GetCrumb returns the crumb for jenkins,
// jenkins uses crumb to prevent CSRF(cross-site request forgery),
// format: "Jenkins-Crumb:a70290b6423777f0a4c771d4805637ac36d5fd78336a20d48d72167ef5f13b9a"
// ref: https://www.jenkins.io/doc/upgrade-guide/2.176/#upgrading-to-jenkins-lts-2-176-3
// ref: https://stackoverflow.com/questions/44711696/jenkins-403-no-valid-crumb-was-included-in-the-request
func (c *Client) GetCrumb() (string, error) {
	request := gorequest.New()
	getCrumbURL := c.Opts.GetJenkinsAccessURL() + `/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,":",//crumb)`
	resp, body, errs := request.Get(getCrumbURL).End()
	log.Debugf("GetCrumb url: %s", getCrumbURL)
	if len(errs) != 0 {
		return "", fmt.Errorf("failed to get crumb: %s", errs)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get crumb, here is response: %s", body)
	}

	return strings.TrimSpace(body), nil
}

// GetCrumbHeader behaves like GetCrumb, but it returns the header key and value
func (c *Client) GetCrumbHeader() (headerKey, headerValue string, err error) {
	// crumb format: "Jenkins-Crumb:a70290b6423777f0a4c771d4805637ac36d5fd78336a20d48d72167ef5f13b9a"
	crumb, err := c.GetCrumb()
	if err != nil {
		return "", "", err
	}
	crumbMap := strings.Split(crumb, ":")
	if len(crumbMap) != 2 {
		return "", "", fmt.Errorf("failed to get crumb, here is response: %s", crumb)
	}
	return crumbMap[0], crumbMap[1], nil
}

//go:embed job-template.xml
var jobTemplate string

// CreateItem creates a job in jenkins with the given job xml
func (c *Client) CreateItem(jobXmlContent string) error {
	request := gorequest.New()
	resp, body, errs := request.Post(c.Opts.GetJenkinsAccessURL()+"/createItem").
		Set("Content-Type", "application/xml").
		Query("name=" + c.Opts.JenkinsJobName).Send(jobXmlContent).End()

	if len(errs) != 0 {
		return fmt.Errorf("failed to create job: %s", errs)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("failed to create job, here is response: %s", body)
	}

	return nil
}
