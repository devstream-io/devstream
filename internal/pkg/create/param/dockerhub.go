package param

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func getDockerHubUsername() (string, error) {
	prompt := promptui.Prompt{
		Label:    "What is your DockerHub username",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Failed to get DockerHub username %v\n", err)
		return "", err
	}

	return result, nil
}

func getDockerHubToken() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Please input your DockerHub Personal Access Token",
		Mask:     '*',
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Failed to get DockerHub token %v\n", err)
		return "", err
	}

	return result, nil
}
