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

type Option func(*Helm)

func NewHelm(param *HelmParam, option ...Option) (*Helm, error) {
	hClient, err := helmclient.New(
		&helmclient.Options{
			Namespace:        param.Chart.Namespace,
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
		},
	)
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
	atomic := true
	if !param.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(param.Chart.Timeout)
	if err != nil {
		return nil, err
	}
	spec := &helmclient.ChartSpec{
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
	h := &Helm{
		Entry:     entry,
		ChartSpec: spec,
		Client:    hClient,
	}

	for _, op := range option {
		op(h)
	}

	if err = h.AddOrUpdateChartRepo(*entry); err != nil {
		return nil, err
	}
	return h, nil
}

func WithEntry(entry *repo.Entry) Option {
	return func(r *Helm) {
		r.Entry = entry
	}
}

func WithChartSpec(spec *helmclient.ChartSpec) Option {
	return func(r *Helm) {
		r.ChartSpec = spec
	}
}

func WithClient(client helmclient.Client) Option {
	return func(r *Helm) {
		r.Client = client
	}
}

func (h *Helm) AddOrUpdateChartRepo(entry repo.Entry) error {
	return h.Client.AddOrUpdateChartRepo(entry)
}

func (h *Helm) InstallOrUpgradeChart() error {
	_, err := h.Client.InstallOrUpgradeChart(context.TODO(), h.ChartSpec)
	return err
}

func (h *Helm) UninstallHelmChartRelease() (err error) {
	if err = h.Client.UninstallReleaseByName(h.ChartSpec.ReleaseName); err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Warn("Release is not found, maybe it has been deleted.")
			return nil
		}
		return err
	}
	return nil
}

// GetAnnotationName will return label key for service created by helm
func GetAnnotationName() string {
	return "meta.helm.sh/release-name"
}

// GetAnnotationName will return label key for service created by helm
func GetLabelName() string {
	return "app.kubernetes.io/instance"
}
