package pluginengine

import (
	"fmt"
	"regexp"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func HandleOutputsReferences(smgr statemanager.Manager, options map[string]interface{}) []error {
	errorsList := make([]error, 0)

	for optionKey, optionValue := range options {
		// only process string values in the options
		// since all outputs references are strings, not ints, not booleans, not maps
		if optionValueStr, ok := optionValue.(string); ok {
			match, toolName, instanceID, outputReferenceKey := getToolNamePluginOutputKey(optionValueStr)
			// do nothing, if the value string isn't in the format of a valid output reference
			if !match {
				continue
			}

			outputs, err := smgr.GetOutputs(statemanager.GenerateStateKeyByToolNameAndPluginKind(toolName, instanceID))
			if err != nil {
				errorsList = append(errorsList, err)
				continue
			}

			if val, ok := outputs.(map[string]interface{})[outputReferenceKey]; ok {
				options[optionKey] = replaceOutputKeyWithValue(optionValueStr, val.(string))
			} else {
				errorsList = append(errorsList, fmt.Errorf("can't find Output reference key %s", outputReferenceKey))
			}
		}

		// recursive if the value is a map (which means Tool.Option is a nested map)
		optionValueMap, ok := optionValue.(map[string]interface{})
		if ok {
			errorsList = append(errorsList, HandleOutputsReferences(smgr, optionValueMap)...)
		}
	}
	return errorsList
}

// getToolNamePluginKindAndOutputReferenceKey returns (false, "", "", "") if regex doesn't match
// if match, returns (true, name, instanceID, key)
func getToolNamePluginOutputKey(s string) (bool, string, string, string) {
	regex := `.*\${{\s*([^.]*)\.([^.]*)\.outputs\.([^.\s]*)\s*}}.*`
	r := regexp.MustCompile(regex)
	if !r.MatchString(s) {
		return false, "", "", ""
	}
	results := r.FindStringSubmatch(s)
	return true, results[1], results[2], results[3]
}

func replaceOutputKeyWithValue(s, val string) string {
	regex := `\${{\s*[^.]*\.[^.]*\.outputs\.[^.]*\s*}}`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, val)
}
