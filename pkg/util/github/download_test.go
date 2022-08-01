package github_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/argoproj/gitops-engine/pkg/utils/io"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("DownloadAsset", func() {
	Context("does DownloadAsset 200", func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/r/o/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1}]`)
			})
		})

		It("does ListReleases with wrong url ", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Org: "or",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})

	Context("does Downloaset step 1", func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/devstream-io/dtm-scaffolding-golang/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name":"t", "name":"n", "assets": [{"id":1}]}]`)
			})
		})

		It("does ListReleases with correct url ", func() {
			ghClient, err := github.NewClientWithOption(github.OptNotNeedAuth, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})

	Context("does DownloadAsset step 2", func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/oo/rr/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name":"t", "name":"n", "assets": [{"id":1}]}]`)
			})
		})

		It("does ListReleases with correct url ", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "oo",
				Repo:  "rr",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})

	Context("does DownloadAsset step 2", func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/a/b/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[{"id":1, "tag_name":"t", "name":"n", "assets":[]}`)
			})
		})

		It("does ListReleases with correct url but empty asset", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "a",
				Repo:  "b",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})

	Context("does DownloadAsset step 3", func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/ooo/rrr/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[
					{
						"id":1, 
						"tag_name":"t", 
						"name":"n", 
						"assets": [{"id":1, "name":"a"}]
					}]`)
			})
		})

		It("does get download url without browser_download_url", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "ooo",
				Repo:  "rrr",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})
	var WorkPath = "./"
	Context("does DownloadAsset step 4", func() {
		BeforeEach(func() {
			DeferCleanup(os.RemoveAll, "./"+github.DefaultWorkPath)
			mux.HandleFunc("/repos/ow/re/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `[
					{
						"id":1, 
						"tag_name":"t", 
						"name":"n", 
						"assets": [{"id":1, "name":"a", "browser_download_url":"u"}]
					}]`)
			})
		})

		It("does get download url browser_download_url with unsupported proto scheme", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "ow",
				Repo:  "re",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", "f")
			Expect(err).NotTo(Succeed())
		})

	})

	Context("does DownloadAsset step 4", func() {
		var filename = "f"
		var downloadUrl = "/download"
		BeforeEach(func() {
			mux.HandleFunc("/repos/own/rep/releases", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `[
					{
						"id":1, 
						"tag_name":"t", 
						"name":"n", 
						"assets": [{"id":1, "name":"a", "browser_download_url":"%s"}]
					}]`, serverURL+github.BaseURLPath+downloadUrl)
			})
			mux.HandleFunc(downloadUrl, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "%s", "download content")
			})
		})

		It("does get download url browser_download_url with supported proto scheme", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner:    "own",
				Repo:     "rep",
				WorkPath: WorkPath,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadAsset("t", "a", filename)
			fmt.Println(err)
			Expect(err).To(Succeed())
		})

		AfterEach(func() {
			io.DeleteFile(filepath.Join(WorkPath, filename))
		})
	})
})

var _ = Describe("DownloadLatestCodeAsZipFile", func() {
	var owner, repo, org, workPath = "owner", "repo", "org", "./"
	Context("does DownloadLatestCodeAsZipFile", func() {

		It("corrent url ", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner:    owner,
				Repo:     repo,
				Org:      org,
				WorkPath: workPath,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadLatestCodeAsZipFile()
			Expect(err).To(Succeed())
		})

		It("produce an error", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner:    owner,
				Repo:     repo,
				Org:      org,
				WorkPath: "//",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.DownloadLatestCodeAsZipFile()
			Expect(err).NotTo(Succeed())
		})

		AfterEach(func() {
			DeferCleanup(io.DeleteFile, workPath+github.DefaultLatestCodeZipfileName)
		})

	})

})
