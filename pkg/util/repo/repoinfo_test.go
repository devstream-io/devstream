package repo_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/repo"
)

var _ = Describe("RepoInfo struct", func() {
	var (
		repoName, branch, owner, org string
		repoInfo                     *repo.RepoInfo
	)
	BeforeEach(func() {
		repoName = "test_repo"
		branch = "test_branch"
		owner = "test_owner"
		org = "test_org"
		repoInfo = &repo.RepoInfo{
			Repo:   repoName,
			Branch: branch,
			Owner:  owner,
			Org:    org,
		}
	})
	Context("GetRepoNameWithBranch method", func() {
		It("should return repo-with-branch", func() {
			result := repoInfo.GetRepoNameWithBranch()
			Expect(result).Should(Equal(fmt.Sprintf("%s-%s", repoName, branch)))
		})
	})
	Context("GetRepoOwner method", func() {
		It("should return owner", func() {
			result := repoInfo.GetRepoOwner()
			Expect(result).Should(Equal(repoInfo.Org))
		})
	})
	Context("GetRepoPath method", func() {
		It("should return repo path", func() {
			result := repoInfo.GetRepoPath()
			Expect(result).Should(Equal(fmt.Sprintf("%s/%s", repoInfo.GetRepoOwner(), repoName)))
		})
	})
})
