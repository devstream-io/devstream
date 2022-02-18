package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Update(options *map[string]interface{}) (map[string]interface{}, error) {
	opt, err := parseAndValidateOptions(options)
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
	if err = client.UpdateSingleFile(opt.PathWithNamespace, opt.Branch, commitMessage, ciFileName, ciFileContent); err != nil {
		return nil, err
	}

	return buildState(opt), nil
}
