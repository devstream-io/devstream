package scm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"errors"

	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GetGitFileStats func", func() {
	var (
		mockScmClient *scm.MockScmClient
		gitFileMap    git.GitFileContentMap
	)

	BeforeEach(func() {
		mockScmClient = &scm.MockScmClient{}
		gitFileMap = git.GitFileContentMap{
			"testFile": []byte("test_Content"),
		}
	})
	When("get scm pathinfo error", func() {
		BeforeEach(func() {
			mockScmClient.GetPathInfoError = errors.New("test")
		})
		It("should return empty map", func() {
			_, err := scm.GetGitFileStats(mockScmClient, gitFileMap)
			Expect(err).Should(HaveOccurred())
		})
	})
	When("get scm pathinfo return", func() {
		var (
			path, sha, branch string
		)
		BeforeEach(func() {
			path = "test_path"
			sha = "test_sha"
			branch = "test_branch"
			gitFileStatus := []*git.RepoFileStatus{
				{
					Path:   path,
					SHA:    sha,
					Branch: branch,
				},
			}
			mockScmClient.GetPathInfoReturnValue = gitFileStatus
		})
		It("should return fileMap", func() {
			status, err := scm.GetGitFileStats(mockScmClient, gitFileMap)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(status)).Should(Equal(1))
			Expect(status).Should(Equal(map[string]any{
				"testFile": map[string]any{
					"localSHA": "bbed2c0c2935a0e860cd4f6212aba4d6",
					"scm": []map[string]string{
						{
							"scmSHA":    "test_sha",
							"scmBranch": "test_branch",
						},
					},
				},
			}))
		})
	})
})
