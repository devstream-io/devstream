package github_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/merico-dev/stream/internal/pkg/util/github"
)

var _ = Describe("GitHub", func() {
	var testRelease = "v0.0.1"
	var testAsset = "argocdapp_0.0.1.so"
	var workPath = ".github-repo-scaffolding-golang"
	Context("Client without auth enabled", func() {
		var ghClient *github.Client
		var err error
		BeforeEach(func() {
			ghClient, err = github.NewClient(&github.Option{
				Owner:    "merico-dev",
				Repo:     "stream",
				NeedAuth: false,
				WorkPath: workPath,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})

		It("Should get assets", func() {
			err := ghClient.DownloadAsset(testRelease, testAsset)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err = os.RemoveAll(filepath.Join(workPath, testAsset))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
