package github

import (
	"fmt"
	"net/http"
	"testing"
)

type fileTest struct {
	baseTest
	content      []byte
	filePath     string
	targetBranch string
	wantErr      bool
}

func TestClient_CreateFile(t *testing.T) {
	mux, serverUrl, teardown := setup(t)
	defer teardown()
	tests := []fileTest{
		// TODO: Add test cases.
		{
			baseTest{"base ", getClientWithOption(
				t, &Option{Owner: "o", Repo: "r", Org: "or"}, serverUrl,
			),
				"/repos/or/r/contents/a", http.MethodPut, true, `{"message":"Initialize the repository","content":"Yw==","branch":"b"}`, ""},
			[]byte{'c'}, "a", "b", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux.HandleFunc(tt.registerUrl, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.wantMethod)
				fmt.Fprint(w, tt.respBody)
			})
			if err := tt.client.CreateFile(tt.content, tt.filePath, tt.targetBranch); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
