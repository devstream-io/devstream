package tool

import "fmt"

var toolArgocd = tool{
	Name: "Argo CD",
	IfExists: func() bool {
		// TODO(dh)
		return false
	},

	Install: func() error {
		if !confirm("Argo CD") {
			return fmt.Errorf("user cancelled")
		}

		if err := execCommand([]string{"helm", "repo", "add", "argo", "https://argoproj.github.io/argo-helm"}); err != nil {
			return err
		}
		if err := execCommand([]string{"helm", "install", "argo/argo-cd", "-n", "argocd", "--create-namespace"}); err != nil {
			return err
		}
		return nil
	},
}
