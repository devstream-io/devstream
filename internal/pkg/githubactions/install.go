package githubactions

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
)

var workflows = []Workflow{
	{"pr builder by DevStream", "pr-builder.yml", prBuilder},
	{"master builder by DevStream", "master-builder.yml", masterBuilder},
}

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) (bool, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		log.Fatalln(err)
	}

	for _, pipeline := range workflows {
		createFile(&Param{
			&ctx,
			getGitHubClient(&ctx),
			&opt,
			&pipeline,
		})
	}

	return true, nil
}
