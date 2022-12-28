package gitlab

import (
	"github.com/imdario/mergo"
	"github.com/xanzy/go-gitlab"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type commitTree struct {
	commitMessage string
	commitBranch  string
	gitlabFileMap map[gitlab.FileActionValue]git.GitFileContentMap
}

func newCommitTree(commitMessage string, branch string) *commitTree {
	return &commitTree{
		commitMessage: commitMessage,
		commitBranch:  branch,
		gitlabFileMap: make(map[gitlab.FileActionValue]git.GitFileContentMap),
	}
}

func (t *commitTree) addCommitFile(action gitlab.FileActionValue, scmPath string, content []byte) {
	actionMap, ok := t.gitlabFileMap[action]
	if !ok {
		t.gitlabFileMap[action] = git.GitFileContentMap{
			scmPath: content,
		}
	} else {
		actionMap[scmPath] = content
	}
}

func (t *commitTree) addCommitFilesFromMap(action gitlab.FileActionValue, gitMap git.GitFileContentMap) {
	actionMap, ok := t.gitlabFileMap[action]
	if !ok {
		t.gitlabFileMap[action] = gitMap
	} else {
		err := mergo.Merge(&actionMap, gitMap, mergo.WithOverride)
		if err != nil {
			log.Debugf("gitlab add commit files failed: %+v", err)
		}
	}
}

func (t *commitTree) getFilesCount() int {
	var fileCount int
	for _, files := range t.gitlabFileMap {
		fileCount += len(files)
	}
	return fileCount
}

func (t *commitTree) createCommitInfo() *gitlab.CreateCommitOptions {
	var commitActionsOptions = make([]*gitlab.CommitActionOptions, 0, t.getFilesCount())
	for action, fileMap := range t.gitlabFileMap {
		for fileName, content := range fileMap {
			commitActionsOptions = append(commitActionsOptions, &gitlab.CommitActionOptions{
				Action:   gitlab.FileAction(action),
				FilePath: gitlab.String(fileName),
				Content:  gitlab.String(string(content)),
			})
		}
	}
	return &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(t.commitBranch),
		CommitMessage: gitlab.String(t.commitMessage),
		Actions:       commitActionsOptions,
	}
}
