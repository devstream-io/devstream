package start

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/manifoldco/promptui"
)

func Start() error {
	fmt.Println("Let's get started.")

	err := installToolsIfNotExist()
	if err != nil {
		return err
	}

	fmt.Println("Enjoy!")
	return nil
}

func installToolsIfNotExist() error {
	if !dockerExists() {
		if err := installDocker(); err != nil {
			return err
		}
	}

	if !minikubeExists() {
		if err := installMinikube(); err != nil {
			return err
		}
	}

	if !helmExists() {
		if err := installHelm(); err != nil {
			return err
		}
	}

	if !argocdExists() {
		if err := installArgocd(); err != nil {
			return err
		}
	}

	return nil
}

func dockerExists() bool {
	// TODO(daniel-hutao)
	return false
}

func minikubeExists() bool {
	_, err := exec.LookPath("minikube")
	return err == nil
}

func helmExists() bool {
	_, err := exec.LookPath("helm")
	return err == nil
}

func argocdExists() bool {
	// TODO(daniel-hutao)
	return false
}

func installDocker() error {
	prompt := promptui.Prompt{
		Label:     "Install Docker now",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Docker must be installed. Quit now.")
		return err
	}

	cmd := exec.Command("brew", "install", "docker", "--cask")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to install Docker. Error: %s", err)
	}

	cmd = exec.Command("open", "-a", "Docker")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to start up Docker. Error: %s", err)
		return err
	}

	return waitForDockerRun()
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

func installMinikube() error {
	prompt := promptui.Prompt{
		Label:     "Install Minikube now",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Minikube must be installed. Quit now.")
		return err
	}

	cmd := exec.Command("brew", "install", "minikube")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to install Minikube. Error: %s", err)
		return err
	}
	return nil
}

func installHelm() error {
	prompt := promptui.Prompt{
		Label:     "Install Helm now",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Helm must be installed. Quit now.")
		return err
	}

	cmd := exec.Command("brew", "install", "helm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to install Helm. Error: %s", err)
		return err
	}
	return nil
}

func installArgocd() error {
	prompt := promptui.Prompt{
		Label:     "Install Argo CD now",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Argo CD must be installed. Quit now.")
		return err
	}

	cmd := exec.Command("helm", "repo", "add", "argo", "https://argoproj.github.io/argo-helm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to add Helm repo. Error: %s", err)
		return err
	}

	cmd = exec.Command("helm", "install", "argo/argo-cd", "-n", "argocd", "--create-namespace")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Failed to install Argo CD. Error: %s", err)
		return err
	}
	return nil
}
