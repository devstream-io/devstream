package jiragithub

// Options is the struct for configurations of the jiragithub plugin.
type Options struct {
	Owner          string
	Org            string
	Repo           string
	JiraBaseUrl    string
	JiraUserEmail  string
	JiraProjectKey string
	Branch         string
}
