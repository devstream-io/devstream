package jenkins

import (
	"fmt"
	"testing"

	"github.com/parnurzeal/gorequest"
)

const jenkinsCredentialID = "credential-jenkins-pipeline-kubernetes-by-devstream"

func initJenkins() *Jenkins {
	jenkinsURL := "http://localhost:32001/"
	username := "admin"
	password := "B1OvhHMnPPxXz4kFwODYIh"
	jenkins, err := NewJenkins(jenkinsURL, username, password)
	if err != nil {
		panic(err)
	}
	return jenkins
}

func TestGetCredentialsUsername(t *testing.T) {
	cd, err := initJenkins().GetCredentialsUsername(jenkinsCredentialID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Printf("%+v\n", cd)
}

func TestDeleteJob(t *testing.T) {
	req := gorequest.New().Post("http://localhost:32001/job/jenkins-plugin-test/doDelete")
	j := initJenkins()
	if err := j.SetCrumb(req); err != nil {
		t.Errorf("Error: %v", err)
	}
	req.SetBasicAuth("admin", "B1OvhHMnPPxXz4kFwODYIh")
	req.Debug = true
	status, body, err := req.End()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if status.StatusCode != 200 {
		t.Errorf("Error: %v", status)
		fmt.Printf("%+v\n", body)
	}

}
