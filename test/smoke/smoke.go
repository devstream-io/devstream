package main

import (
	"log"
	"os"

	"github.com/merico-dev/stream/internal/pkg/pluginengine"
	dsos "github.com/merico-dev/stream/internal/pkg/util/os"
)

const configFileForSmoke = "config.yaml"

func main() {
	if _, err := os.Stat("dtm"); err != nil {
		log.Fatal(err)
	}

	// TODO(daniel-hutao): How to deal with the GitHub token with githubactions?
	err := dsos.ExecInSystem(".", []string{"./dtm", "apply", "-f", configFileForSmoke}, nil, true)
	if err != nil {
		log.Fatal(err)
	}

	healthy := pluginengine.CheckHealthy(configFileForSmoke)

	if healthy {
		log.Println("all tools are healthy")
	}

	log.Fatalf("some tools are not healthy")
}
