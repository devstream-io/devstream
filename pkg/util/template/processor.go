package template

import (
	"regexp"
)

// Processor process content before render
type processFunc func([]byte) []byte

// AddDotForVariablesInConfigProcessor will add dot before varName
// this is because our variables' syntax is [[ varName ]]
// while Go's template is [[ .varName ]]
func AddDotForVariablesInConfigProcessor(bytes []byte) []byte {
	regex := `\[\[\s*(\w+)\s*\]\]`
	r := regexp.MustCompile(regex)
	return r.ReplaceAll(bytes, []byte("[[ .$1 ]]"))
}

// AddQuoteForVariablesInConfigProcessor will add quote for special varName
// When [[ ]] has two words and the first word don't contain quota, add quotes to the second word
// e.g. [[ env GITHUB_TOKEN]] -> [[ env "GITHUB_TOKEN" ]]
// [[ env 'GITHUB_TOKEN' ]] -> do nothing
// [[ env "GITHUB_TOKEN" ]] -> do nothing
// [[ "env" "GITHUB_TOKEN" ]] -> do nothing
// [[ GITHUB_TOKEN ]] -> do nothing
func AddQuoteForVariablesInConfigProcessor(bytes []byte) []byte {
	regex := `\[\[\s*([^'"]\w+)\s+([^'"]\w+)\s*\]\]`
	r := regexp.MustCompile(regex)
	return r.ReplaceAll(bytes, []byte("[[ $1 \"$2\" ]]"))
}
