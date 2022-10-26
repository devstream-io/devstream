package github_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	githubCommon "github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	githubUtil "github.com/devstream-io/devstream/pkg/util/scm/github"
)

var (
	mux       *http.ServeMux
	serverURL string
	teardown  func()
)

const basePath = "/api-v3"

var _ = BeforeSuite(func() {
	mux, serverURL, teardown = githubUtil.Setup()
})

var _ = AfterSuite(func() {
	teardown()
})

func CreateClientWithOr(opt *git.RepoInfo) *githubUtil.Client {
	c, err := githubUtil.NewClientWithOption(opt, serverURL)
	Expect(c).NotTo(Equal(nil))
	Expect(err).To(Succeed())
	return c
}

func TestPlanmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitHub Suite")
}

func newTestClient(baseUrl string, repoInfo *git.RepoInfo) *githubUtil.Client {
	githubClient := githubCommon.NewClient(nil)
	url, _ := url.Parse(baseUrl + basePath + "/")

	githubClient.BaseURL = url
	githubClient.UploadURL = url

	return &githubUtil.Client{
		RepoInfo: repoInfo,
		Client:   githubClient,
		Context:  context.Background(),
	}
}
