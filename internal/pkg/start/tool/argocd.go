package tool

import (
	"fmt"
	"os/exec"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	helmUtil "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolArgocd = tool{
	Name: "Argo CD",
	IfExists: func() bool {
		cmd := exec.Command("helm", "status", "argocd", "-n", "argocd")
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
		if err = kubeClient.UpsertNameSpace("argocd"); err != nil {
			return err
		}

		// install argocd by helm
		argocdHelmOpts := &helm.Options{
			Chart: helmCommon.Chart{
				ChartPath:   "",
				ChartName:   "argo/argo-cd",
				Version:     "",
				Timeout:     "10m",
				Wait:        types.Bool(true),
				UpgradeCRDs: types.Bool(true),
				ReleaseName: "argocd",
				Namespace:   "argocd",
			},
			Repo: helmCommon.Repo{
				URL:  "https://argoproj.github.io/argo-helm",
				Name: "argo",
			},
		}
		h, err := helmUtil.NewHelm(argocdHelmOpts.GetHelmParam())
		if err != nil {
			return err
		}
		if err = h.InstallOrUpgradeChart(); err != nil {
			return err
		}

		return nil
	},
}
