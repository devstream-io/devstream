package helm

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	repositoryCache  = filepath.Join(os.TempDir(), ".helmcache")
	repositoryConfig = filepath.Join(os.TempDir(), ".helmrepo")
)

// Helm is helm implementation
type Helm struct {
	*repo.Entry
	*helmclient.ChartSpec
	helmclient.Client
}

type Option func(*Helm)

// NewHelm creates a new Helm
func NewHelm(param *HelmParam, option ...Option) (*Helm, error) {
	isDebugMode := viper.GetBool("debug")
	if isDebugMode {
		log.Info("Helm is running in debug mode.")
	}

	hClient, err := helmclient.New(
		&helmclient.Options{
			Namespace:        param.Chart.Namespace,
			RepositoryCache:  repositoryCache,
			RepositoryConfig: repositoryConfig,
			Debug:            isDebugMode,
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
	if !*param.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(param.Chart.Timeout)
	if err != nil {
		return nil, err
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
		Wait:             *param.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      *param.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}
	if param.Chart.ChartPath != "" {
		chartSpec.ChartName = param.Chart.ChartPath
		if err = cacheChartPackage(param.Chart.ChartPath); err != nil {
			return nil, err
		}
	}

	h := &Helm{
		Entry:     entry,
		ChartSpec: chartSpec,
		Client:    hClient,
	}

	for _, op := range option {
		op(h)
	}

	if param.Chart.ChartPath == "" {
		if err = h.AddOrUpdateChartRepo(*entry); err != nil {
			return nil, err
		}
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

// GetLabelName will return label key for service created by helm
func GetLabelName() string {
	return "app.kubernetes.io/instance"
}
