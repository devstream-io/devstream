package github

import (
	"fmt"
	"net/http"
	"testing"

	"gotest.tools/assert/cmp"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type releaseTest struct {
	BaseTest
	wantTag string
	wantErr bool
}

func TestClient_GetLatestReleaseTagName(t *testing.T) {
	mux, serverUrl, teardown := Setup()
	defer teardown()

	tests := []releaseTest{
		{
			BaseTest{"base err != nil", GetClientWithOption(
				t, &git.RepoInfo{Owner: ""}, serverUrl,
			),
				"/repos2/o/r/releases/latest", http.MethodGet, false, "", ""},
			"", true},
		{
			BaseTest{"base 200", GetClientWithOption(
				t, &git.RepoInfo{Owner: "", Org: "o", Repo: "r"}, serverUrl,
			),
				"/repos/o/r/releases/latest", http.MethodGet, false, "", `{"id":3,"tag_name":"v1.0.0"}`},
			"v1.0.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				DoTestMethod(t, r, tt.wantMethod)
				DoTestBody(t, r, "")
				fmt.Fprint(w, `{"id":3,"tag_name":"v1.0.0"}`)
			})
			tag, err := tt.client.GetLatestReleaseTagName()
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetLatestReleaseTagName returned error: %v", err)
			}
			if !cmp.Equal(tag, tt.wantTag)().Success() {
				t.Errorf("Repositories.GenerateReleaseNotes returned %+v, want %+v", tag, tt.wantTag)
			}
		})
	}
}
