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

var _ = Describe("repo methods", func() {
	var (
		repoInfo                                     *git.RepoInfo
		gitlabClient                                 *gitlab.Client
		server                                       *ghttp.Server
		repoName, branch, reqPath, owner, visibility string
		expectReqBody                                []byte
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		server.SetAllowUnhandledRequests(true)
		owner = "test_user"
		repoName = "test_repo"
		branch = "test_branch"
		visibility = "internal"
		repoInfo = &git.RepoInfo{
			BaseURL:    server.URL(),
			Branch:     branch,
			Repo:       repoName,
			Owner:      owner,
			Visibility: visibility,
		}
		client, err := gitlabCommon.NewClient(
			"test", gitlabCommon.WithBaseURL(server.URL()))
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
	})
	Context("InitRepo method", func() {
		BeforeEach(func() {
			reqPath = fmt.Sprintf("%sprojects", apiRootPath)
			expectReqBody = []byte(fmt.Sprintf(`{"auto_devops_enabled":false,"default_branch":"%s","description":"Bootstrapped by DevStream.","name":"%s","visibility":"%s","merge_requests_enabled":true,"snippets_enabled":true}`, branch, repoName, visibility))
			server.RouteToHandler("POST", reqPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", reqPath),
				ghttp.VerifyBody(expectReqBody),
				ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
			))
		})
		It("should create repo", func() {
			err := gitlabClient.InitRepo()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Context("DeleteRepo method", func() {
		BeforeEach(func() {
			reqPath = fmt.Sprintf("%sprojects/%s/%s", apiRootPath, owner, repoName)
			server.RouteToHandler("DELETE", reqPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", reqPath),
				ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
			))
		})
		It("should delete repo", func() {
			err := gitlabClient.DeleteRepo()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Context("DescribeRepo method", func() {
		BeforeEach(func() {
			reqPath = fmt.Sprintf("%sprojects/%s/%s", apiRootPath, owner, repoName)
			server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", reqPath),
				ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
			))
		})
		It("should return repo info", func() {
			_, err := gitlabClient.DescribeRepo()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
