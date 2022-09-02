package gitlab_test

import (
	"fmt"
	"net/http"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	gitlabCommon "github.com/xanzy/go-gitlab"
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
		client, err := gitlabCommon.NewClient("test", gitlabCommon.WithBaseURL(server.URL()))
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
		server.SetAllowUnhandledRequests(true)
	})
	Context("FileExists method", func() {
		BeforeEach(func() {
			reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
		})
		When("get files url return err", func() {
			BeforeEach(func() {
				server.SetUnhandledRequestStatusCode(http.StatusNotFound)
			})
			It("should return false", func() {
				exist, err := gitlabClient.FileExists(testFile)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(exist).Should(BeFalse())
			})
		})
		When("git files return normal", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
				))
			})
			It("should return true", func() {
				exist, err := gitlabClient.FileExists(testFile)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(exist).Should(BeTrue())
			})
		})
	})
	Context("GetLocationInfo method", func() {
		When("gitlab return normal", func() {
			BeforeEach(func() {
				reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
				server.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil),
				))
			})
			It("should work", func() {
				fileInfo, err := gitlabClient.GetLocationInfo(testFile)
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
				_, err := gitlabClient.GetLocationInfo(testFile)
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
			_, err := gitlabClient.GetLocationInfo(testFile)
			Expect(err).Error().Should(HaveOccurred())
		})

	})
	AfterEach(func() {
		server.Close()
	})
})
