package jenkins

import (
	_ "embed"
	"fmt"

	"github.com/parnurzeal/gorequest"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// this file is used to write custom implementation of jenkins client

// GetCrumb returns the crumb for jenkins,
// jenkins uses crumb to prevent CSRF(cross-site request forgery),
// ref: https://www.jenkins.io/doc/upgrade-guide/2.176/#upgrading-to-jenkins-lts-2-176-3
// ref: https://stackoverflow.com/questions/44711696/jenkins-403-no-valid-crumb-was-included-in-the-request
func (j *Jenkins) GetCrumb() (crumbHeaderKey, crumbHeaderValue string, cookie string, err error) {
	// crumb response format:
	type CrumbResponse struct {
		Crumb             string `json:"crumb"`
		CrumbRequestField string `json:"crumbRequestField"`
	}
	var crumbResp CrumbResponse

	// get crumb
	request := gorequest.New()
	getCrumbURL := j.JenkinsURL + `/crumbIssuer/api/json`
	resp, body, errs := request.Get(getCrumbURL).
		SetBasicAuth(j.BasicAuth.Username, j.BasicAuth.Password).
		EndStruct(&crumbResp)

	// check error
	if len(errs) != 0 {
		return "", "", "", fmt.Errorf("failed to get jenkins crumb: %s", errs)
	}
	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("failed to get jenkins crumb, here is response: %s", body)
	}

	log.Debugf("crumb: %+v", crumbResp)

	return crumbResp.CrumbRequestField, crumbResp.Crumb, resp.Header.Get("set-cookie"), nil
}

// SetCrumb sets the jenkins crumb to the request
func (j *Jenkins) SetCrumb(req *gorequest.SuperAgent) error {
	crumbHeaderKey, crumbHeaderValue, cookie, err := j.GetCrumb()
	if err != nil {
		return err
	}

	// all these three should be set in the request, or it will cause 403
	req.Set(crumbHeaderKey, crumbHeaderValue)
	req.Set("Cookie", cookie)
	req.SetBasicAuth(j.BasicAuth.Username, j.BasicAuth.Password)

	return nil
}
