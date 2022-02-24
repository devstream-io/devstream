package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/log"
)

func Read(options *map[string]interface{}) (map[string]interface{}, error) {
	gis, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	api := gis.GetApi()
	log.Infof("API is: %s.", api.Name)

	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)
	retMap, err := gis.VerifyWorkflows(ws)
	if err != nil {
		return nil, err
	}

	errFlag := false
	for name, err := range retMap {
		if err != nil {
			errFlag = true
			log.Errorf("The workflow/file %s got some error: %s", name, err)
		}
		log.Infof("The workflow/file %s is ok", name)
	}
	if errFlag {
		return nil, nil
	}

	return gis.buildReadState(api)
}
