package configloader

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// render variables if var file is not empty
func renderVariables(varFileName string, configFileBytes []byte) ([]byte, error) {
	if varFileName == "" {
		return configFileBytes, nil
	}

	// load variables file
	variables, err := loadVariablesFilesIntoMap(varFileName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// handle variables format
	configFileContentString := addDotForVariablesInConfig(string(configFileBytes))

	// render config with variables
	result, err := renderConfigWithVariables(configFileContentString, variables)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return result, nil
}

func loadVariablesFilesIntoMap(varFileName string) (map[string]interface{}, error) {
	fileBytes, err := ioutil.ReadFile(varFileName)
	if err != nil {
		return nil, err
	}

	log.Debugf("Variables file: \n%s\n", string(fileBytes))

	variables := make(map[string]interface{})
	err = yaml.Unmarshal(fileBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// this is because our variables' syntax is [[ varName ]]
// while Go's template is [[ .varName ]]
func addDotForVariablesInConfig(s string) string {
	// regex := `\[\[\s*(.*)\s*\]\]`
	// r := regexp.MustCompile(regex)
	// return r.ReplaceAllString(s, "[[ .$1 ]]")
	regex := `\[\[\s*`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "[[ .")
}

func renderConfigWithVariables(fileContent string, variables map[string]interface{}) ([]byte, error) {
	tpl, err := template.New("configfile").Delims("[[", "]]").Parse(fileContent)
	if err != nil {
		return nil, err
	}

	var results bytes.Buffer
	err = tpl.Execute(&results, variables)
	if err != nil {
		return nil, err
	}

	return results.Bytes(), err
}
