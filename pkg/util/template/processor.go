package template

import "regexp"

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

func AddDotForVariablesInConfigProcessor() Processor {
	return func(bytes []byte) ([]byte, error) {
		return []byte(addDotForVariablesInConfig(string(bytes))), nil
	}
}

func (r *rendererWithGetter) AddDotForVariablesInConfigProcessor() *rendererWithGetter {
	return r.AddProcessor(AddDotForVariablesInConfigProcessor())
}
