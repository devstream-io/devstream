package main

import (
	"fmt"
	"log"
	"os"

	"github.com/merico-dev/stream/test/smoke/argocd"
)

func main() {
	if _, err := os.Stat("dtm"); err != nil {
		log.Fatal(err)
	}

	// TODO(daniel-hutao): How to deal with the GitHub token with githubactions?
	err := ExecInSystem(".", []string{"./dtm", "apply", "-f", "config.yaml"}, nil, true)
	if err != nil {
		log.Fatal(err)
	}

	if err := checkPlugins(); err != nil {
		log.Fatal(err)
	}
}

func checkPlugins() error {
	plugins := []Plugin{argocd.NewArgocd()}

	for _, p := range plugins {
		health, err := p.Health()
		if err != nil {
			return err
		}
		if !health {
			return fmt.Errorf("plugin %s is not healthy", p.Name())
		}
	}
	return nil
}
