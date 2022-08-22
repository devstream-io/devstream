package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
)

var _ = Describe("Options struct", func() {
	var (
		testOpts      helm.Options
		testChartName string
		testRepoName  string
		testNameSpace string
		expectMap     map[string]interface{}
		emptyBool     *bool
	)

	BeforeEach(func() {
		testChartName = "test_chart"
		testRepoName = "test_repo"
		testNameSpace = "test_nameSpace"
		testOpts = helm.Options{
			Chart: helmCommon.Chart{
				ChartName: testChartName,
				Namespace: testNameSpace,
			},
			Repo: helmCommon.Repo{
				Name: testRepoName,
			},
		}
		expectMap = map[string]interface{}{
			"repo": map[string]interface{}{
				"name": "test_repo",
				"url":  "",
			},
			"chart": map[string]interface{}{
				"version":      "",
				"release_name": "",
				"wait":         emptyBool,
				"chart_name":   "test_chart",
				"namespace":    "test_nameSpace",
				"timeout":      "",
				"upgradeCRDs":  emptyBool,
				"values_yaml":  "",
			},
		}
	})

	Context("GetHelmParam method", func() {
		It("should pass chart and repo field", func() {
			helmParam := testOpts.GetHelmParam()
			Expect(helmParam.Chart).Should(Equal(testOpts.Chart))
			Expect(helmParam.Repo).Should(Equal(testOpts.Repo))
		})
	})

	Context("GetNamespace method", func() {
		It("should return chart's nameSpace", func() {
			Expect(testOpts.GetNamespace()).Should(Equal(testOpts.Chart.Namespace))
		})
	})

	Context("GetReleaseName method", func() {
		It("should return chart's ReleaseName", func() {
			Expect(testOpts.GetReleaseName()).Should(Equal(testOpts.Chart.ReleaseName))
		})
	})

	Context("Encode method", func() {
		It("should return opts map", func() {
			result, err := testOpts.Encode()
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(expectMap))
		})
	})
})

var _ = Describe("NewOptions func", func() {

	var (
		inputOptions  plugininstaller.RawOptions
		testRepoName  string
		testChartName string
	)

	BeforeEach(func() {
		testRepoName = "test_repo"
		testChartName = "test_chart"
		inputOptions = map[string]interface{}{
			"repo": map[string]interface{}{
				"name": testRepoName,
			},
			"chart": map[string]interface{}{
				"chart_name": testChartName,
			},
		}
	})

	It("should work normal", func() {
		opts, err := helm.NewOptions(inputOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(opts.Chart.ChartName).Should(Equal(testChartName))
		Expect(opts.Repo.Name).Should(Equal(testRepoName))
	})
})
