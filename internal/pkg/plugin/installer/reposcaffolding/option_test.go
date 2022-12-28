package reposcaffolding

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("Options struct", func() {
	var (
		opts *Options
	)
	BeforeEach(func() {
		opts = &Options{
			SourceRepo: &git.RepoInfo{
				Repo:     "source_repo",
				Owner:    "source_owner",
				RepoType: "github",
			},
			DestinationRepo: &git.RepoInfo{
				Repo:     "dst_repo",
				Owner:    "dst_owner",
				RepoType: "github",
			},
		}
	})
	Context("NewOptions method", func() {
		var (
			rawOptions configmanager.RawOptions
		)
		BeforeEach(func() {
			rawOptions = configmanager.RawOptions{
				"sourceRepo": map[string]string{
					"owner":   "test_user",
					"name":    "test_repo",
					"scmType": "github",
				},
				"destinationRepo": map[string]string{
					"owner":   "dst_user",
					"name":    "dst_repo",
					"scmType": "github",
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
			repo, ok := renderConfig["Repo"]
			Expect(ok).Should(BeTrue())
			Expect(repo).Should(Equal(map[string]string{"Owner": "dst_owner", "Name": "dst_repo"}))
		})
	})
})
