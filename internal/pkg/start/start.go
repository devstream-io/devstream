package start

import (
	"fmt"
)

func Start() error {
	fmt.Println("Let's get started.")

	err := installToolsIfNotExist()
	if err != nil {
		return err
	}

	fmt.Println("Enjoy it!☺️")
	return nil
}

func installToolsIfNotExist() error {
	for _, t := range tools {
		if !t.Exists() {
			if err := t.Install(); err != nil {
				return err
			}
		}
	}
	return nil
}
