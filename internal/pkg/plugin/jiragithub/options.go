package jiragithub

// Options is the struct for configurations of the jiragithub plugin.
type Options struct {
	Owner          string `validate:"required"`
	Org            string `validate:"required"`
	Repo           string `validate:"required"`
	JiraBaseUrl    string
	JiraUserEmail  string
	JiraProjectKey string
	Branch         string `validate:"required"`
}
