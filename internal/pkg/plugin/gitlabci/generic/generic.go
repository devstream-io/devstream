package generic

import (
	"io"
	"net/http"
)

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

func buildState(opts *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opts.PathWithNamespace,
		"branch":            opts.Branch,
		"templateURL":       opts.TemplateURL,
		"templateVariables": opts.TemplateVariables,
	}
}

func download(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
}
