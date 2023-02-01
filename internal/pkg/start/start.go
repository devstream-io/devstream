package start

import (
	"fmt"

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
	// TODO(daniel-hutao)
	return false
}

func helmExists() bool {
	// TODO(daniel-hutao)
	return false
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
		return nil
	}

	// TODO(daniel-hutao): install Docker
	fmt.Println("Docker installing...")
	fmt.Println("Docker installed.")
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
		return nil
	}

	// TODO(daniel-hutao): install Minikube
	fmt.Println("Minikube installing...")
	fmt.Println("Minikube installed.")
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
		return nil
	}

	// TODO(daniel-hutao): install Helm
	fmt.Println("Helm installing...")
	fmt.Println("Helm installed.")
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
		return nil
	}

	// TODO(daniel-hutao): install Argo CD
	fmt.Println("Argo CD installing...")
	fmt.Println("Argo CD installed.")
	return nil
}
