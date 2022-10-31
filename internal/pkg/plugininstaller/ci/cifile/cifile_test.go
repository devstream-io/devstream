package cifile

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile/server"
)

var _ = Describe("Options struct", func() {
	var (
		ciConfig                          *CIConfig
		ciFilePath, ciType, ciFileContent string
	)
	BeforeEach(func() {
		ciType = "jenkins"
		ciFileContent = "test_content"
		ciConfig = &CIConfig{
			Type: server.CIServerType(ciType),
		}
		ciFilePath = ciConfig.newCIServerClient().CIFilePath()
	})
	Context("SetContent method", func() {
		When("contentMap is existed", func() {
			BeforeEach(func() {
				ciConfig.SetContent(ciFileContent)
			})
			It("should update map", func() {
				existV, ok := ciConfig.ConfigContentMap[ciFilePath]
				Expect(ok).Should(BeTrue())
				Expect(existV).Should(Equal(ciFileContent))
				testContent := "another_test"
				ciConfig.SetContent(testContent)
				changedV, ok := ciConfig.ConfigContentMap[ciFilePath]
				Expect(ok).Should(BeTrue())
				Expect(changedV).Should(Equal(testContent))
			})
		})
	})

	Context("SetContentMap method", func() {
		It("should set contentMap", func() {
			c := map[string]string{
				"test": "gg",
			}
			ciConfig.SetContentMap(c)
			v, ok := ciConfig.ConfigContentMap["test"]
			Expect(ok).Should(BeTrue())
			Expect(v).Should(Equal("gg"))
			Expect(len(ciConfig.ConfigContentMap)).Should(Equal(1))
		})
	})

	Context("getGitfileMap method", func() {
		When("ConfigContentMap is not empty", func() {
			When("render contentMap failed", func() {
				BeforeEach(func() {
					ciFileContent = "rendered_test[[ .APP ]]"
					ciConfig.SetContent(ciFileContent)
					ciConfig.Vars = map[string]interface{}{
						"GG": "not_exist",
					}
				})
				It("should return render err", func() {
					_, err := ciConfig.getGitfileMap()
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("map has no entry for key"))
				})
			})
			When("not render map content", func() {
				BeforeEach(func() {
					ciConfig.SetContent(ciFileContent)
					ciConfig.Vars = map[string]interface{}{}
				})
				It("should return render err", func() {
					gitFileMap, err := ciConfig.getGitfileMap()
					Expect(err).Error().ShouldNot(HaveOccurred())
					v, ok := gitFileMap[ciFilePath]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(Equal([]byte(v)))
				})
			})
		})
		When("content and location are all empty", func() {
			BeforeEach(func() {
				ciConfig.ConfigContentMap = map[string]string{}
				ciConfig.ConfigLocation = ""
			})
			It("should return error", func() {
				_, err := ciConfig.getGitfileMap()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("ci can't get valid ci files"))
			})
		})
	})

	Context("renderContent method", func() {
		BeforeEach(func() {
			ciConfig.Vars = map[string]interface{}{
				"APP": "here",
			}
			ciFileContent = "test,[[ .APP ]]"
		})
		It("should work normal", func() {
			content, err := ciConfig.renderContent(ciFileContent)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(content).Should(Equal("test,here"))
		})
	})

	Context("getConfigContentFromLocation method", func() {
		var localFilePath string
		When("get ci file failed", func() {
			BeforeEach(func() {
				localFilePath = "not_exist"
				ciConfig.ConfigContentMap = map[string]string{}
				ciConfig.ConfigLocation = localFilePath
			})
			It("should return error", func() {
				_, err := ciConfig.getConfigContentFromLocation()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("no such file or directory"))
			})
		})
		When("get ci file success", func() {
			BeforeEach(func() {
				ciFileContent = "get ci file content"
				localFilePath, err := os.CreateTemp("", "test")
				Expect(err).ShouldNot(HaveOccurred())
				testContent := []byte(ciFileContent)
				err = os.WriteFile(localFilePath.Name(), testContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
				ciConfig.ConfigContentMap = map[string]string{}
				ciConfig.ConfigLocation = localFilePath.Name()
			})
			It("should return error", func() {
				gitMap, err := ciConfig.getConfigContentFromLocation()
				Expect(err).Error().ShouldNot(HaveOccurred())
				v, ok := gitMap[ciConfig.newCIServerClient().CIFilePath()]
				Expect(ok).Should(BeTrue())
				Expect(v).Should(Equal([]byte(ciFileContent)))
			})
		})
	})
})

var _ = Describe("CIFileVarsMap struct", func() {
	var m CIFileVarsMap
	Context("Set method", func() {
		m = CIFileVarsMap{}
		k := "test_key"
		v := "test_val"
		m.Set(k, v)
		expectV, ok := m[k]
		Expect(ok).Should(BeTrue())
		Expect(expectV).Should(Equal(v))
	})
})
