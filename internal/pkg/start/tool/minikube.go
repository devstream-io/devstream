package tool

import (
	"fmt"
	"os/exec"
)

var toolMinikube = tool{
	Name: "Minikube",

	IfExists: func() bool {
		_, err := exec.LookPath("minikube")
		return err == nil
	},

	Install: func() error {
		if !confirm("Minikube") {
			return fmt.Errorf("user cancelled")
		}

		if err := execCommand([]string{"brew", "install", "minikube"}); err != nil {
			return err
		}
		return nil
	},

	IfStopped: func() bool {
		cmd := exec.Command("minikube", "status")
		return cmd.Run() != nil
	},

	Start: func() error {
		if err := execCommand([]string{"minikube", "start"}); err != nil {
			return err
		}
		return nil
	},
}
