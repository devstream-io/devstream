package github_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("NewPullRequest", func() {
	const (
		fromBranch, toBranch = "fb", "tb"
		owner, repo          = "owner", "repo"
		rightOrg, wrongOrg   = "org", "/"
	)

	var (
		s    *ghttp.Server
		org  string
		opts *github.Option
	)

	JustBeforeEach(func() {
		opts = &github.Option{
			Owner: owner,
			Repo:  repo,
			Org:   org,
		}
	})

	AfterEach(func() {
		s.Close()
	})

	When("Create", func() {
		BeforeEach(func() {
			s = ghttp.NewServer()
			org = wrongOrg
		})
		It("url is incorrect", func() {
			s.SetAllowUnhandledRequests(true)
			c, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(Equal(nil))
			n, err := c.NewPullRequest(fromBranch, toBranch)
			Expect(err).To(HaveOccurred())
			fmt.Println(err)
			Expect(n).To(Equal(0))
		})
	})
	When("Create", func() {
		BeforeEach(func() {
			s = ghttp.NewServer()
			org = rightOrg
		})
		It("url is correct", func() {
			u := github.BaseURLPath + fmt.Sprintf("/repos/%v/%v/pulls", org, repo)
			s.RouteToHandler("POST", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"number": 1}`)
			})
			c, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(Equal(nil))
			n, err := c.NewPullRequest(fromBranch, toBranch)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(1))
		})
	})
})

var _ = Describe("MergePullRequest", func() {
	const (
		fromBranch, toBranch = "fb", "tb"
		owner, repo          = "owner", "repo"
		rightOrg, wrongOrg   = "org", "/"
		number               = 1
	)

	var (
		s    *ghttp.Server
		org  string
		opts *github.Option
	)

	JustBeforeEach(func() {
		opts = &github.Option{
			Owner: owner,
			Repo:  repo,
			Org:   org,
		}
	})

	AfterEach(func() {
		s.Close()
	})

	When("Merge", func() {
		BeforeEach(func() {
			s = ghttp.NewServer()
			org = wrongOrg
		})
		It("url is incorrect", func() {
			s.SetAllowUnhandledRequests(true)
			c, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(Equal(nil))
			err = c.MergePullRequest(number, github.MergeMethodMerge)
			Expect(err).To(HaveOccurred())
		})
	})
	When("Merge", func() {
		BeforeEach(func() {
			s = ghttp.NewServer()
			org = rightOrg
		})
		It("url is correct but merged is false", func() {
			u := github.BaseURLPath + fmt.Sprintf("/repos/%v/%v/pulls/%d/merge", org, repo, number)
			s.RouteToHandler("PUT", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			c, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(Equal(nil))
			err = c.MergePullRequest(number, github.MergeMethodMerge)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("merge failed"))
		})
	})
	When("Merge", func() {
		BeforeEach(func() {
			s = ghttp.NewServer()
			org = rightOrg
		})
		It("url is correct", func() {
			u := github.BaseURLPath + fmt.Sprintf("/repos/%v/%v/pulls/%d/merge", org, repo, number)
			s.RouteToHandler("PUT", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"merged": true}`)
			})
			c, err := github.NewClientWithOption(opts, s.URL())
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(Equal(nil))
			err = c.MergePullRequest(number, github.MergeMethodMerge)
			Expect(err).To(Succeed())
		})
	})
})
