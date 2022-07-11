package jenkins

import _ "embed"

//go:embed markupformatter.yaml
var markupformatterYaml string

func init() {
	_ = markupformatterYaml
}
