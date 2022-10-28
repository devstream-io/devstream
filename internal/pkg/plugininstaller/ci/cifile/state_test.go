package cifile

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("getSCMFileStatus func", func() {
	var (
		mockScmClient *scm.MockScmClient
		scmPath       string
	)

	BeforeEach(func() {
		mockScmClient = &scm.MockScmClient{}
		scmPath = "test"
	})
	When("get scm pathinfo error", func() {
		BeforeEach(func() {
			mockScmClient.GetPathInfoError = errors.New("test")
		})
		It("should return empty map", func() {
			fileList := getSCMFileStatus(mockScmClient, scmPath)
			Expect(len(fileList)).Should(BeZero())
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
			fileList := getSCMFileStatus(mockScmClient, scmPath)
			Expect(len(fileList)).Should(Equal(1))
			file := fileList[0]
			v, ok := file["scmSHA"]
			Expect(ok).Should(BeTrue())
			Expect(v).Should(Equal(sha))
			v, ok = file["scmBranch"]
			Expect(ok).Should(BeTrue())
			Expect(v).Should(Equal(branch))
		})
	})
})
