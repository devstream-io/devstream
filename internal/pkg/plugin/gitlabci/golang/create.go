package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Create(options map[string]interface{}) (map[string]interface{}, error) {
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

	if err = client.CommitSingleFile(opt.PathWithNamespace, opt.Branch, commitMessage, ciFileName, ciFileContent); err != nil {
		return nil, err
	}

	return buildState(opt), nil
}
