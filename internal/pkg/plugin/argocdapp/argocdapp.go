package argocdapp

import (
	_ "embed"
)

const Name = "argocdapp"

//go:embed tpl/argocd.tpl.yaml
var templateFileLoc string
