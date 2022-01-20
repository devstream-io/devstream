package githubactions

import "github.com/merico-dev/stream/internal/pkg/util/log"

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := githubActions.GetLanguage()
	log.Infof("Language is: %s.", language.String())
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())

	for _, w := range ws {
		if err := githubActions.renderTemplate(w); err != nil {
			return false, err
		}
		if err := githubActions.AddWorkflow(w); err != nil {
			return false, err
		}
	}

	return true, nil
}
