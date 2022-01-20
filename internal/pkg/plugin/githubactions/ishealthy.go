package githubactions

import "github.com/merico-dev/stream/internal/pkg/util/log"

func IsHealthy(options *map[string]interface{}) (bool, error) {
	ghActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := ghActions.GetLanguage()
	log.Infof("Language is: %s.", language.String())

	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())
	retMap, err := ghActions.VerifyWorkflows(ws)
	if err != nil {
		return false, err
	}

	errFlag := false
	for name, err := range retMap {
		if err != nil {
			errFlag = true
			log.Errorf("The workflow/file %s is not ok: %s", name, err)
		}
		log.Successf("The workflow/file %s is ok", name)
	}
	if errFlag {
		return false, nil
	}

	return true, nil
}
