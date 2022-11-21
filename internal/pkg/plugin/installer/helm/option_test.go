package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
)

var _ = Describe("Options struct", func() {
	var (
		testOpts      helm.Options
		testChartName string
		testRepoName  string
		testNameSpace string
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
})

var _ = Describe("NewOptions func", func() {

	var (
		inputOptions  configmanager.RawOptions
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
				"chartName": testChartName,
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
