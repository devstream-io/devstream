package helm_test

import (
	"fmt"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/devstream-io/devstream/pkg/util/helm"
)

var _ = Describe("InstallOrUpgradeChart func", func() {
	It("should work noraml", func() {
		atomic := true
		if !*helmParam.Chart.Wait {
			atomic = false
		}
		tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
		Expect(err).ShouldNot(HaveOccurred())
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

		h, err := helm.NewHelm(helmParam, helm.WithChartSpec(chartSpec), helm.WithClient(&mockClient{}))
		Expect(err).ShouldNot(HaveOccurred())
		err = h.InstallOrUpgradeChart()
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("AddOrUpdateChartRepo func", func() {
	It("should work noraml", func() {
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
		Expect(err).ShouldNot(HaveOccurred())
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

		h, err := helm.NewHelm(helmParam, helm.WithEntry(entry), helm.WithChartSpec(chartSpec), helm.WithClient(&mockClient{}))
		Expect(err).ShouldNot(HaveOccurred())
		err = h.AddOrUpdateChartRepo(*entry)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("UninstallHelmChartRelease func", func() {
	It("should work", func() {
		atomic := true
		if !*helmParam.Chart.Wait {
			atomic = false
		}
		tmout, err := time.ParseDuration(helmParam.Chart.Timeout)
		Expect(err).ShouldNot(HaveOccurred())
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
		h, err := helm.NewHelm(helmParam, helm.WithChartSpec(chartSpec), helm.WithClient(&mockClient{}))
		Expect(err).ShouldNot(HaveOccurred())

		err = h.UninstallHelmChartRelease()
		Expect(err).ShouldNot(HaveOccurred())

		// mock error
		h, err = helm.NewHelm(helmParam, helm.WithChartSpec(chartSpec), helm.WithClient(&mockClient{
			UninstallReleaseByNameError: fmt.Errorf("data error"),
		}))
		Expect(err).ShouldNot(HaveOccurred())

		err = h.UninstallHelmChartRelease()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("data error"))
	})
})
