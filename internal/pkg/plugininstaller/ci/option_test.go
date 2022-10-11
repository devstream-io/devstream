package ci

import (
	"fmt"
	"net/http"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/server"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
)

var _ = Describe("Options struct", func() {
	Context("NewOptions method", func() {
		var (
			rawOptions plugininstaller.RawOptions
		)
		When("options is valid", func() {
			BeforeEach(func() {
				rawOptions = plugininstaller.RawOptions{
					"ci": map[string]interface{}{
						"type":    "gitlab",
						"content": "test",
					},
				}
			})
			It("should not raise error", func() {
				_, err := NewOptions(rawOptions)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("buildGitMap method", func() {
		var (
			opts *Options
			repo *common.Repo
		)
		BeforeEach(func() {
			opts = &Options{}
			repo = &common.Repo{
				Owner:    "test",
				Repo:     "test_repo",
				Branch:   "test_branch",
				RepoType: "gitlab",
			}
			opts.ProjectRepo = repo
		})
		When("all ci config is empty", func() {
			BeforeEach(func() {
				opts.CIConfig = &CIConfig{}
			})
			It("should raise error", func() {
				_, err := opts.buildGitMap()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("LocalPath field is setted with github ci files", func() {
			var (
				localDir, localFile string
				testContent         []byte
			)
			BeforeEach(func() {
				tempDir := GinkgoT().TempDir()
				localDir = fmt.Sprintf("%s/%s", tempDir, ".github/workflows")
				err := os.MkdirAll(localDir, os.ModePerm)
				Expect(err).Error().ShouldNot(HaveOccurred())
				tempFile, err := os.CreateTemp(localDir, "testFile")
				Expect(err).Error().ShouldNot(HaveOccurred())
				localFile = tempFile.Name()
				testContent = []byte("_test")
				err = os.WriteFile(localFile, testContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			When("LocalPath is directory", func() {
				BeforeEach(func() {
					opts.CIConfig = &CIConfig{
						Type:      "github",
						LocalPath: localDir,
					}
				})
				It("should get all files content", func() {
					gitMap, err := opts.buildGitMap()
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(gitMap).ShouldNot(BeEmpty())
					expectedKey := fmt.Sprintf("%s/%s", ".github/workflows", path.Base(localFile))
					v, ok := gitMap[expectedKey]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(Equal(testContent))
				})
			})
			When("localPath is file", func() {
				BeforeEach(func() {
					opts.CIConfig = &CIConfig{
						Type:      "github",
						LocalPath: localFile,
					}
				})
				It("should get file content", func() {
					gitMap, err := opts.buildGitMap()
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(gitMap).ShouldNot(BeEmpty())
					expectedKey := fmt.Sprintf("%s/%s", ".github/workflows", path.Base(localFile))
					v, ok := gitMap[expectedKey]
					Expect(ok).Should(BeTrue())
					Expect(v).Should(Equal(testContent))
				})
			})
		})
		When("Content field is setted", func() {
			var (
				testContent string
			)
			BeforeEach(func() {
				testContent = "testJenkins"
				opts.CIConfig = &CIConfig{
					Type:    "jenkins",
					Content: testContent,
				}
			})
			It("should return gitmap", func() {
				gitMap, err := opts.buildGitMap()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(len(gitMap)).Should(Equal(1))
				v, ok := gitMap["Jenkinsfile"]
				Expect(ok).Should(BeTrue())
				Expect(v).Should(Equal([]byte(testContent)))
			})
		})

		When("RemoteURL field is setted", func() {
			var (
				templateVal string
				s           *ghttp.Server
			)
			BeforeEach(func() {
				s = ghttp.NewServer()
				testContent := "testGitlabCI [[ .App ]]"
				templateVal = "template variable"
				s.RouteToHandler("GET", "/", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, testContent)
				})
				opts.CIConfig = &CIConfig{
					Type:      "gitlab",
					RemoteURL: s.URL(),
					Vars: map[string]interface{}{
						"App": templateVal,
					},
				}
			})
			AfterEach(func() {
				s.Close()
			})
			It("should get gitmap", func() {
				gitMap, err := opts.buildGitMap()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(len(gitMap)).Should(Equal(1))
				v, ok := gitMap[".gitlab-ci.yml"]
				Expect(ok).Should(BeTrue())
				Expect(string(v)).Should(Equal(fmt.Sprintf("testGitlabCI %s", templateVal)))
			})
		})
	})

	Context("FillDefaultValue method", func() {
		var (
			defaultOpts, opts *Options
		)
		BeforeEach(func() {
			opts = &Options{}
			defaultCIConfig := &CIConfig{
				Type:      "github",
				RemoteURL: "http://www.test.com",
			}
			defaultRepo := &common.Repo{
				Owner:    "test",
				Repo:     "test_repo",
				Branch:   "test_branch",
				RepoType: "gitlab",
			}
			defaultOpts = &Options{
				CIConfig:    defaultCIConfig,
				ProjectRepo: defaultRepo,
			}
		})
		When("ci config and repo are all empty", func() {
			It("should set default value", func() {
				opts.FillDefaultValue(defaultOpts)
				Expect(opts.CIConfig).ShouldNot(BeNil())
				Expect(opts.ProjectRepo).ShouldNot(BeNil())
				Expect(opts.CIConfig.RemoteURL).Should(Equal("http://www.test.com"))
				Expect(opts.ProjectRepo.Repo).Should(Equal("test_repo"))
			})
		})
		When("ci config and repo has some value", func() {
			BeforeEach(func() {
				opts.CIConfig = &CIConfig{
					RemoteURL: "http://exist.com",
				}
				opts.ProjectRepo = &common.Repo{
					Branch: "new_branch",
				}
			})
			It("should update empty value", func() {
				opts.FillDefaultValue(defaultOpts)
				Expect(opts.CIConfig).ShouldNot(BeNil())
				Expect(opts.ProjectRepo).ShouldNot(BeNil())
				Expect(opts.CIConfig.RemoteURL).Should(Equal("http://exist.com"))
				Expect(opts.CIConfig.Type).Should(Equal(server.CIServerType("github")))
				Expect(opts.ProjectRepo.Branch).Should(Equal("new_branch"))
				Expect(opts.ProjectRepo.Repo).Should(Equal("test_repo"))
			})
		})
	})
})
