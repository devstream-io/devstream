package gitlab

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func TestClient_CreateProject(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		baseURL string
		Client  *Client
		opts    *CreateProjectOptions
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base", baseURL, client,
			&CreateProjectOptions{
				Namespace: "",
			}, true},
		{"base Namespace is not empty", baseURL, client,
			&CreateProjectOptions{
				Namespace: "test",
			}, true},
		{"base Visibility is `public`", baseURL, client,
			&CreateProjectOptions{
				Visibility: "public",
			}, true},
		{"base Visibility is `internal`", baseURL, client,
			&CreateProjectOptions{
				Visibility: "internal",
			}, true},
		{"base Visibility is `private`", baseURL, client,
			&CreateProjectOptions{
				Visibility: "private",
			}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.Client.CreateProject(tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteProject(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteProject("")
	assert.NotNil(t, err)

	err = client.DeleteProject("test")
	assert.NotNil(t, err)

}

func TestClient_DescribeProject(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", fake_gitlab_token)
	defer func() {
		os.Unsetenv("GITLAB_TOKEN")
	}()
	client, err := NewClient(WithBaseURL(baseURL))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		Client  *Client
		project string
		want    *gitlab.Project
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base", client, "", nil, true},
		{"base", client, "test", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.Client.DescribeProject(tt.project)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DescribeProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DescribeProject() = %v, want %v", got, tt.want)
			}
		})
	}
}
