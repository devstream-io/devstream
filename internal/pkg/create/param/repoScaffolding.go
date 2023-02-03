package param

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	arrowUp                 = "↑"
	fingerUp                = "👆"
	backToLanguageSelection = arrowUp + " Back to language selection"
)

func selectRepoScaffolding() (language, framework, url string, err error) {
	languagesRepoMap := make(map[string][]RepoScaffolding)
	for _, repo := range ListRepoScaffolding() {
		languagesRepoMap[repo.Language] = append(languagesRepoMap[repo.Language], repo)
	}

	// use a loop to allow user to go back to language selection
	for {
		language, err = selectLanguage(languagesRepoMap)
		if err != nil {
			return
		}

		fmt.Println("\nPlease choose a framework next.")

		framework, url, err = selectFrameworks(languagesRepoMap[language])
		if err != nil {
			return
		}
		if framework != backToLanguageSelection {
			break
		}
	}
	language = strings.ToLower(language)
	framework = strings.ToLower(framework)

	return language, framework, url, nil

}

func selectLanguage(languagesRepoMap map[string][]RepoScaffolding) (language string, err error) {
	var langFrameworkList []struct {
		Language   string
		Frameworks []string
	}

	for lang, repos := range languagesRepoMap {
		var frameworks []string
		for _, repo := range repos {
			frameworks = append(frameworks, repo.Framework)
		}
		langFrameworkList = append(langFrameworkList, struct {
			Language   string
			Frameworks []string
		}{lang, frameworks})
	}

	combinedFuncMap := promptui.FuncMap
	combinedFuncMap["join"] = func(list []string, sep string) string {
		return strings.Join(list, sep)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "🐰 {{ .Language | blue }}",
		Inactive: "  {{ .Language | cyan }}",
		Selected: "🐰 Language: {{ .Language | blue | cyan }}",
		Details: `
--------- Supported Frameworks For {{ .Language}} ----------
{{ join .Frameworks ", " }}`,
		FuncMap: combinedFuncMap,
	}

	searcher := func(input string, index int) bool {
		lang := langFrameworkList[index]
		name := strings.Replace(strings.ToLower(lang.Language), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose your language",
		Items:     langFrameworkList,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	return langFrameworkList[i].Language, nil
}

func selectFrameworks(repos []RepoScaffolding) (framework, url string, err error) {
	repos = append(repos, RepoScaffolding{
		Name:      backToLanguageSelection,
		Framework: backToLanguageSelection,
		URL:       fingerUp, Description: fingerUp, Language: fingerUp})
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "🐰 {{ .Framework | blue }}",
		Inactive: "  {{ .Framework | cyan }}",
		Selected: "🐰 Framework: {{ .Framework | blue | cyan }}",
		Details: `
--------- Language/Framework Details ----------
{{ "Repo Template Name:" | faint }}	{{ .Name }}
{{ "Language:" | faint }}	{{ .Language }}
{{ "Framework:" | faint }}	{{ .Framework }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		repo := repos[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Choose your framework",
		Items:     repos,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	return repos[i].Framework, repos[i].URL, nil
}
