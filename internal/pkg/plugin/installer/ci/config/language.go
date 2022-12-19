package config

import (
	"strings"

	"github.com/devstream-io/devstream/pkg/util/types"
)

type LanguageOption struct {
	Name      string `mapstructure:"name"`
	Version   string `mapstructure:"version"`
	FrameWork string `mapstructure:"frameWork"`
}

var (
	// aliasMap is used to map language/framework alias to standard name devstream used
	aliasMap = map[string]string{
		"golang":      "go",
		"spring-boot": "springboot",
	}
	languageDefaultOptionMap = map[string]*GeneralDefaultOption{
		"java": {
			Test: &TestOption{
				Command:       []string{"mvn -B test"},
				ContainerName: "maven:3.8.1-jdk-8",
				Enable:        types.Bool(true),
			},
			defaultVersion: "8",
		},
		"go": {
			Test: &TestOption{
				Enable:          types.Bool(true),
				ContainerName:   "golang:1.18",
				Command:         []string{"go test ./..."},
				CoverageCommand: "go tool cover -func=coverage.out >> coverage.cov",
				CoverageStatusCommand: `cat coverage.cov
body=$(cat coverage.cov)
body="${body//'%'/'%25'}"
body="${body//$'\n'/'%0A'}"
body="${body//$'\r'/'%0D'}"
echo ::set-output name=body::$body`,
			},
			defaultVersion: "1.17",
		},
		"python": {
			Test: &TestOption{
				Command: []string{
					"python -m pip install --upgrade pip",
					"pip install -r requirements.txt",
					"python3 -m unittest",
				},
				ContainerName: "python:3.10.9",
				Enable:        types.Bool(true),
			},
			defaultVersion: "3.9",
		},
		"nodejs": {
			Test: &TestOption{
				Command: []string{
					"npm ci",
					"npm run build --if-present",
					"npm test",
				},
				ContainerName: "nodejs:18",
				Enable:        types.Bool(true),
			},
			defaultVersion: "18",
		},
	}
	frameworkDefaultOptionMap = map[string]*GeneralDefaultOption{
		"gin":        languageDefaultOptionMap["go"],
		"springboot": languageDefaultOptionMap["java"],
		"django":     languageDefaultOptionMap["python"],
	}
	framworkLanguageMap = map[string]string{
		"gin":        "go",
		"springboot": "java",
		"django":     "python",
	}
)

// GetGeneralDefaultOpt will return language/frameWork default ci/cd options
func (l *LanguageOption) GetGeneralDefaultOpt() *GeneralDefaultOption {
	// if frameWork is configured, use frameWork first
	// else use languange name for defalt options
	var defaultOpt *GeneralDefaultOption
	if l.FrameWork != "" {
		frameWork := strings.TrimSpace(strings.ToLower(l.FrameWork))
		if frameWorkAlias, exist := aliasMap[frameWork]; exist {
			frameWork = frameWorkAlias
		}
		if l.Name == "" {
			l.Name = framworkLanguageMap[l.FrameWork]
		}
		defaultOpt = frameworkDefaultOptionMap[frameWork]
	}

	if l.Name != "" {
		lang := strings.TrimSpace(strings.ToLower(l.Name))
		if aliasLang, exist := aliasMap[lang]; exist {
			lang = aliasLang
			l.Name = aliasLang
		}
		if defaultOpt == nil {
			defaultOpt = languageDefaultOptionMap[lang]
		}
	}
	if defaultOpt != nil {
		if l.Version == "" {
			l.Version = defaultOpt.defaultVersion
		}
	}
	return defaultOpt
}

// IS
func (l LanguageOption) IsConfigured() bool {
	return l.Name != "" || l.FrameWork != ""
}
