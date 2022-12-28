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

var _ = Describe("CreateCommitInfo method", func() {
	var (
		repoInfo                           *git.RepoInfo
		gitlabClient                       *gitlab.Client
		testFile, repoName, reqPath, owner string
		server                             *ghttp.Server
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		testFile = "test_file"
		owner = "test_user"
		repoName = "test_repo"
		repoInfo = &git.RepoInfo{
			BaseURL: server.URL(),
			Branch:  "test",
			Repo:    repoName,
			Owner:   owner,
		}
		client, err := gitlabCommon.NewClient(
			"test", gitlabCommon.WithBaseURL(server.URL()),
			// don't retry http request when test
			gitlabCommon.WithCustomRetryMax(0),
		)
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
		server.SetAllowUnhandledRequests(true)
	})
	Context("GetPathInfo method", func() {
		When("gitlab return normal", func() {
			BeforeEach(func() {
				reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
				))
			})
			It("should work", func() {
				fileInfo, err := gitlabClient.GetPathInfo(testFile)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(fileInfo).ShouldNot(BeNil())
				Expect(fileInfo[0].Branch).Should(BeEmpty())
			})
		})
		When("gitlab return error", func() {
			BeforeEach(func() {
				reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusBadGateway, nil),
				))
			})
			It("should return err", func() {
				_, err := gitlabClient.GetPathInfo(testFile)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("DeleteFiles method", func() {
		BeforeEach(func() {
			reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
			server.RouteToHandler("POST", reqPath, ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", reqPath),
				ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
			))
		})
		It("should work normal", func() {
			_, err := gitlabClient.GetPathInfo(testFile)
			Expect(err).Error().Should(HaveOccurred())
		})

	})
	AfterEach(func() {
		server.Close()
	})
})
