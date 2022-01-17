package helm

import (
	"context"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"log"
	"strings"
	"time"
)

type Helm struct {
	*repo.Entry
	*helmclient.ChartSpec
	helmclient.Client
}

func NewHelm(param *HelmParam) (*Helm, error) {
	var hClient helmclient.Client
	var err error
	if hClient, err = helmclient.New(
		&helmclient.Options{
			Namespace:        param.Chart.Namespace,
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
		},
	); err != nil {
		return nil, err
	}

	entry := &repo.Entry{
		Name:                  param.Repo.Name,
		URL:                   param.Repo.URL,
		Username:              "",
		Password:              "",
		CertFile:              "",
		KeyFile:               "",
		CAFile:                "",
		InsecureSkipTLSverify: false,
		PassCredentialsAll:    false,
	}

	chartSpec := &helmclient.ChartSpec{
		ReleaseName:     param.Chart.ReleaseName,
		ChartName:       param.Chart.Name,
		Namespace:       param.Chart.Namespace,
		ValuesYaml:      "",
		Version:         "",
		CreateNamespace: param.Chart.CreateNamespace,
		DisableHooks:    false,
		Replace:         false,
		// TODO(daniel-hutao): default to true now, maybe exposed at config.yaml later
		Wait:             true,
		DependencyUpdate: false,
		// TODO(daniel-hutao): default to 5min now, maybe exposed at config.yaml later
		Timeout:      5 * time.Minute,
		GenerateName: false,
		NameTemplate: "",
		Atomic:       false,
		SkipCRDs:     false,
		// TODO(daniel-hutao): default to true now, maybe exposed at config.yaml later
		UpgradeCRDs:   true,
		SubNotes:      false,
		Force:         false,
		ResetValues:   false,
		ReuseValues:   false,
		Recreate:      false,
		MaxHistory:    0,
		CleanupOnFail: false,
		DryRun:        false,
	}

	helm := &Helm{
		Entry:     entry,
		ChartSpec: chartSpec,
		Client:    hClient,
	}

	if err = helm.AddOrUpdateChartRepo(*helm.Entry); err != nil {
		return nil, err
	}

	return helm, nil
}

func (h *Helm) InstallOrUpgradeChart() error {
	_, err := h.Client.InstallOrUpgradeChart(context.TODO(), h.ChartSpec)
	return err
}

func (h *Helm) UninstallHelmChart() error {
	var err error
	if err = h.UninstallReleaseByName(h.ChartSpec.ReleaseName); err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Println("release is not found, maybe it has been deleted")
			return nil
		}
		return err
	}

	return nil
}
