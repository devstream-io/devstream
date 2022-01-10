package githubactions

import "log"

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := githubActions.GetLanguage()
	log.Printf("Language is: %s.", language.String())
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())

	for _, pipeline := range ws {
		err := githubActions.AddWorkflow(pipeline)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
