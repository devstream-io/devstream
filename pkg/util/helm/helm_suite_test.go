package helm_test

import (
	"context"
	"testing"

	helmclient "github.com/mittwald/go-helm-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

func TestHelm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg Util Helm Test Suite")
}

type mockClient struct {
	helmclient.Client
	AddOrUpdateChartRepoError   error
	UninstallReleaseByNameError error
	InstallOrUpgradeChartError  error
}

func (c *mockClient) AddOrUpdateChartRepo(enrty repo.Entry) error {

	return c.AddOrUpdateChartRepoError
}

func (c *mockClient) UninstallReleaseByName(name string) error {
	return c.UninstallReleaseByNameError
}

func (c *mockClient) InstallOrUpgradeChart(ctx context.Context, spec *helmclient.ChartSpec) (*release.Release, error) {
	if c.InstallOrUpgradeChartError != nil {
		return nil, c.InstallOrUpgradeChartError
	}
	var mockedRelease = release.Release{Name: "test"}
	return &mockedRelease, nil
}

var helmParam = &helm.HelmParam{
	helm.Repo{
		Name: "helm",
		URL:  "test1",
	},
	helm.Chart{
		ReleaseName: "helm:v1.0.0",
		Timeout:     "1m",
		Wait:        types.Bool(false),
		UpgradeCRDs: types.Bool(false),
	},
}
