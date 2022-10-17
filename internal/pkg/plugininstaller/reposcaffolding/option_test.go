package reposcaffolding

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("Options struct", func() {
	var (
		opts *Options
	)
	BeforeEach(func() {
		opts = &Options{
			SourceRepo: &scm.Repo{
				Repo:     "source_repo",
				Owner:    "source_owner",
				RepoType: "github",
			},
			DestinationRepo: &scm.Repo{
				Repo:     "dst_repo",
				Owner:    "dst_owner",
				RepoType: "github",
			},
		}
	})
	Context("NewOptions method", func() {
		var (
			rawOptions plugininstaller.RawOptions
		)
		BeforeEach(func() {
			rawOptions = plugininstaller.RawOptions{
				"sourceRepo": map[string]string{
					"owner":    "test_user",
					"repo":     "test_repo",
					"repoType": "github",
				},
				"destinationRepo": map[string]string{
					"owner":    "dst_user",
					"repo":     "dst_repo",
					"repoType": "github",
				},
			}
		})
		It("should work normal", func() {
			opts, err := NewOptions(rawOptions)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(opts.SourceRepo.Owner).Should(Equal("test_user"))
			Expect(opts.SourceRepo.Repo).Should(Equal("test_repo"))
			Expect(opts.DestinationRepo.Owner).Should(Equal("dst_user"))
			Expect(opts.DestinationRepo.Repo).Should(Equal("dst_repo"))
		})
	})

	Context("renderTplConfig method", func() {
		var (
			coverAppName string
		)
		BeforeEach(func() {
			coverAppName = "cover app"
			opts.Vars = map[string]interface{}{
				"AppName": coverAppName,
			}
		})
		It("should return with vars", func() {
			renderConfig := opts.renderTplConfig()
			Expect(len(renderConfig)).ShouldNot(BeZero())
			appName, ok := renderConfig["AppName"]
			Expect(ok).Should(BeTrue())
			Expect(appName).Should(Equal(coverAppName))
		})
	})

	Context("downloadAndRenderScmRepo method", func() {
		var (
			scmClient *scm.MockScmClient
		)
		When("download repo from scm failed", func() {
			var (
				errMsg string
			)
			BeforeEach(func() {
				errMsg = "test download repo failed"
				scmClient = &scm.MockScmClient{
					DownloadRepoError: fmt.Errorf(errMsg),
				}
			})
			It("should return error", func() {
				_, err := opts.downloadAndRenderScmRepo(scmClient)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("download repo success", func() {
			var (
				tempDir string
			)
			BeforeEach(func() {
				tempDir = GinkgoT().TempDir()
				scmClient = &scm.MockScmClient{
					DownloadRepoValue: tempDir,
				}
			})
			It("should work normal", func() {
				gitMap, err := opts.downloadAndRenderScmRepo(scmClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(len(gitMap)).Should(BeZero())
			})
		})
	})
})
