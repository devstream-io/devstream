package githubactions

import "fmt"

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := githubActions.GetLanguage()
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())

	fmt.Println("lang is ", language.Name, language.Version, language.String())
	fmt.Println("ws is ", ws)

	for _, pipeline := range ws {
		err := githubActions.AddWorkflow(pipeline)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
