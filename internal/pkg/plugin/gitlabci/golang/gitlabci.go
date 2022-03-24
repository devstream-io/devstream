package golang

const (
	ciFileName    string = ".gitlab-ci.yml"
	commitMessage string = "managed by DevStream"
)

type Options struct {
	PathWithNamespace string
	Branch            string
}

func buildState(opts *Options) map[string]interface{} {
	return map[string]interface{}{
		"pathWithNamespace": opts.PathWithNamespace,
		"branch":            opts.Branch,
	}
}
