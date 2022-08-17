package golang

import _ "embed"

//go:embed pr.tpl.yml
var prPipeline string

//go:embed main.tpl.yml
var mainPipeline string
