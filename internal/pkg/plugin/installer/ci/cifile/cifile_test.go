package cifile

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/downloader"
)

var _ = Describe("Options struct", func() {
	var (
		c                                 *CIFileConfig
		ciFilePath, ciType, ciFileContent string
	)
	BeforeEach(func() {
		ciType = "jenkins"
		ciFileContent = "test_content"
		c = &CIFileConfig{
			Type: server.CIServerType(ciType),
		}
		ciFilePath = c.newCIServerClient().CIFilePath()
	})
	Context("SetContent method", func() {
		When("contentMap is existed", func() {
			BeforeEach(func() {
				c.SetContent(ciFileContent)
			})
			It("should update map", func() {
				existV, ok := c.ConfigContentMap[ciFilePath]
				Expect(ok).Should(BeTrue())
				Expect(existV).Should(Equal(ciFileContent))
				testContent := "another_test"
				c.SetContent(testContent)
				changedV, ok := c.ConfigContentMap[ciFilePath]
				Expect(ok).Should(BeTrue())
				Expect(changedV).Should(Equal(testContent))
			})
		})
	})

	Context("SetContentMap method", func() {
		It("should set contentMap", func() {
			testMap := map[string]string{
				"test": "gg",
			}
			c.SetContentMap(testMap)
			v, ok := c.ConfigContentMap["test"]
			Expect(ok).Should(BeTrue())
			Expect(v).Should(Equal("gg"))
			Expect(len(c.ConfigContentMap)).Should(Equal(1))
		})
	})

	Context("getGitfileMap method", func() {
		When("ConfigContentMap is not empty", func() {
			When("render contentMap failed", func() {
				BeforeEach(func() {
					ciFileContent = "rendered_test[[ .APP ]]"
					c.SetContent(ciFileContent)
					c.Vars = map[string]interface{}{
						"GG": "not_exist",
					}
				})
				It("should return render err", func() {
					_, err := c.getGitfileMap()
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("map has no entry for key"))
				})
			})
			When("not render map content", func() {
				BeforeEach(func() {
					c.SetContent(ciFileContent)
					c.Vars = map[string]interface{}{}
				})
				It("should return render err", func() {
					gitFileMap, err := c.getGitfileMap()
					Expect(err).Error().ShouldNot(HaveOccurred())
					v, ok := gitFileMap[ciFilePath]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(Equal([]byte(v)))
				})
			})
		})
		When("content and location are all empty", func() {
			BeforeEach(func() {
				c.ConfigContentMap = map[string]string{}
				c.ConfigLocation = ""
			})
			It("should return error", func() {
				_, err := c.getGitfileMap()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("ci can't get valid ci files"))
			})
		})
	})

	Context("renderContent method", func() {
		BeforeEach(func() {
			c.Vars = map[string]interface{}{
				"APP": "here",
			}
			ciFileContent = "test,[[ .APP ]]"
		})
		It("should work normal", func() {
			content, err := c.renderContent(ciFileContent)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(content).Should(Equal("test,here"))
		})
	})

	Context("getConfigContentFromLocation method", func() {
		var localFilePath string
		When("get ci file failed", func() {
			BeforeEach(func() {
				localFilePath = "not_exist"
				c.ConfigContentMap = map[string]string{}
				c.ConfigLocation = downloader.ResourceLocation(localFilePath)
			})
			It("should return error", func() {
				_, err := c.getConfigContentFromLocation()
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
				c.ConfigContentMap = map[string]string{}
				c.ConfigLocation = downloader.ResourceLocation(localFilePath.Name())
			})
			It("should return error", func() {
				gitMap, err := c.getConfigContentFromLocation()
				Expect(err).Error().ShouldNot(HaveOccurred())
				v, ok := gitMap[c.newCIServerClient().CIFilePath()]
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
