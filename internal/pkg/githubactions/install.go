package githubactions

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
)

var workflows = []Workflow{
	{"pr builder by DevStream", "pr-builder.yaml", prBuilder},
	{"master builder by DevStream", "master-builder.yaml", masterBuilder},
}

// Install sets up GitHub Actions workflows.
func Install(options *map[string]interface{}) {
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
}
