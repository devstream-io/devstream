package helm

import (
	"context"
	"errors"
	"testing"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/types"
)

var (
	NormalError   = errors.New("normal error")
	NotFoundError = errors.New("release name not found")
)

var mockedRelease = release.Release{Name: "test"}
var helmParam = &HelmParam{
	Repo{
		Name: "helm",
		URL:  "test1",
	},
	Chart{
		ReleaseName: "helm:v1.0.0",
		Timeout:     "1m",
		Wait:        types.Bool(false),
		UpgradeCRDs: types.Bool(false),
	},
}

type DefaultMockClient struct {
	helmclient.Client
}

func (c *DefaultMockClient) AddOrUpdateChartRepo(enrty repo.Entry) error {
	return nil
}

func (c *DefaultMockClient) UninstallReleaseByName(name string) error {
	return NotFoundError
}

func (c *DefaultMockClient) InstallOrUpgradeChart(ctx context.Context, spec *helmclient.ChartSpec) (*release.Release, error) {
	return &mockedRelease, nil
}

type DefaultMockClient2 struct {
	DefaultMockClient
}

func (c *DefaultMockClient2) UninstallReleaseByName(name string) error {
	return NormalError
}

type DefaultMockClient3 struct {
	DefaultMockClient
}

func (c *DefaultMockClient3) UninstallReleaseByName(name string) error {
	return nil
}

type DefaultMockClient4 struct {
	DefaultMockClient
}

func (c *DefaultMockClient4) AddOrUpdateChartRepo(enrty repo.Entry) error {
	return NormalError
}

func TestNewHelm(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    bool
		wantHelm   bool
		helmClient helmclient.Client
	}{
		{
			name:       "base",
			wantErr:    false,
			wantHelm:   true,
			helmClient: &DefaultMockClient{},
		},
		{
			name:       "newHelm with NormalError",
			wantErr:    true,
			wantHelm:   false,
			helmClient: &DefaultMockClient4{},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHelm(helmParam, WithClient(tt.helmClient))

			if tt.wantErr {
				require.Errorf(t, err, "error: %v must not be nil\n", err)
			} else {
				require.NoErrorf(t, err, "error: %v must be nil\n", err)
			}

			if tt.wantHelm {
				require.NotNilf(t, got, "got: %v must not be nil\n", got)
			} else {
				require.Nilf(t, got, "got: %v must be nil\n", got)
			}
		})

	}

}

func TestNewHelmWithOption(t *testing.T) {
	entry := &repo.Entry{
		Name:                  helmParam.Repo.Name,
		URL:                   helmParam.Repo.URL,
		Username:              "",
		Password:              "",
		CertFile:              "",
		KeyFile:               "",
		CAFile:                "",
		InsecureSkipTLSverify: false,
		PassCredentialsAll:    false,
	}
	atomic := true
	if !*helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	require.NoErrorf(t, err, "error: %v must be nil\n", err)
	spec := &helmclient.ChartSpec{
		ReleaseName:      helmParam.Chart.ReleaseName,
		ChartName:        helmParam.Chart.ChartName,
		Namespace:        helmParam.Chart.Namespace,
		ValuesYaml:       helmParam.Chart.ValuesYaml,
		Version:          helmParam.Chart.Version,
		CreateNamespace:  false,
		DisableHooks:     false,
		Replace:          true,
		Wait:             *helmParam.Chart.Wait,
		DependencyUpdate: false,
		Timeout:          tmout,
		GenerateName:     false,
		NameTemplate:     "",
		Atomic:           atomic,
		SkipCRDs:         false,
		UpgradeCRDs:      *helmParam.Chart.UpgradeCRDs,
		SubNotes:         false,
		Force:            false,
		ResetValues:      false,
		ReuseValues:      false,
		Recreate:         false,
		MaxHistory:       0,
		CleanupOnFail:    false,
		DryRun:           false,
	}

	mockClient := &DefaultMockClient{}

	got, err := NewHelm(helmParam, WithClient(mockClient))
	require.NoErrorf(t, err, "err: %v\n", err)

	want := &Helm{
		Entry:     entry,
		ChartSpec: spec,
		Client:    mockClient,
	}
	require.Equalf(t, got, want, "NewHelm() = \n%+v\n, want \n%+v\n", got, want)
}
