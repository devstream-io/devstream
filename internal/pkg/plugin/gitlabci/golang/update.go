package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Update(options map[string]interface{}) (map[string]interface{}, error) {
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

	// the only difference between "Create" and "Update"
	if err = client.UpdateSingleFile(opts.PathWithNamespace, opts.Branch, commitMessage, ciFileName, ciFileContent); err != nil {
		return nil, err
	}

	return buildState(opts), nil
}
