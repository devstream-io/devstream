package step

import (
	"strings"

	"github.com/devstream-io/devstream/pkg/util/types"
)

type language struct {
	Name      string `mapstructure:"name"`
	Version   string `mapstructure:"version"`
	FrameWork string `mapstructure:"frameWork"`
}

type generalDefaultOption struct {
	testOption *test
}

var languageDefaultOptionMap = map[string]*generalDefaultOption{
	"java": {
		testOption: &test{
			Command:       []string{"mvn -B test"},
			ContainerName: "maven:3.8.1-jdk-8",
			Enable:        types.Bool(true),
		},
	},
	"go": {
		testOption: &test{
			Enable:          types.Bool(true),
			Command:         []string{"go test ./..."},
			CoverageCommand: "go tool cover -func=coverage.out >> coverage.cov",
			CoverageStatusCommand: `cat coverage.cov
body=$(cat coverage.cov)
body="${body//'%'/'%25'}"
body="${body//$'\n'/'%0A'}"
body="${body//$'\r'/'%0D'}"
echo ::set-output name=body::$body`,
		},
	},
	"python": {
		testOption: &test{
			Command: []string{
				"python -m pip install --upgrade pip",
				"pip install -r requirements.txt",
				"python3 -m unittest",
			},
			Enable: types.Bool(true),
		},
	},
	"nodejs": {
		testOption: &test{
			Command: []string{
				"npm ci",
				"npm run build --if-present",
				"npm test",
			},
			Enable: types.Bool(true),
		},
	},
}

func (l *language) getGeneralDefaultOption() *generalDefaultOption {
	lang := strings.TrimSpace(strings.ToLower(l.Name))
	return languageDefaultOptionMap[lang]
}
