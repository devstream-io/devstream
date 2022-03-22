package golang

import "github.com/merico-dev/stream/pkg/util/gitlab"

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	opts, err := parseAndValidateOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := gitlab.NewClient()
	if err != nil {
		return nil, err
	}

	exists, err := client.FileExists(opts.PathWithNamespace, opts.Branch, ciFileName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return buildState(opts), nil
}
