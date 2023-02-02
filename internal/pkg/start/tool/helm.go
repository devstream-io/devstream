package tool

import (
	"fmt"
	"os/exec"
)

var toolHelm = tool{
	Name: "Helm",

	IfExists: func() bool {
		_, err := exec.LookPath("helm")
		return err == nil
	},

	Install: func() error {
		if !confirm("Helm") {
			return fmt.Errorf("user cancelled")
		}

		if err := execCommand([]string{"brew", "install", "helm"}); err != nil {
			return err
		}
		return nil
	},
}
