package java

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	client, err := opts.newGitlabClient()
	if err != nil {
		return false, err
	}
	commitInfo := &git.CommitInfo{
		CommitMsg:    commitMessage,
		CommitBranch: opts.Branch,
		GitFileMap:   git.GitFileContentMap{ciFileName: []byte("")},
	}
	if err = client.DeleteFiles(commitInfo); err != nil {
		return false, err
	}
	return false, nil
}
