package github_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/argoproj/gitops-engine/pkg/utils/io"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

var _ = Describe("DownloadAsset", func() {
	const (
		owner, repoName              = "owner", "repo"
		rightOrg, wrongOrg           = "org", "/"
		tagName, assetName, fileName = "t", "a", "f3"
	)

	var (
		s        *ghttp.Server
		org      string
		opts     *git.RepoInfo
		workPath string
	)

	JustBeforeEach(func() {
		opts = &git.RepoInfo{
			Owner: owner,
			Repo:  repoName,
			Org:   org,
		}
		if len(workPath) != 0 {
			opts.WorkPath = workPath
		}
	})
	When("// 1. get releases: the ListReleases url is incorrect", func() {
		BeforeEach(func() {
			org = wrongOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			s.SetAllowUnhandledRequests(true)
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(HaveOccurred())
		})
	})

	When("// 1. get releases: the ListReleases url is correct but assets is empty", func() {
		BeforeEach(func() {
			org = rightOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `[{"id":1, "tag_name": "t"}]`)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("500"))
		})
	})

	When("// 2. get assets: tagName is equal and the assets is empty", func() {
		BeforeEach(func() {
			org = rightOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name": "t", "name": "n", "assets": []}]`)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("assets is empty"))
		})
	})

	When("// 2. get assets: tagName is not equal and the assets is empty", func() {
		BeforeEach(func() {
			org = rightOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name": "tttt", "name": "n", "assets": []}]`)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("release with tag <%s> was not found", tagName)))
		})
	})

	When("// 3. get download url: browser_download_url is empty", func() {
		BeforeEach(func() {
			org = rightOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name": "t", "name": "a", "assets": [{"browser_download_url": ""}]}]`)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("failed to got the download url for %s, maybe it not exists", assetName)))
		})
	})

	When("// 4. download: filename is '.' ", func() {
		BeforeEach(func() {
			org = rightOrg
			s = ghttp.NewServer()
		})

		It("should return error", func() {
			downloadUrl := s.URL() + github.BaseURLPath + "/download"
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `[
							{
								"id":1,
								"tag_name":"t",
								"name":"n",
								"assets": [{"id":1, "name":"a", "browser_download_url":"%s"}]
							}]`, downloadUrl)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, ".")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("filename must not be dir"))
		})
	})

	When("// 4. download", func() {
		BeforeEach(func() {
			org = rightOrg
			workPath = os.TempDir()
			s = ghttp.NewServer()
		})

		It("should return no error", func() {
			downloadUrl := s.URL() + github.BaseURLPath + "/download"
			u := fmt.Sprintf("/repos/%s/%s/releases", org, repoName)
			s.RouteToHandler("GET", github.BaseURLPath+"/download", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "file content")
			})
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `[
							{
								"id":1,
								"tag_name":"t",
								"name":"n",
								"assets": [{"id":1, "name":"a", "browser_download_url":"%s"}]
							}]`, downloadUrl)
			})
			ghClient, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset(tagName, assetName, fileName)
			Expect(err).To(Succeed())
		})
	})

	AfterEach(func() {
		s.Close()
		DeferCleanup(io.DeleteFile, workPath+"/"+fileName)
	})
})
