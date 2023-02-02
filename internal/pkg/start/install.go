package start

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/manifoldco/promptui"
)

func installDocker() error {
	if !confirm("Docker") {
		return fmt.Errorf("user cancelled")
	}

	if err := execCommand([]string{"brew", "install", "docker", "--cask"}); err != nil {
		return err
	}
	if err := execCommand([]string{"open", "-a", "Docker"}); err != nil {
		return err
	}

	return waitForDockerRun()
}

func installMinikube() error {
	if !confirm("Minikube") {
		return fmt.Errorf("user cancelled")
	}

	if err := execCommand([]string{"brew", "install", "minikube"}); err != nil {
		return err
	}
	return nil
}

func installHelm() error {
	if !confirm("Helm") {
		return fmt.Errorf("user cancelled")
	}

	if err := execCommand([]string{"brew", "install", "helm"}); err != nil {
		return err
	}
	return nil
}

func installArgocd() error {
	if !confirm("Argo CD") {
		return fmt.Errorf("user cancelled")
	}

	if err := execCommand([]string{"helm", "repo", "add", "argo", "https://argoproj.github.io/argo-helm"}); err != nil {
		return err
	}
	if err := execCommand([]string{"helm", "install", "argo/argo-cd", "-n", "argocd", "--create-namespace"}); err != nil {
		return err
	}
	return nil
}

func waitForDockerRun() error {
	fmt.Println("\nPlease make sure your docker has been started. The OS may ask you to authorize it manually.")
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

func confirm(name string) bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Install %s now.", name),
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
