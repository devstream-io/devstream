package pluginengine

import (
	"fmt"
	"regexp"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

const (
	// a legal output reference should be in the format of ${{ abc }}
	OUTPUT_REFERENCE_PREFIX = "${{"
	OUTPUT_REFERENCE_SUFFIX = "}}"

	// e.g., ${{ TOOL_NAME.PLUGIN.outputs.some_key }}

	// it has 4 sections, separated by "."
	OUTPUT_REFERENCE_TOTAL_SECTIONS = 4
	SECTION_SEPARATOR               = "."

	// the first section is the tool's name
	TOOL_NAME = 0

	// the second section is the plugin's kind
	PLUGIN = 1

	// the third section is a constant string "outputs"
	// and the last section is the key to refer to
	OUTPUT_REFERENCE_KEY = 3
)

func HandleOutputsReferences(smgr statemanager.Manager, options map[string]interface{}) []error {
	errorsList := make([]error, 0)

	for optionKey, optionValue := range options {
		// only process string values in the options
		// since all outputs references are strings, not ints, not booleans, not maps
		if optionValueStr, ok := optionValue.(string); ok {
			match, toolName, pluginKind, outputReferenceKey := getToolNamePluginOutputKey(optionValueStr)
			// do nothing, if the value string isn't in the format of a valid output reference
			if !match {
				continue
			}

			outputs, err := smgr.GetOutputs(statemanager.GenerateStateKeyByToolNameAndPluginKind(toolName, pluginKind))
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
// if match, returns (true, name, kind, key)
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
