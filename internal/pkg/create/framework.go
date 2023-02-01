package create

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type framework struct {
	Name string
}

//TODO(daniel-hutao): support Python first
func getFramework() (string, error) {
	frameworks := []framework{
		{Name: "flask"},
		{Name: "django"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Framework ----------
{{ "Name:" | faint }}	{{ .Name }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := frameworks[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose your framework",
		Items:     frameworks,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return "", fmt.Errorf("got unsupported framework")
	}

	retFram := frameworks[i].Name
	fmt.Printf("You choosed: %s.\n", retFram)
	return retFram, nil
}
