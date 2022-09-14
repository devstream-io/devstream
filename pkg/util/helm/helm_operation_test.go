package helm

import (
	"testing"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/repo"
)

func TestInstallOrUpgradeChart(t *testing.T) {
	atomic := true
	if !*helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	require.NoErrorf(t, err, "error: %v must be nil\n", err)
	chartSpec := &helmclient.ChartSpec{
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

	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// mockClient := mockhelmclient.NewMockClient(ctrl)
	// if mockClient == nil {
	// 	t.Fail()
	// }
	h, err := NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))

	require.NoErrorf(t, err, "error: %v\n", err)
	// mockClient.EXPECT().InstallOrUpgradeChart(context.TODO(), chartSpec).Return(&mockedRelease, nil)

	err = h.InstallOrUpgradeChart()
	require.NoError(t, err)
}

func TestAddOrUpdateChartRepo(t *testing.T) {
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
	chartSpec := &helmclient.ChartSpec{
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

	h, err := NewHelm(helmParam, WithEntry(entry), WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))
	require.NoErrorf(t, err, "error: %v\n", err)

	err = h.AddOrUpdateChartRepo(*entry)
	require.NoErrorf(t, err, "error: %v\n", err)
}

func TestHelm_UninstallHelmChartRelease(t *testing.T) {
	atomic := true
	if !*helmParam.Chart.Wait {
		atomic = false
	}
	tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
	require.NoErrorf(t, err, "error: %v must be nil\n", err)
	chartSpec := &helmclient.ChartSpec{
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
	// base
	h, err := NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient3{}))
	require.NoErrorf(t, err, "error: %v\n", err)

	err = h.UninstallHelmChartRelease()
	require.NoError(t, err)

	// mock error not found
	h, err = NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient{}))
	require.NoErrorf(t, err, "error: %v\n", err)

	err = h.UninstallHelmChartRelease()
	require.NoError(t, err)

	// mock error
	h, err = NewHelm(helmParam, WithChartSpec(chartSpec), WithClient(&DefaultMockClient2{}))
	require.NoErrorf(t, err, "error: %v\n", err)

	err = h.UninstallHelmChartRelease()
	require.Error(t, err, "error not found")

	require.Equalf(t, err, NormalError, "got: %+v\n, want %+v\n", err, NormalError)
}
