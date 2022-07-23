package jenkins

import (
	"context"
	"fmt"

	"github.com/bndr/gojenkins"
)

const (
	domain          = "_"
	credentialScope = "GLOBAL"
)

type (
	// Jenkins is a Jenkins client,
	// it is a wrapper around gojenkins.Jenkins and implements additional methods
	Jenkins struct {
		*gojenkins.Jenkins

		*BasicAuth
		JenkinsURL string
	}

	BasicAuth struct {
		Username string
		Password string
	}
)

// NewJenkins creates a new Jenkins client and validates the connection
func NewJenkins(jenkinsURL, username, password string) (*Jenkins, error) {
	// validate
	if jenkinsURL == "" {
		return nil, fmt.Errorf("jenkinsURL is required")
	}
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// create gojenkins client
	gojenkinsClient := gojenkins.CreateJenkins(nil, jenkinsURL, username, password)
	ctx := context.Background()
	// init and validate client
	if _, err := gojenkinsClient.Init(ctx); err != nil {
		return nil, err
	}

	return &Jenkins{
		Jenkins: gojenkinsClient,
		BasicAuth: &BasicAuth{
			Username: username,
			Password: password,
		},
		JenkinsURL: jenkinsURL,
	}, nil
}

func (j *Jenkins) GetCredentialManager() *gojenkins.CredentialsManager {
	return &gojenkins.CredentialsManager{
		J: j.Jenkins,
	}
}

func (j *Jenkins) CreateCredentialsUsername(username, password, id, description string) error {
	cred := gojenkins.UsernameCredentials{
		ID:          id,
		Scope:       credentialScope,
		Username:    username,
		Password:    password,
		Description: description,
	}

	// create credential
	ctx := context.Background()
	cm := j.GetCredentialManager()
	err := cm.Add(ctx, domain, cred)
	if err != nil {
		return fmt.Errorf("could not create credential: %v", err)
	}

	// get credential to validate creation
	getCred := gojenkins.UsernameCredentials{}
	if err = cm.GetSingle(ctx, domain, cred.ID, &getCred); err != nil {
		return fmt.Errorf("could not get credential: %v", err)
	}

	return nil
}

// GetCredentialsUsername returns the credentials of type username-password with the given id,
// it returns an error if the credential does not exist
func (j *Jenkins) GetCredentialsUsername(id string) (*gojenkins.UsernameCredentials, error) {
	getCred := gojenkins.UsernameCredentials{}
	ctx := context.Background()
	cm := j.GetCredentialManager()
	err := cm.GetSingle(ctx, domain, id, &getCred)
	if err != nil {
		return nil, fmt.Errorf("could not get credential: %v", err)
	}
	return &getCred, nil
}

func (j *Jenkins) DeleteCredentialsUsername(id string) error {
	return j.GetCredentialManager().Delete(context.Background(), domain, id)
}
