package create

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type language struct {
	Name string
}

func getLanguage() (string, error) {
	languages := []language{
		{Name: "Python"},
		{Name: "Java"},
		{Name: "Go"},
		{Name: "Node.js"},
		{Name: "C++"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Language ----------
{{ "Name:" | faint }}	{{ .Name }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := languages[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose your language",
		Items:     languages,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return "", fmt.Errorf("got unsupported language")
	}

	retLang := languages[i].Name
	fmt.Printf("You choosed: %s.\n", retLang)
	return retLang, nil
}
