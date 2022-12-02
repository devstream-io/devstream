package template

import "regexp"

// this is because our variables' syntax is [[ varName ]]
// while Go's template is [[ .varName ]]
func addDotForVariablesInConfig(s string) string {
	// add dot if there is only one word in [[ ]]
	regex := `\[\[\s*(\w+)\s*\]\]`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "[[ .$1 ]]")
}

func AddDotForVariablesInConfigProcessor() Processor {
	return func(bytes []byte) ([]byte, error) {
		return []byte(addDotForVariablesInConfig(string(bytes))), nil
	}
}

// When [[ ]] has two words and the first word don't contain quota, add quotes to the second word
// e.g. [[ env GITHUB_TOKEN]] -> [[ env "GITHUB_TOKEN" ]]
// [[ env 'GITHUB_TOKEN' ]] -> do nothing
// [[ env "GITHUB_TOKEN" ]] -> do nothing
// [[ "env" "GITHUB_TOKEN" ]] -> do nothing
// [[ GITHUB_TOKEN ]] -> do nothing
func addQuoteForVariablesInConfig(s string) string {
	regex := `\[\[\s*([^'"]\w+)\s+([^'"]\w+)\s*\]\]`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "[[ $1 \"$2\" ]]")
}

func AddQuoteForVariablesInConfigProcessor() Processor {
	return func(bytes []byte) ([]byte, error) {
		return []byte(addQuoteForVariablesInConfig(string(bytes))), nil
	}
}

func (r *rendererWithGetter) AddDotForVariablesInConfigProcessor() *rendererWithGetter {
	return r.AddProcessor(AddDotForVariablesInConfigProcessor())
}

func (r *rendererWithGetter) AddQuoteForVariablesInConfigProcessor() *rendererWithGetter {
	return r.AddProcessor(AddQuoteForVariablesInConfigProcessor())
}
