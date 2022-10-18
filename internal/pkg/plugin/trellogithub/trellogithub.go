package trellogithub

import (
	"context"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

type TrelloGithub struct {
	ctx     context.Context
	client  *github.Client
	options *Options
}

type TrelloItemId struct {
	boardId     string
	todoListId  string
	doingListId string
	doneListId  string
}

func NewTrelloGithub(options configmanager.RawOptions) (*TrelloGithub, error) {
	ctx := context.Background()

	var opts Options
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	ghOptions := &git.RepoInfo{
		Owner:    opts.Owner,
		Org:      opts.Org,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	return &TrelloGithub{
		ctx:     ctx,
		client:  ghClient,
		options: &opts,
	}, nil
}

// CompareFiles compare files between local and remote
func (tg *TrelloGithub) CompareFiles(wsFiles, filesInRemoteDir []string) map[string]error {
	lostFiles := mapset.NewSet[string](wsFiles...).Difference(mapset.NewSet[string](filesInRemoteDir...)).ToSlice()
	// all files exist
	if len(lostFiles) == 0 {
		log.Info("All workflows exist.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
		return retMap
	}
	// some files lost
	retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
	for _, f := range lostFiles {
		log.Warnf("Lost file: %s.", f)
		retMap[f] = fmt.Errorf("not found")
	}
	return retMap
}
