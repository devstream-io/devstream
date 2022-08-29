package jenkins

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/bndr/gojenkins"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type GitlabCredentials struct {
	XMLName     xml.Name `xml:"com.dabsquared.gitlabjenkins.connection.GitLabApiTokenImpl"`
	ID          string   `xml:"id"`
	Scope       string   `xml:"scope"`
	Description string   `xml:"description"`
	APIToken    string   `xml:"apiToken"`
}

func (j *jenkins) CreateGiltabCredential(id, gitlabToken string) error {
	cred := GitlabCredentials{
		ID:          id,
		Scope:       credentialScope,
		APIToken:    gitlabToken,
		Description: "this credential is created by devstream",
	}

	cm := &gojenkins.CredentialsManager{
		J: &j.Jenkins,
	}
	err := cm.Add(j.ctx, domain, cred)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Debugf("jenkins credential %s has been created", id)
			return nil
		}
		return fmt.Errorf("could not create credential: %v", err)
	}

	// get credential to validate creation
	getCred := GitlabCredentials{}
	if err = cm.GetSingle(j.ctx, domain, cred.ID, &getCred); err != nil {
		return fmt.Errorf("could not get credential: %v", err)
	}
	return nil
}
