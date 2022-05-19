package githubactions

const (
	CommitMessage       = "GitHub Actions workflow, created by DevStream"
	PRBuilderFileName   = "pr-builder.yml"
	MainBuilderFileName = "main-builder.yml"
)

// Language is the struct containing details of a programming language specified in the GitHub Actions Workflow.
type Language struct {
	Name    string `validate:"required"`
	Version string
}
