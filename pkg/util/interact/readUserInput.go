package interact

import (
	"fmt"
	"os"

	"github.com/tcnksm/go-input"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// AskUserIfContinue asks the user if he wants to continue
// default is false
func AskUserIfContinue(query string) (continued bool) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

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
	return userInput == "y"
}
