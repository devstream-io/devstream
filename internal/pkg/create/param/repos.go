package param

func ListRepoScaffolding() []RepoScaffolding {
	return []RepoScaffolding{
		{Name: "dtm-repo-scaffolding-golang-gin",
			URL:         "https://github.com/devstream-io/dtm-repo-scaffolding-golang-gin",
			Language:    "Golang",
			Framework:   "Gin",
			Description: "This is a scaffolding for Golang web app based on Gin framework",
		},
		{Name: "dtm-repo-scaffolding-python-flask",
			URL:         "https://github.com/devstream-io/dtm-repo-scaffolding-python-flask",
			Language:    "Python",
			Framework:   "Flask",
			Description: "This is a scaffolding for Python web app based on Flask framework",
		},
		{Name: "dtm-repo-scaffolding-java-springboot",
			URL:         "https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot",
			Language:    "Java",
			Framework:   "SpringBoot",
			Description: "This is a scaffolding for Java web app based on SpringBoot framework",
		},
	}
}

type RepoScaffolding struct {
	Name        string
	URL         string
	Language    string
	Framework   string
	Description string
}
