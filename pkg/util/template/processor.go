package template

import "regexp"

// this is because our variables' syntax is [[ varName ]]
// while Go's template is [[ .varName ]]
func AddDotForVariablesInConfig(s string) string {
	// regex := `\[\[\s*(.*)\s*\]\]`
	// r := regexp.MustCompile(regex)
	// return r.ReplaceAllString(s, "[[ .$1 ]]")
	regex := `\[\[\s*`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "[[ .")
}

type dotProcessor struct{}

func (p *dotProcessor) Process(bytes []byte) ([]byte, error) {
	return []byte(AddDotForVariablesInConfig(string(bytes))), nil
}

func AddDotForVariablesInConfigProcessor() Processor {
	return &dotProcessor{}
}

type processorFuncWrapper struct {
	f func([]byte) ([]byte, error)
}

func (p *processorFuncWrapper) Process(bytes []byte) ([]byte, error) {
	return p.f(bytes)
}

func NewProcessorFuncWrapper(f func([]byte) ([]byte, error)) Processor {
	return &processorFuncWrapper{f: f}
}
