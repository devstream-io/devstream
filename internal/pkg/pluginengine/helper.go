package pluginengine

import (
	"fmt"
	"os"

	"github.com/tcnksm/go-input"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func readUserInput() string {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Continue? [y/n]"
	userInput, err := ui.Ask(query, &input.Options{
		Required: true,
		Default:  "n",
		Loop:     true,
		ValidateFunc: func(s string) error {
			if s != "y" && s != "n" {
				return fmt.Errorf("input must be y or n")
			}
			return nil
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return userInput
}
