package githubactions

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
)

var pipelines = []Pipeline{
	Pipeline{"pr builder by DevStream", "pr-builder.yaml", prBuilder},
	Pipeline{"master builder by DevStream", "master-builder.yaml", masterBuilder},
}

func Install(options *map[string]interface{}) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		log.Fatalln(err)
	}

	for _, pipeline := range pipelines {
		createFile(&Param{
			&ctx,
			getGitHubClient(&ctx),
			&opt,
			&pipeline,
		})
	}
}
