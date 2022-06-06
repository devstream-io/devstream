package configloader

import (
	"bytes"
	"html/template"
	"regexp"
)

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
