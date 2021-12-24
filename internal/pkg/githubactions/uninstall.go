package githubactions

// Uninstall remove GitHub Actions workflows.
func Uninstall(options *map[string]interface{}) (bool, error) {
	githubActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	for _, pipeline := range workflows {
		err := githubActions.DeleteWorkflow(pipeline)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
