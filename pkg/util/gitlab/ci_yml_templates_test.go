package gitlab

import (
	"os"
	"strings"
	"testing"
)

var gitlab_token = "FAKE_GITLAB_TOKEN"

func TestClient_GetGitLabCIGolangTemplate(t *testing.T) {
	baseURL := ""
	os.Setenv("GITLAB_TOKEN", gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name        string
		client      *Client
		wantContent bool
		wantErr     bool
	}{
		// TODO: Add test cases.
		{"base", client, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetGitLabCIGolangTemplate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetGitLabCIGolangTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, "golang") {
				t.Errorf("Client.GetGitLabCIGolangTemplate's content doesn't include 'golang'")
			}
		})
	}
}
