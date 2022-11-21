package git_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("RepoFileStatus struct", func() {
	var (
		repoFileStatus    *git.RepoFileStatus
		path, branch, sha string
	)
	BeforeEach(func() {
		path = "testPath"
		branch = "test_branch"
		sha = "s"
		repoFileStatus = &git.RepoFileStatus{
			Path:   path,
			Branch: branch,
			SHA:    sha,
		}
	})
	Context("EncodeToGitHubContentOption method", func() {
		It("should return github contentInfo", func() {
			commitMsg := "commit"
			result := repoFileStatus.EncodeToGitHubContentOption(commitMsg)
			Expect(*result.SHA).Should(Equal(sha))
			Expect(*result.Message).Should(Equal(commitMsg))
			Expect(*result.Branch).Should(Equal(branch))
		})
	})
})

var _ = Describe("CalculateGitHubBlobSHA func", func() {
	var content string
	It("should return as expect", func() {
		content = "test Content"
		Expect(git.CalculateGitHubBlobSHA([]byte(content))).Should(Equal("d9c012c6ecfcc8ce04a6538cc43490b1d5401241"))
	})
})

var _ = Describe("CalculateLocalFileSHA func", func() {
	var content string
	It("should return as expect", func() {
		content = "test Content"
		Expect(git.CalculateLocalFileSHA([]byte(content))).Should(Equal("f73d59b513c429a33da4f7efe70c7af3"))
	})
})
