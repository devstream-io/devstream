package jiragithub

// Options is the struct for configurations of the jiragithub plugin.
type Options struct {
	Owner          string `validate:"required_without=Org"`
	Org            string `validate:"required_without=Owner"`
	Repo           string `validate:"required"`
	JiraBaseUrl    string `validate:"required"`
	JiraUserEmail  string `validate:"required"`
	JiraProjectKey string `validate:"required"`
	Branch         string `validate:"required"`
}
