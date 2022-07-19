package gitlab

import (
	"os"
	"testing"
)

var fake_gitlab_token = "FAKE_GITLAB_TOKEN"
var baseURL = ""

func TestWithBaseURL(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	got := WithBaseURL(baseURL)
	if got == nil {
		t.Error("got must not be nil\n")
	}
	client, err := NewClient(got)
	if err != nil {
		t.Error(err)
	}
	if client.baseURL != baseURL {
		t.Errorf("client.baseURL = %s\n, want = %s\n", client.baseURL, baseURL)
	}
}

func TestNewClient(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()

	baseURL := ""
	opts := []OptionFunc{
		WithBaseURL(baseURL),
	}
	opts2 := []OptionFunc{
		WithBaseURL(DefaultGitlabHost),
	}
	testDataName := "last case: empty token"
	tests := []struct {
		name        string
		opts        []OptionFunc
		isNilClient bool
		wantErr     bool
	}{
		// TODO: Add test cases.
		{"base", opts, false, false},
		{"base not empty baseUrl", opts2, false, false},
		{testDataName, opts, true, true},
	}
	for _, tt := range tests {
		if tt.name == testDataName {
			os.Setenv("GITLAB_TOKEN", "")
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.isNilClient {
				t.Errorf("NewClient() = %+v\n, want %+v\n", got, tt.isNilClient)
			}
		})
	}
}
