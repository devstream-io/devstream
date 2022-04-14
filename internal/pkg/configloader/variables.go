package configloader

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const defaultVarFileName = "variables.yaml"

func renderVariables(varFileName string, configFileBytes []byte) ([]byte, error) {
	// if the var file is default (user didn't overwrite the value with --var-file option)
	// and the default var file doesn't exist, do nothing
	// it's OK to not use a var file
	if defaultVarFileName == varFileName {
		if _, err := os.Stat(defaultVarFileName); errors.Is(err, os.ErrNotExist) {
			return configFileBytes, nil
		}
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

// this is because our variables syntax is [[ varName ]]
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
