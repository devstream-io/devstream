package github_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("Pullrequest", func() {
	Context(("do NewPullRequest 200"), func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/devstream-io/dtm-scaffolding-golang/pulls", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"number":1}`)
			})
		})

		It("do create a new pr", func() {
			ghClient, err := github.NewClientWithOption(github.OptNotNeedAuth, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			n, err := ghClient.NewPullRequest("from", "to")
			Expect(n).To(Equal(1))
			Expect(err).To(Succeed())
		})

	})

	Context(("do NewPullRequest 404"), func() {

		It("do create a new pr with wrong url", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "",
				Org:   "or",
				Repo:  "r",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			n, err := ghClient.NewPullRequest("from", "to")
			Expect(n).To(Equal(0))
			Expect(err).NotTo(Equal(nil))
		})
	})

	Context(("do MergePullRequest"), func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/devstream-io/dtm-scaffolding-golang/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `
					{
					"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
					"merged": true,
					"message": "Pull Request successfully merged"
					}`)
			})
		})

		It("200", func() {
			ghClient, err := github.NewClientWithOption(github.OptNotNeedAuth, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.MergePullRequest(1, github.MergeMethodRebase)
			Expect(err).To(Succeed())
		})
	})

	Context(("do MergePullRequest"), func() {

		It("404", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "",
				Org:   "or",
				Repo:  "r",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.MergePullRequest(1, github.MergeMethodRebase)
			Expect(err).To(HaveOccurred())
		})
	})

	Context(("do MergePullRequest"), func() {

		BeforeEach(func() {
			mux.HandleFunc("/repos/or/r/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `
					{
					"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
					"merged": false,
					"message": "Pull Request successfully merged"
					}`)
			})
		})

		It("return merged false", func() {
			ghClient, err := github.NewClientWithOption(&github.Option{
				Owner: "",
				Org:   "or",
				Repo:  "r",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
			err = ghClient.MergePullRequest(1, github.MergeMethodRebase)
			Expect(err).NotTo(Succeed())
		})
	})
})
