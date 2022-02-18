package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Read(options *map[string]interface{}) (map[string]interface{}, error) {
	opt, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.NewClient()
	if err != nil {
		return nil, err
	}

	if err = client.FileExists(opt.PathWithNamespace, opt.Branch, ciFileName); err != nil {
		return nil, err
	}

	return buildState(opt), nil
}
