package jenkins

import (
	_ "embed"

	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

//go:embed tpl/casc.tpl.groovy
var cascGroovyScript string

func (jenkins *jenkins) ExecuteScript(script string) (string, error) {
	now := time.Now().Unix()
	verifier := fmt.Sprintf("verifier-%d", now)
	output := ""
	fullScript := fmt.Sprintf("%s\nprint println('%s')", script, verifier)

	data := url.Values{}
	data.Set("script", fullScript)

	ar := gojenkins.NewAPIRequest("POST", "/scriptText", bytes.NewBufferString(data.Encode()))
	if err := jenkins.Requester.SetCrumb(jenkins.ctx, ar); err != nil {
		return output, err
	}
	ar.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	ar.Suffix = ""

	r, err := jenkins.Requester.Do(jenkins.ctx, ar, &output, nil)
	if err != nil {
		return "", fmt.Errorf("couldn't execute groovy script, logs '%s'", output)
	}
	log.Debugf("------> %s\n%s", output, script)
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return output, errors.Errorf("invalid status code '%d', logs '%s'", r.StatusCode, output)
	}

	if !strings.Contains(output, verifier) {
		return output, fmt.Errorf("jenkins run script return error: %s", output)
	}

	return output, nil
}

func (jenkins *jenkins) ConfigCasc(cascConfig string) error {
	groovyCascScript, err := template.Render(
		"jenkins casc", cascGroovyScript, map[string]string{
			"CascConfig": cascConfig,
		},
	)
	if err != nil {
		log.Debugf("jenkins render casc failed: %s", err)
		return err
	}
	_, err = jenkins.ExecuteScript(groovyCascScript)
	return err
}
