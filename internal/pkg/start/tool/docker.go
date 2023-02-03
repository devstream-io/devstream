package tool

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/manifoldco/promptui"
)

var toolDocker = tool{
	Name: "Docker",

	IfExists: func() bool {
		_, err := exec.LookPath("docker")
		return err == nil
	},

	Install: func() error {
		if !confirm("Docker") {
			return fmt.Errorf("user cancelled")
		}

		if err := execCommand([]string{"brew", "install", "docker", "--cask"}); err != nil {
			return err
		}
		return nil
	},

	IfStopped: func() bool {
		cmd := exec.Command("docker", "version")
		return cmd.Run() != nil
	},

	Start: func() error {
		if err := execCommand([]string{"open", "-a", "Docker"}); err != nil {
			return err
		}

		return waitForDockerRun()
	},
}

func waitForDockerRun() error {
	fmt.Println("\nI've tried to start Docker for you.")
	time.Sleep(time.Second)
	fmt.Println("But the OS may ask you to authorize it manually.")
	time.Sleep(time.Second)
	fmt.Println("Please make sure your docker has been started.")
	fmt.Println()
	time.Sleep(time.Second)
	prompt := promptui.Prompt{
		Label:     "I've verified that Docker is running properly by using the `docker version` command.",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Please make sure docker starts properly first.")
		return err
	}

	return nil
}
