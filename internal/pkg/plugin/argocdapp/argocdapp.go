package argocdapp

import (
	_ "embed"
)

//go:embed tpl/argocd.tpl.yaml
var templateFileLoc string
