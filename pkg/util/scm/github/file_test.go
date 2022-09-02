package github_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

var _ = Describe("Github files methods", func() {
	var (
		s                                *ghttp.Server
		repoName, owner, branch, reqPath string
		c                                *github.Client
	)
	BeforeEach(func() {
		repoName = "test_repo"
		owner = "test_owner"
		branch = "test_branch"
		s = ghttp.NewServer()
		repoInfo := &git.RepoInfo{
			Repo:   repoName,
			Branch: branch,
			Owner:  owner,
		}
		c = newTestClient(s.URL(), repoInfo)
	})
	AfterEach(func() {
		s.Close()
	})
	Context("GetLocationInfo func", func() {
		var (
			testFile, filePath, fileSHA string
		)
		When("file exist", func() {
			BeforeEach(func() {
				testFile = "testFile"
				filePath = "/test_path"
				fileSHA = "s"
				reqPath = fmt.Sprintf("%s/repos/test_owner/test_repo/contents/%s", basePath, testFile)
				query := fmt.Sprintf("ref=%s", branch)
				log.Infof("%s<->%s", reqPath, query)
				mockRsp := map[string]string{
					"type": "test",
					"path": filePath,
					"sha":  fileSHA,
				}
				s.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath, query),
					ghttp.RespondWithJSONEncoded(http.StatusOK, mockRsp),
				))
			})
			It("should work", func() {
				info, err := c.GetLocationInfo(testFile)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(info).ShouldNot(BeNil())
				Expect(len(info)).Should(Equal(1))
				Expect(info[0].Branch).Should(Equal(branch))
				Expect(info[0].SHA).Should(Equal(fileSHA))
				Expect(info[0].Path).Should(Equal(filePath))
			})
		})
		When("file not exist", func() {
			BeforeEach(func() {
				testFile = "not_exist"
				fileSHA = "s"
				reqPath = fmt.Sprintf("%s/repos/test_owner/test_repo/contents/%s", basePath, testFile)
				query := fmt.Sprintf("ref=%s", branch)
				log.Infof("%s<->%s", reqPath, query)
				s.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath, query),
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, nil),
				))
			})
			It("should return empty list", func() {
				info, err := c.GetLocationInfo(testFile)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(info)).Should(Equal(0))
			})
		})
	})
})
