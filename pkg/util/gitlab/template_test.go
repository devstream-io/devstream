package gitlab_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	gitlabCommon "github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
)

var _ = Describe("repo methods", func() {
	var (
		repoInfo                         *git.RepoInfo
		gitlabClient                     *gitlab.Client
		server                           *ghttp.Server
		repoName, branch, reqPath, owner string
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		server.SetAllowUnhandledRequests(true)
		owner = "test_user"
		repoName = "test_repo"
		branch = "test_branch"
		repoInfo = &git.RepoInfo{
			BaseURL: server.URL(),
			Branch:  branch,
			Repo:    repoName,
			Owner:   owner,
		}
		client, err := gitlabCommon.NewClient(
			"test", gitlabCommon.WithBaseURL(server.URL()))
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
	})
	Context("GetGitLabCIGolangTemplate method", func() {
		When("ci file exist", func() {
			BeforeEach(func() {
				reqPath = fmt.Sprintf("%stemplates/gitlab_ci_ymls/Go", apiRootPath)
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
				))
			})
			It("should return content if file exist", func() {
				_, err := gitlabClient.GetGitLabCIGolangTemplate()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("ci file not exist", func() {
			BeforeEach(func() {
				reqPath = fmt.Sprintf("%stemplates/gitlab_ci_ymls/Go", apiRootPath)
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, nil),
				))
			})
			It("should return content if file exist", func() {
				_, err := gitlabClient.GetGitLabCIGolangTemplate()
				Expect(err).Should(HaveOccurred())
			})

		})
	})
	AfterEach(func() {
		server.Close()
	})
})
