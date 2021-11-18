package argocd

import (
	"context"
	"log"
	"time"

	helmClient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

func createHelmClient(namespace string) helmClient.Client {
	client, err := helmClient.New(
		&helmClient.Options{
			Namespace:        namespace,
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
		},
	)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return client
}

func addArgoHelmRepo(c *helmClient.Client, param *Param) {
	chartRepo := repo.Entry{
		Name: param.Repo.Name,
		URL:  param.Repo.URL,
	}

	if err := (*c).AddOrUpdateChartRepo(chartRepo); err != nil {
		log.Fatalf(err.Error())
	}
}

func installOrUpdateArgoHelmChart(c *helmClient.Client, param *Param) {
	chartSpec := helmClient.ChartSpec{
		ReleaseName: param.Chart.ReleaseName,
		ChartName:   param.Chart.Name,
		Namespace:   param.Chart.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
		Timeout:     3 * time.Minute,
	}
	if _, err := (*c).InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		log.Fatalf(err.Error())
	}
}
