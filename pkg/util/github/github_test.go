package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	// "github.com/devstream-io/devstream/pkg/util/github"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-v3"
)

var (
	optNotNeedAuth = &Option{
		Owner: "",
		Org:   "devstream-io",
		Repo:  "dtm-scaffolding-golang",
	}
	optNeedAuth = &Option{
		Owner:    "",
		Org:      "devstream-io",
		Repo:     "dtm-scaffolding-golang",
		NeedAuth: true,
	}
)

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) (mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	return mux, server.URL, server.Close
}

func getClientWithOption(t *testing.T, opt *Option, severUrl string) *Client {
	client, err := NewClient(opt)
	if err != nil {
		t.Error(err)
	}

	url, _ := url.Parse(severUrl + baseURLPath + "/")

	client.Client.BaseURL = url
	client.Client.UploadURL = url
	return client
}

var _ = Describe("GitHub", func() {
	Context(("Client with cacahe"), func() {
		var ghClient *Client
		var err error

		BeforeEach(func() {
			ghClient, err = NewClient(optNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})

		It("with cacahe client", func() {
			ghClient, err = NewClient(optNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
	})

	Context("Client without auth enabled", func() {
		var ghClient *Client
		var err error
		It("", func() {
			ghClient, err = NewClient(optNotNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
	})

	Context("Client with auth enabled but not github token", func() {
		var ghClient *Client
		var err error
		It("", func() {
			ghClient, err = NewClient(optNeedAuth)
			fmt.Printf("--->: %+v\n", ghClient)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Client with auth enabled and github token", func() {
		var ghClient *Client
		var err error
		BeforeEach(func() {
			os.Setenv("GITHUB_TOKEN", "GITHUB_TOKEN")
		})
		It("", func() {
			ghClient, err = NewClient(optNeedAuth)
			Expect(err).NotTo(HaveOccurred())
			Expect(ghClient).NotTo(Equal(nil))
		})
		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})
	})
})
