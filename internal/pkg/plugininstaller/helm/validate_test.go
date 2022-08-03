package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var _ = Describe("Validate func", func() {
	var testOption plugininstaller.RawOptions

	When("options is not valid", func() {
		BeforeEach(func() {
			testOption = map[string]interface{}{
				"chart": map[string]string{},
				"repo":  map[string]string{},
			}
		})
		It("should return error", func() {
			_, err := helm.Validate(testOption)
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("options is valid", func() {
		BeforeEach(func() {
			testOption = map[string]interface{}{
				"chart": map[string]string{
					"chart_name": "test",
				},
				"repo": map[string]string{
					"url":  "http://test.com",
					"name": "test",
				},
			}
		})
		It("should return success", func() {
			opt, err := helm.Validate(testOption)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(opt).ShouldNot(BeEmpty())
		})
	})
})

var _ = Describe("SetDefaultConfig func", func() {
	var (
		testChartName string
		testRepoURL   string
		testRepoName  string
		testBool      *bool
		defaultConfig helm.Options
		testOptions   plugininstaller.RawOptions
		expectChart   map[string]interface{}
		expectRepo    map[string]interface{}
	)
	BeforeEach(func() {
		testChartName = "test_chart"
		testRepoName = "test_repo"
		testRepoURL = "http://test.com"
		testBool = types.Bool(true)
		testOptions = map[string]interface{}{
			"chart": map[string]string{},
			"repo":  map[string]string{},
		}
		defaultConfig = helm.Options{
			Chart: helmCommon.Chart{
				ChartName:       testChartName,
				Wait:            testBool,
				UpgradeCRDs:     testBool,
				CreateNamespace: testBool,
			},
			Repo: helmCommon.Repo{
				URL:  testRepoURL,
				Name: testRepoName,
			},
		}
		expectChart = map[string]interface{}{
			"chart_name":       testChartName,
			"wait":             testBool,
			"namespace":        "",
			"version":          "",
			"release_name":     "",
			"values_yaml":      "",
			"timeout":          "",
			"create_namespace": testBool,
			"upgradeCRDs":      testBool,
		}
		expectRepo = map[string]interface{}{
			"url":  testRepoURL,
			"name": testRepoName,
		}
	})
	It("should update default value", func() {
		updateFunc := helm.SetDefaultConfig(&defaultConfig)
		o, err := updateFunc(testOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		oRepo, exist := o["repo"]
		Expect(exist).Should(BeTrue())
		oChart, exist := o["chart"]
		Expect(exist).Should(BeTrue())
		Expect(oRepo).Should(Equal(expectRepo))
		Expect(oChart).Should(Equal(expectChart))
	})
})
