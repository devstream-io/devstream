package tool

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/k8s"
)

const (
	argocdRepoName         = "argo"
	argocdRepoURL          = "https://argoproj.github.io/argo-helm"
	argocdNamespace        = "argocd"
	argocdChartReleaseName = "argocd"
	argocdChartName        = argocdRepoName + "/" + "argo-cd"
)

var toolArgocd = tool{
	Name: "Argo CD",
	IfExists: func() bool {
		cmd := exec.Command("helm", "status", argocdChartReleaseName, "-n", argocdNamespace)
		return cmd.Run() == nil
	},

	Install: func() error {
		if !confirm("Argo CD") {
			return fmt.Errorf("user cancelled")
		}

		// create namespace if not exist
		kubeClient, err := k8s.NewClient()
		if err != nil {
			return err
		}
		if err = kubeClient.UpsertNameSpace(argocdNamespace); err != nil {
			return err
		}

		// install argocd by helm
		err = execCommand([]string{"helm", "repo", "add", argocdRepoName, argocdRepoURL})
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			return err
		}

		if err = execCommand([]string{"helm", "install", argocdChartReleaseName, argocdChartName, "-n", argocdNamespace, "--wait", "--create-namespace"}); err != nil {
			return err
		}

		return nil
	},
}
