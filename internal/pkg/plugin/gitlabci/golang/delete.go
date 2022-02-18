package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Delete(options *map[string]interface{}) (bool, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return false, err
	}

	client, err := gitlab.NewClient()
	if err != nil {
		return false, err
	}

	if err = client.DeleteSingleFile(opt.PathWithNamespace, opt.Branch, commitMessage, ciFileName); err != nil {
		return false, err
	}

	return true, nil
}
