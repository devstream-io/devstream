package configmanager

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/template"
)

func parseNestedVars(origin map[string]any) (map[string]any, error) {
	unparsed := make(map[string]any, len(origin))
	for k, v := range origin {
		unparsed[k] = v
	}

	parsed := make(map[string]any, len(origin))
	updated := true // if any vars were updated in one loop

	// loop until:
	// 1. all vars have been parsed
	// 2. or can not render more vars by existing parsed vars
	for len(unparsed) > 0 && updated {
		updated = false
		// use "parsed map" to parse each <key, value> in "unparsed map".
		// once one value doesn't contain other keys, put it into "parsed map".

		// a util func which implements the "put" steps
		putParsedValue := func(k string, v any) {
			parsed[k] = v
			delete(unparsed, k)
			updated = true
		}

		for k, v := range unparsed {
			// if value is not string or doesn't contain var, just put
			vString, ok := v.(string)
			if !ok || !ifContainVar(vString) {
				putParsedValue(k, v)
				continue
			}

			// parse one value with "parsed map"
			valueParsed, err := template.NewRenderClient(&template.TemplateOption{},
				template.ContentGetter, template.AddDotForVariablesInConfigProcessor).
				Render(fmt.Sprintf("%v", v), parsed)
			// if no error(means this value doesn't import vars in "unparsed map"), put
			if err == nil {
				putParsedValue(k, valueParsed)
			}
		}
	}

	// check if vars map is correct
	if len(unparsed) > 0 {
		errString := "failed to parse var "
		var errKeyValues []string
		for k := range unparsed {
			errKeyValues = append(errKeyValues, fmt.Sprintf(`<"%s": "%s">`, k, origin[k]))
		}
		errString += strings.Join(errKeyValues, ", ")
		return nil, fmt.Errorf(errString)
	}

	return parsed, nil

}

// check if value contains other vars
func ifContainVar(value string) bool {
	// this function is humble, some case could not be checked probably
	// e.g. "[[123]]" "[[]]" will be regard as a var
	if strings.Contains(value, "[[") && strings.Contains(value, "]]") {
		return true
	}
	return false
}
