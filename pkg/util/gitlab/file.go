package gitlab

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) PushLocalPath(repoPath, branch, pathWithNamespace, commitMsg string) error {
	var files = make(map[string][]byte)

	if err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}

		if info.IsDir() {
			log.Debugf("Found dir: %s.", path)
			return nil
		}

		log.Debugf("Found file: %s.", path)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		repoPath := strings.Join(strings.Split(path, "/")[2:], "/")
		files[repoPath] = content
		return nil
	}); err != nil {
		return err
	}
	return c.CommitMultipleFiles(pathWithNamespace, branch, commitMsg, files)
}
