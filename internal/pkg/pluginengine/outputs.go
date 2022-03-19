package pluginengine

import (
	"fmt"
	"strings"

	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

const (
	// a legal output reference should be in the format of ${{ abc }}
	OUTPUT_REFERENCE_PREFIX = "${{"
	OUTPUT_REFERENCE_SUFFIX = "}}"

	// e.g., ${{ TOOL_NAME.PLUGIN_KIND.outputs.some_key }}

	// it has 4 sections, separated by "."
	OUTPUT_REFERENCE_TOTAL_SECTIONS = 4
	SECTION_SEPARATOR               = "."

	// the first section is the tool's name
	TOOL_NAME = 0

	// the second section is the plugin's kind
	PLUGIN_KIND = 1

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
			// do nothing, if the value string isn't in the format of a valid output reference
			if !isValidOutputsReferenceFormat(optionValueStr) {
				continue
			}

			toolName, pluginKind, outputReferenceKey := getToolNamePluginKindAndOutputReferenceKey(optionValueStr)

			outputs, err := smgr.GetOutputs(statemanager.GenerateStateKeyByToolNameAndPluginKind(toolName, pluginKind))
			if err != nil {
				errorsList = append(errorsList, err)
				continue
			}

			if val, ok := outputs.(map[string]interface{})[outputReferenceKey]; ok {
				options[optionKey] = val
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

// isValidOutputsReference returns true if:
// - the str is in the format of ${{ xyz }}
// - xyz has at least OUTPUT_REFERENCE_TOTAL_SECTIONS number of sections if we split it by "period"
func isValidOutputsReferenceFormat(rawOutputReference string) bool {
	if !strings.HasPrefix(rawOutputReference, OUTPUT_REFERENCE_PREFIX) || !strings.HasSuffix(rawOutputReference, OUTPUT_REFERENCE_SUFFIX) {
		return false
	}

	outputReferenceStr := stripOutputReferencePrefixAndSuffix(rawOutputReference)

	sections := strings.Split(outputReferenceStr, SECTION_SEPARATOR)
	if len(sections) < OUTPUT_REFERENCE_TOTAL_SECTIONS {
		return false
	}

	return true
}

// stripOutputReferencePrefixSuffix returns "abc" given an input "${{ abc }}"
func stripOutputReferencePrefixAndSuffix(s string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, OUTPUT_REFERENCE_PREFIX), OUTPUT_REFERENCE_SUFFIX))
}

// getToolNamePluginKindAndOutputReferenceKey returns "name, kind, key" given input "name.kind.outputs.key"
func getToolNamePluginKindAndOutputReferenceKey(s string) (string, string, string) {
	outputReferenceStr := stripOutputReferencePrefixAndSuffix(s)
	sections := strings.Split(outputReferenceStr, SECTION_SEPARATOR)
	return sections[0], sections[1], sections[3]
}
