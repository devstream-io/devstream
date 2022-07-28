package github

import (
	"fmt"
	"net/http"
	"testing"
)

type newBranchTest struct {
	BaseTest
	baseBranch string
	newBranch  string
	wantErr    bool
}

type delBranchTest struct {
	BaseTest
	branch  string
	wantErr bool
}

func TestClient_NewBranch(t *testing.T) {
	respBody := `
	{
		"ref": "refs/heads/b",
		"url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		"object": {
			"type": "commit",
			"sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
			"url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		}
	}`
	mux, serverUrl, teardown := Setup()
	defer teardown()

	tests := []newBranchTest{
		// TODO: Add test cases.
		{
			BaseTest{"base", GetClientWithOption(
				t, &Option{Owner: "o", Repo: "r", Org: "or"}, serverUrl,
			),
				"/repos/or/r/git/ref/heads/b", http.MethodGet, false, "", ""},
			"b", "", true,
		},
		{
			BaseTest{"base set wrong register url for GetRef api in mock server", GetClientWithOption(
				t, &Option{Owner: "o", Repo: "r"}, serverUrl,
			),
				"repos", http.MethodGet, false, "", ""},
			"b", "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				t.Logf("test name: %s, hit path: %s", tt.name, r.URL.Path)
				fmt.Fprint(w, respBody)
			})
			if err := tt.client.NewBranch(tt.baseBranch, tt.newBranch); (err != nil) != tt.wantErr {
				t.Errorf("Client.NewBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteBranch(t *testing.T) {
	mux, serverUrl, teardown := Setup()
	defer teardown()

	tests := []delBranchTest{
		// TODO: Add test cases.
		{
			BaseTest{"base", GetClientWithOption(
				t, &Option{Owner: "o", Repo: "r", Org: "or"}, serverUrl,
			),
				"/repos/or/r/git/ref/heads/b", http.MethodGet, false, "", ""},
			"b", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			err := tt.client.DeleteBranch(tt.branch)
			t.Log(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
