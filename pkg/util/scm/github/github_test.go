package github

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

const (
	// BaseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	BaseURLPath = "/api-v3"
)

var (
	OptNotNeedAuth = &git.RepoInfo{
		Owner: "",
		Org:   "devstream-io",
		Repo:  "dtm-repo-scaffolding-golang",
	}
	OptNeedAuth = &git.RepoInfo{
		Owner:    "",
		Org:      "devstream-io",
		Repo:     "dtm-repo-scaffolding-golang",
		NeedAuth: true,
	}
)

type BaseTest struct {
	name        string //test name
	client      *Client
	registerUrl string // url resigtered in mock server
	wantMethod  string // wanted method in mock server
	wantReqBody bool   // want request body nor not in mock server
	reqBody     string // content of request body
	respBody    string // content of response body
}

func DoTestMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func DoTestBody(t *testing.T, r *http.Request, want string) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// Setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func Setup() (mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(BaseURLPath+"/", http.StripPrefix(BaseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	return mux, server.URL, server.Close
}

func GetClientWithOption(t *testing.T, opt *git.RepoInfo, severUrl string) *Client {
	client, err := NewClient(opt)
	if err != nil {
		t.Error(err)
	}

	url, _ := url.Parse(severUrl + BaseURLPath + "/")

	client.Client.BaseURL = url
	client.Client.UploadURL = url
	return client
}

func NewClientWithOption(opt *git.RepoInfo, severUrl string) (*Client, error) {
	client, err := NewClient(opt)
	if err != nil {
		return nil, err
	}

	url, _ := url.Parse(severUrl + BaseURLPath + "/")

	client.Client.BaseURL = url
	client.Client.UploadURL = url
	return client, nil
}

var _ = Describe("GitHub", func() {
	Context("Client with cacahe", func() {
		var ghClient *Client
		var err error

		BeforeEach(func() {
			ghClient, err = NewClient(OptNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})

		It("with cacahe client", func() {
			ghClient, err = NewClient(OptNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
	})

	Context("Client without auth enabled", func() {
		var ghClient *Client
		var err error
		It("", func() {
			ghClient, err = NewClient(OptNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
	})

	Context("Client with auth enabled but not github token", func() {
		var existToken string
		BeforeEach(func() {
			existToken = os.Getenv("GITHUB_TOKEN")
			err := os.Unsetenv("GITHUB_TOKEN")
			Expect(err).NotTo(HaveOccurred())
		})
		It("", func() {
			_, err := NewClient(OptNeedAuth)
			Expect(err).To(HaveOccurred())
		})
		AfterEach(func() {
			if existToken != "" {
				err := os.Setenv("GITHUB_TOKEN", existToken)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Context("Client with auth enabled and github token", func() {
		var ghClient *Client
		var err error
		BeforeEach(func() {
			os.Setenv("GITHUB_TOKEN", "GITHUB_TOKEN")
		})
		It("", func() {
			ghClient, err = NewClient(OptNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})
	})
})
