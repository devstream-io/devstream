package githubactions

import "log"

func IsHealthy(options *map[string]interface{}) (bool, error) {
	ghActions, err := NewGithubActions(options)
	if err != nil {
		return false, err
	}

	language := ghActions.GetLanguage()
	log.Printf("Language is: %s.", language.String())

	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(language.String())
	retMap, err := ghActions.VerifyWorkflows(ws)
	if err != nil {
		return false, err
	}

	errFlag := false
	for name, err := range retMap {
		if err != nil {
			errFlag = true
			log.Printf("The workflow/file %s is not ok: %s", name, err)
		}
		log.Printf("The workflow/file %s is ok", name)
	}
	if errFlag {
		return false, nil
	}

	return true, nil
}
