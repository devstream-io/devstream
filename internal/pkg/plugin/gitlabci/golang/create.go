package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	opts, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.NewClient()
	if err != nil {
		return nil, err
	}

	ciFileContent, err := client.GetGitLabCIGolangTemplate()
	if err != nil {
		return nil, err
	}

	if err = client.CommitSingleFile(opts.PathWithNamespace, opts.Branch, commitMessage, ciFileName, ciFileContent); err != nil {
		return nil, err
	}

	return buildState(opts), nil
}
