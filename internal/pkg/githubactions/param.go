package githubactions

// Workflow is the struct for a GitHub Actions workflow.
type Workflow struct {
	commitMessage    string
	workflowFileName string
	workflowContent  string
}
