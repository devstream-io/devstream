package reposcaffolding

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("Options struct", func() {
	var (
		opts *options
	)
	BeforeEach(func() {
		opts = &options{
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
