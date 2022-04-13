package github_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("GitHub", func() {
	// var testTag = "v0.0.1"
	// var testAsset = "dtm-scaffolding-golang-v0.0.1.tar.gz"
	var workPath = ".github-repo-scaffolding-golang"
	Context("Client without auth enabled", func() {
		var ghClient *github.Client
		var err error
		BeforeEach(func() {
			ghClient, err = github.NewClient(&github.Option{
				Owner:    "",
				Org:      "devstream-io",
				Repo:     "dtm-scaffolding-golang",
				NeedAuth: false,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})

		It("Should get assets", func() {
			// The test case below will trigger "API rate limit exceeded" error.
			// But it's useful to test locally in the future, so just comment it.

			// err := ghClient.DownloadAsset(testTag, testAsset)
			// Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err = os.RemoveAll(workPath)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Client with auth enabled", func() {

	})
})
