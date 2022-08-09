package python

import _ "embed"

//go:embed lint.yaml
var lintPipeline string

//go:embed test.yaml
var testPipeline string

//go:embed docker.yaml
var dockerPipeline string
