package github

import (
	"fmt"
	"net/http"
	"testing"

	"gotest.tools/assert/cmp"
)

func TestClient_GetLatestReleaseTagName(t *testing.T) {
	mux, serverUrl, teardown := setup(t)
	defer teardown()

	tests := []struct {
		name        string
		client      *Client
		registerUrl string
		wantMethod  string
		wantReqBody bool
		reqBody     string
		respBody    string
		wantTag     string
		wantErr     bool
	}{
		{"base err != nil", getClientWithOption(
			t, &Option{Owner: ""}, serverUrl,
		), "/repos2/o/r/releases/latest", http.MethodGet, false, "", "", "", true},
		{"base 200", getClientWithOption(
			t, &Option{Owner: "", Org: "o", Repo: "r"}, serverUrl,
		), "/repos/o/r/releases/latest", http.MethodGet, false, "", `{"id":3,"tag_name":"v1.0.0"}`, "v1.0.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.wantMethod)
				testBody(t, r, "")
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
