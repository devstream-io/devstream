package argocd

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	helmClient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

type ArgoCD struct {
	client *helmClient.Client
	param  *Param
}

func NewArgoCD(options *map[string]interface{}) (*ArgoCD, error) {
	var param Param
	if err := mapstructure.Decode(*options, &param); err != nil {
		return nil, err
	}

	client, err := helmClient.New(
		&helmClient.Options{
			Namespace:        param.Chart.Namespace,
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
		},
	)
	if err != nil {
		return nil, err
	}

	return &ArgoCD{
		client: &client,
		param:  &param,
	}, nil
}

func (a *ArgoCD) addHelmRepo() error {
	chartRepo := repo.Entry{
		Name: a.param.Repo.Name,
		URL:  a.param.Repo.URL,
	}

	if err := (*a.client).AddOrUpdateChartRepo(chartRepo); err != nil {
		return err
	}
	return nil
}

func (a *ArgoCD) installOrUpgradeHelmChart() error {
	log.Println("Adding and updating argocd helm chart repo ...")
	if err := a.addHelmRepo(); err != nil {
		return err
	}

	chartSpec := helmClient.ChartSpec{
		ReleaseName:     a.param.Chart.ReleaseName,
		ChartName:       a.param.Chart.Name,
		Namespace:       a.param.Chart.Namespace,
		CreateNamespace: a.param.Chart.CreateNamespace,
		UpgradeCRDs:     true,
		Wait:            true,
		Timeout:         3 * time.Minute,
	}

	_, err := (*a.client).InstallOrUpgradeChart(context.Background(), &chartSpec)
	if err != nil {
		return err
	}

	return nil
}

// uninstallHelmChartIgnoreReleaseNotFound will return nil when:
// 1. The argocd helm chart release uninstall successful
// 2. The argocd helm chart release not found
func (a *ArgoCD) uninstallHelmChartIgnoreReleaseNotFound() error {
	err := a.uninstallHelmChart()
	// Log: < Release not loaded: argocd: release: not found >
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "not found") {
		log.Println("argocd release is not found, maybe it has been deleted")
		return nil
	}
	return err
}

func (a *ArgoCD) uninstallHelmChart() error {
	err := (*a.client).UninstallReleaseByName(a.param.Chart.ReleaseName)
	if err != nil {
		return err
	}
	return nil
}
