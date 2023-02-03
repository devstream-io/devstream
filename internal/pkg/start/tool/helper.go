package tool

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

func confirm(name string) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Install %s now", name),
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Printf("%s must be installed. Quit now.", name)
		return false
	}
	return true
}

func execCommand(cmdStr []string) error {
	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
