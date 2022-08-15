package github_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("Secrets", func() {
	var owner, repoName, org = "o", "r", "or"
	var registerUrl string = fmt.Sprintf("/repos/%v/%v/actions/secrets/public-key", org, repoName)
	sk, sv := "sk", "sv"

	Context("does AddRepoSecret", func() {
		It("step1: does GetRepoPublicKey with wrong url", func() {
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))
			err = rightClient.AddRepoSecret(sk, sv)
			Expect(err).NotTo(Succeed())
		})

		It("step2: does AddRepoSecret with correct url", func() {
			registerUrl = fmt.Sprintf("/repos/%v/%v/actions/secrets/public-key", org, repoName)
			mux.HandleFunc(registerUrl, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
			})
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))
			err = rightClient.AddRepoSecret(sk, sv)
			Expect(err).NotTo(Succeed())
		})

		It("step3: does CreateOrUpdateRepoSecret with wrong url", func() {
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))

			err = rightClient.AddRepoSecret(sk, sv)
			Expect(err).NotTo(Succeed())
		})

		It("step3: does CreateOrUpdateRepoSecret with correct url", func() {
			registerUrl = fmt.Sprintf("/repos/%v/%v/actions/secrets/%v", org, repoName, sk)
			mux.HandleFunc(registerUrl, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
			})
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))
			err = rightClient.AddRepoSecret(sk, sv)

			Expect(err).To(Succeed())
		})

	})

	Context("RepoSecretExists", func() {
		It("does RepoSecretExists with wrong url", func() {
			wrongClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  "rrrr",
				Org:   "ororor",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(wrongClient).NotTo(Equal(nil))
			b, err := wrongClient.RepoSecretExists(sk)
			Expect(err).To(Succeed())
			Expect(b).To(Equal(false))
		})

		It("does RepoSecretExists with correct url", func() {
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))
			b, err := rightClient.RepoSecretExists(sk)
			Expect(err).To(Succeed())
			Expect(b).To(Equal(true))
		})
	})

	Context("DeleteRepoSecret", func() {
		It("does DeleteRepoSecret with wrong url", func() {
			wrongClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  "rrrr",
				Org:   "ororor",
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(wrongClient).NotTo(Equal(nil))
			err = wrongClient.DeleteRepoSecret(sk)
			Expect(err).To(Succeed())
		})

		It("does DeleteRepoSecret with correct url", func() {
			rightClient, err := github.NewClientWithOption(&git.RepoInfo{
				Owner: owner,
				Repo:  repoName,
				Org:   org,
			}, serverURL)
			Expect(err).NotTo(HaveOccurred())
			Expect(rightClient).NotTo(Equal(nil))
			err = rightClient.DeleteRepoSecret(sk)
			Expect(err).To(Succeed())
		})
	})
})
