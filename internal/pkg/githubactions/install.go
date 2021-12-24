package githubactions

var workflows = []Workflow{
	{"pr builder by DevStream", "pr-builder.yml", prBuilder},
	{"master builder by DevStream", "master-builder.yml", masterBuilder},
}

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	for _, pipeline := range workflows {
		err := githubActions.AddWorkflow(pipeline)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
