package pluginengine

import (
	"fmt"
	"regexp"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// HandleOutputsReferences renders outputs references in config file recursively.
// The parameter options will be changed.
func HandleOutputsReferences(smgr statemanager.Manager, options configmanager.RawOptions) []error {
	errorsList := make([]error, 0)

	for optionKey, optionInterface := range options {
		switch optionValue := optionInterface.(type) {
		// only process string values in the options
		// since all outputs references are strings, not ints, not booleans, not maps
		case string:
			log.Debugf("Before: %s: %s", optionKey, optionValue)
			match, toolName, instanceID, outputReferenceKey := getToolNamePluginOutputKey(optionValue)
			// do nothing, if the value string isn't in the format of a valid output reference
			if !match {
				continue
			}
			outputs, err := smgr.GetOutputs(statemanager.GenerateStateKeyByToolNameAndInstanceID(toolName, instanceID))
			if err != nil {
				errorsList = append(errorsList, err)
				continue
			}
			if val, ok := outputs[outputReferenceKey]; ok {
				options[optionKey] = replaceOutputKeyWithValue(optionValue, val.(string))
				log.Debugf("After: %s: %s", optionKey, options[optionKey])
			} else {
				errorsList = append(errorsList, fmt.Errorf("can't find Output reference key %s", outputReferenceKey))
			}
		case configmanager.RawOptions:
			// recursive if the value is a map (which means Tool.Option is a nested map)
			log.Debugf("Got nested map: %v", optionValue)
			errorsList = append(errorsList, HandleOutputsReferences(smgr, optionValue)...)
		}
	}

	log.Debugf("Final options: %v.", options)

	return errorsList
}

// getToolNamePluginKindAndOutputReferenceKey returns (false, "", "", "") if regex doesn't match
// if matched, returns (true, name, instanceID, key)
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
