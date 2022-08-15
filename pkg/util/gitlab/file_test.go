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
		reqPath = fmt.Sprintf("%sprojects/%s/%s/repository/files/%s", apiRootPath, owner, repoName, testFile)
		client, err := gitlabCommon.NewClient("test", gitlabCommon.WithBaseURL(server.URL()))
		Expect(err).Error().ShouldNot(HaveOccurred())
		gitlabClient = &gitlab.Client{
			Client:   client,
			RepoInfo: repoInfo,
		}
	})
	When("get files url return err", func() {
		BeforeEach(func() {
			server.SetAllowUnhandledRequests(true)
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
			server.SetAllowUnhandledRequests(true)
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
	AfterEach(func() {
		server.Close()
	})
})
