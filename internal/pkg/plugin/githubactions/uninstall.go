package githubactions

import "github.com/merico-dev/stream/internal/pkg/util/log"

// Uninstall remove GitHub Actions workflows.
func Uninstall(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := githubActions.GetLanguage()
	log.Infof("language is %s", language.String())
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())

	for _, pipeline := range ws {
		err := githubActions.DeleteWorkflow(pipeline)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
