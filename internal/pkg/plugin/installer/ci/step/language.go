package step

import (
	"strings"

	"github.com/devstream-io/devstream/pkg/util/types"
)

type Language struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type generalDefaultOption struct {
	testOption *test
}

var languageDefaultOptionMap = map[string]*generalDefaultOption{
	"java": {
		testOption: &test{
			Command:       "mvn -B test",
			ContainerName: "maven:3.8.1-jdk-8",
			Enable:        types.Bool(true),
		},
	},
	"golang": {
		testOption: &test{
			Enable:          types.Bool(true),
			Command:         "go test ./...",
			CoverageCommand: "go tool cover -func=coverage.out >> coverage.cov",
			CoverageStatusCommand: `cat coverage.cov
body=$(cat coverage.cov)
body="${body//'%'/'%25'}"
body="${body//$'\n'/'%0A'}"
body="${body//$'\r'/'%0D'}"
echo ::set-output name=body::$body`,
		},
	},
}

func (l *Language) getGeneralDefaultOption() *generalDefaultOption {
	lang := strings.TrimSpace(strings.ToLower(l.Name))
	return languageDefaultOptionMap[lang]
}
