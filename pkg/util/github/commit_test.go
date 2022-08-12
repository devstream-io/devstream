package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-github/v42/github"

	"github.com/devstream-io/devstream/pkg/util/repo"
)

type commitTest struct {
	BaseTest
	want    *github.RepositoryCommit
	wantErr bool
}

func TestClient_GetLastCommit(t *testing.T) {
	mux, serverUrl, teardown := Setup()
	defer teardown()
	sha := "s"
	tests := []commitTest{
		// TODO: Add test cases.
		{
			BaseTest{"base 200 ", GetClientWithOption(
				t, &repo.RepoInfo{Owner: "o", Repo: "r", Org: "or"}, serverUrl,
			),
				"/repos/or/r/commits", http.MethodGet, false, "", `[{"sha": "s"}]`},
			&github.RepositoryCommit{SHA: &sha}, false,
		},
		{
			BaseTest{"base 200 with empty result", GetClientWithOption(
				t, &repo.RepoInfo{Owner: "o", Repo: "r"}, serverUrl,
			),
				"/repos/o/r/commits", http.MethodGet, false, "", `[]`},
			nil, true,
		},
		{
			BaseTest{"base 404", GetClientWithOption(
				t, &repo.RepoInfo{Owner: "o"}, serverUrl,
			),
				"/aaa", http.MethodGet, false, "", ""},
			nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				DoTestMethod(t, r, tt.wantMethod)
				fmt.Fprint(w, tt.respBody)
			})
			got, err := tt.client.GetLastCommit()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLastCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetLastCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}
