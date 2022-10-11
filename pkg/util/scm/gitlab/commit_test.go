package gitlab_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	gitlabCommon "github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

var _ = Describe("commit method", func() {
	var (
		repoInfo                                                             *git.RepoInfo
		commitInfo                                                           *git.CommitInfo
		gitlabClient                                                         *gitlab.Client
		server                                                               *ghttp.Server
		repoName, reqPath, branch, owner, commitMsg, gitFile, gitFileContent string
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		server.SetAllowUnhandledRequests(true)
		owner = "test_user"
		repoName = "test_repo"
		branch = "test_branch"
		gitFile = "test_git.file"
		gitFileContent = "test_git_content"
		commitMsg = "test msg"
		reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/commits", apiRootPath, owner, repoName)
		repoInfo = &git.RepoInfo{
			BaseURL: server.URL(),
			Branch:  branch,
			Repo:    repoName,
			Owner:   owner,
		}
		commitInfo = &git.CommitInfo{
			CommitMsg:    commitMsg,
			CommitBranch: branch,
			GitFileMap: git.GitFileContentMap{
				gitFile: []byte(gitFileContent),
			},
		}
		client, err := gitlabCommon.NewClient(
			"test", gitlabCommon.WithBaseURL(server.URL()))
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
	})
	// Context("CreateCommitInfo method", func() {
	// It("should return gitlab commit options", func() {
	// commitInfoData := gitlabClient.CreateCommitInfo(gitlabCommon.FileCreate, commitInfo)
	// Expect(*commitInfoData.Branch).Should(Equal(branch))
	// Expect(*commitInfoData.CommitMessage).Should(Equal(commitMsg))
	// actions := commitInfoData.Actions
	// Expect(len(actions)).Should(Equal(1))
	// action := actions[0]
	// Expect(*action.Action).Should(Equal(gitlabCommon.FileCreate))
	// Expect(*action.FilePath).Should(Equal(gitFile))
	// Expect(*action.Content).Should(Equal(gitFileContent))
	// })
	// })

	Context("CommitActions method", func() {
		BeforeEach(func() {
			server.RouteToHandler("POST", reqPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", reqPath),
				ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
			))
		})
		It("should work normal", func() {
			needRollBack, err := gitlabClient.PushFiles(commitInfo, false)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(needRollBack).Should(BeFalse())
			err = gitlabClient.DeleteFiles(commitInfo)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	AfterEach(func() {
		server.Close()
	})
})
