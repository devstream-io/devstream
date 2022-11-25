package golang

import _ "embed"

const Name = "githubactions-golang"

//go:embed pr.tpl.yml
var prPipeline string

//go:embed main.tpl.yml
var mainPipeline string
