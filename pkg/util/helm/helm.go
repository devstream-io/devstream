package helm

import (
	"context"
	"strings"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/log"
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

	tmout, err := time.ParseDuration(param.Chart.Timeout)
	if err != nil {
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

	// 'Wait' will automatically be set to true when using Atomic.
	atomic := true
	if !param.Chart.Wait {
		atomic = false
	}

	chartSpec := &helmclient.ChartSpec{
		ReleaseName:      param.Chart.ReleaseName,
		ChartName:        param.Chart.ChartName,
		Namespace:        param.Chart.Namespace,
		ValuesYaml:       param.Chart.ValuesYaml,
		Version:          param.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             param.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      param.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
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

func (h *Helm) UninstallHelmChartRelease() error {
	var err error
	if err = h.UninstallReleaseByName(h.ChartSpec.ReleaseName); err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Warn("Release is not found, maybe it has been deleted.")
			return nil
		}
		return err
	}

	return nil
}
