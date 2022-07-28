package gitlab

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) PushLocalPathToBranch(repoPath, branch, pathWithNamespace, commitMsg string) (bool, error) {
	var files = make(map[string][]byte)

	// 1. walk through files
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

		gitlabPath, _ := filepath.Rel(repoPath, path)
		files[gitlabPath] = content
		return nil
	}); err != nil {
		return false, err
	}

	//2. push repo to gitlab
	err := c.CommitMultipleFiles(pathWithNamespace, branch, commitMsg, files)
	needRollBack := false
	if err != nil {
		needRollBack = true
	}
	return needRollBack, err
}
