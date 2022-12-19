package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/google/go-github/v42/github"
)

type RepoFileStatus struct {
	Path   string
	SHA    string
	Branch string
}

func (f *RepoFileStatus) EncodeToGitHubContentOption(commitMsg string) *github.RepositoryContentFileOptions {
	return &github.RepositoryContentFileOptions{
		Message: github.String(commitMsg),
		SHA:     github.String(f.SHA),
		Branch:  github.String(f.Branch),
	}
}

func CalculateGitHubBlobSHA(fileContent []byte) string {
	p := fmt.Sprintf("blob %d\x00", len(fileContent))
	h := sha1.New()
	h.Write([]byte(p))
	h.Write([]byte(fileContent))
	return hex.EncodeToString(h.Sum(nil))
}
