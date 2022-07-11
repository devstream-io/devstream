package jenkinspipelinekubernetes

import (
	_ "embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/parnurzeal/gorequest"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	jenkinsCredentialID   = "credential-devstream-io-jenkins-pipeline-kubernetes"
	jenkinsCredentialDesc = "Jenkins Pipeline secret, created by devstream-io/jenkins-pipeline-kubernetes"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// always try to create credential, there will be no error if it already exists
	if err := CreateCredentialUsernamePassword(&opts); err != nil {
		return nil, err
	}

	// todo: check if the job already exists
	if err := CreateItem(&opts); err != nil {
		return nil, fmt.Errorf("failed to create job: %s", err)
	}

	return (&resource{}).toMap(), nil
}

//func ifCredentialExists(options *Options) bool {
//	// todo
//	return false
//}

func CreateCredentialSecretText(opts *Options) error {

	accessURL := opts.GetJenkinsAccessURL()

	// todo: use gorequest to do the request
	cmdGetCrumb := fmt.Sprintf(`CRUMB=$(curl -s '%s/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,":",//crumb)')`, accessURL)
	cmdCreateCredential := fmt.Sprintf(`
curl -H $CRUMB -X POST '%s/credentials/store/system/domain/_/createCredentials' \
--data-urlencode 'json={
  "": "0",
   "credentials": {
	   "scope": "GLOBAL",
	   "id": "%s",
	   "secret": "%s",
	   "description": "%s",
	   "$class": "org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl"
  }
}'`, accessURL, jenkinsCredentialID, opts.GitHubToken, jenkinsCredentialDesc)

	cmd := exec.Command("sh", "-c", cmdGetCrumb+" && "+cmdCreateCredential)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	return nil
}

func CreateCredentialUsernamePassword(opts *Options) error {

	accessURL := opts.GetJenkinsAccessURL()

	// todo: use gorequest to do the request
	cmdGetCrumb := fmt.Sprintf(`CRUMB=$(curl -s '%s/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,":",//crumb)')`, accessURL)
	cmdCreateCredential := fmt.Sprintf(`
curl -H $CRUMB -X POST '%s/credentials/store/system/domain/_/createCredentials' \
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
}'`, accessURL, jenkinsCredentialID, opts.GitHubToken, jenkinsCredentialDesc)

	cmd := exec.Command("sh", "-c", cmdGetCrumb+" && "+cmdCreateCredential)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create credential secret: %s", err)
	}

	return nil
}

//go:embed job-template.xml
var jobTemplate string

func CreateItem(opts *Options) error {

	jobXml := renderJobXml(jobTemplate, opts)

	// create job
	request := gorequest.New()
	resp, body, errs := request.Post(opts.GetJenkinsAccessURL()+"/createItem").
		Set("Content-Type", "application/xml").
		Query("name=" + opts.JenkinsJobName).Send(jobXml).End()

	if len(errs) != 0 {
		return fmt.Errorf("failed to create job: %s", errs)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("failed to create job, here is response: %s", body)
	}

	return nil
}

// todo: unit test
func renderJobXml(jobTemplate string, opts *Options) string {
	// note: maybe it is better to use html/template to generate the job template,
	// but it is complex and this is the simplest way to do it
	jobXml := strings.Replace(jobTemplate, "{{.GitHubRepoURL}}", opts.GitHubRepoURL, 1)
	jobXml = strings.Replace(jobXml, "{{.CredentialsID}}", jenkinsCredentialID, 1)

	fmt.Println(jobXml)

	return jobXml
}
