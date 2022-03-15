package trellogithub

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/mapz"
	"github.com/merico-dev/stream/pkg/util/slicez"
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

func NewTrelloGithub(options map[string]interface{}) (*TrelloGithub, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(options, &opt)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opt); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	ghOptions := &github.Option{
		Owner:    opt.Owner,
		Repo:     opt.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	return &TrelloGithub{
		ctx:     ctx,
		client:  ghClient,
		options: &opt,
	}, nil
}

func (tg *TrelloGithub) GetApi() *Api {
	return tg.options.Api
}

// CompareFiles compare files between local and remote
func (tg *TrelloGithub) CompareFiles(wsFiles, filesInRemoteDir []string) map[string]error {
	lostFiles := slicez.SliceInSliceStr(wsFiles, filesInRemoteDir)
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
