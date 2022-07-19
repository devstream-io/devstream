package gitlab

import (
	"os"
	"testing"
)

func TestClient_FileExists(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name     string
		client   *Client
		project  string
		branch   string
		filename string
		getFile  bool
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"base 404", client, "", "", "", false, false},
		{"base ref is empty", client, "yrdy", "", "sfdf", false, false},
		{"base branch is Unauthorized", client, "yrdy", "yrdy", "sfdf", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FileExists(tt.project, tt.branch, tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.getFile {
				t.Errorf("Client.FileExists() = %v, want %v", got, tt.getFile)
			}
		})
	}
}
