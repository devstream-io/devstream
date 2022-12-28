package github_test

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		commitInfo                       *git.CommitInfo
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
		commitInfo = &git.CommitInfo{
			CommitMsg:    "test",
			CommitBranch: branch,
			GitFileMap: map[string][]byte{
				"srcPath": []byte("test data"),
			},
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
				info, err := c.GetPathInfo(testFile)
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
				info, err := c.GetPathInfo(testFile)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(info)).Should(Equal(0))
			})
		})
	})

	Context("PushLocalPathToRepo", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})

		It("1. create new branch from main", func() {
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			r, err := c.PushFiles(commitInfo, false)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			Expect(r).To(Equal(false))
		})
		It("2. create new branch from main", func() {
			// u := fmt.Sprintf("/repos/%v/%v/git/ref/heads/%s", org, repo, filePath)
			u := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repoName, strings.Trim(os.TempDir(), "/"))
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			s.RouteToHandler("GET", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "")
			})
			r, err := c.PushFiles(commitInfo, false)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			Expect(r).To(Equal(false))
		})
	})

})
