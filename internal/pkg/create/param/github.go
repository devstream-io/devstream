package param

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func getGitHubUsername() (string, error) {
	prompt := promptui.Prompt{
		Label:    "What is your GitHub username",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Failed to get GitHub username %v\n", err)
		return "", err
	}

	return result, nil
}

func getGitHubRepo() (string, error) {
	prompt := promptui.Prompt{
		Label:    "What GitHub Repo You Want to Create",
		Validate: validate,
		Default:  "firstapp",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Failed to get GitHub repo %v\n", err)
		return "", err
	}

	return result, nil
}

func getGitHubToken() (string, error) {
	prompt := promptui.Prompt{
		Label:    "Please input your GitHub Personal Access Token",
		Mask:     '*',
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Failed to get GitHub token %v\n", err)
		return "", err
	}

	return result, nil
}
