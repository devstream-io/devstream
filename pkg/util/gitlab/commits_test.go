package gitlab

import (
	"os"
	"testing"
)

func TestClient_CommitSingleFile(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		client        *Client
		project       string
		branch        string
		commitMessage string
		filename      string
		content       string
		wantErr       bool
	}{
		// TODO: Add test cases.
		{"base", client, "", "", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.CommitSingleFile(tt.project, tt.branch, tt.commitMessage, tt.filename, tt.content); (err != nil) != tt.wantErr {
				t.Errorf("Client.CommitSingleFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteSingleFile(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		client        *Client
		project       string
		branch        string
		commitMessage string
		filename      string
		wantErr       bool
	}{
		// TODO: Add test cases.
		{"base", client, "", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.DeleteSingleFile(tt.project, tt.branch, tt.commitMessage, tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteSingleFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateSingleFile(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		client        *Client
		project       string
		branch        string
		commitMessage string
		filename      string
		content       string
		wantErr       bool
	}{
		// TODO: Add test cases.
		{"base", client, "", "", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.UpdateSingleFile(tt.project, tt.branch, tt.commitMessage, tt.filename, tt.content); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateSingleFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CommitMultipleFiles(t *testing.T) {
	m := map[string][]byte{
		"a.txt": []byte("a.txt"),
		"b.txt": []byte("b.txt"),
		"c.txt": []byte("c.txt"),
	}
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name          string
		client        *Client
		project       string
		branch        string
		commitMessage string
		files         map[string][]byte
		wantErr       bool
	}{
		// TODO: Add test cases.
		{"base", client, "", "", "", m, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.CommitMultipleFiles(tt.project, tt.branch, tt.commitMessage, tt.files); (err != nil) != tt.wantErr {
				t.Errorf("Client.CommitMultipleFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
