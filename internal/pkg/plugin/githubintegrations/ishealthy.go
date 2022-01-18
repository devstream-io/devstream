package githubintegrations

import "log"

func IsHealthy(options *map[string]interface{}) (bool, error) {
	gis, err := NewGithubIntegrations(options)
	if err != nil {
		return false, err
	}

	api := gis.GetApi()
	log.Printf("api is: %s.", api.Name)

	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)
	retMap, err := gis.VerifyWorkflows(ws)
	if err != nil {
		return false, err
	}

	errFlag := false
	for name, err := range retMap {
		if err != nil {
			errFlag = true
			log.Printf("The workflow/file %s got some error: %s", name, err)
		}
		log.Printf("The workflow/file %s is ok", name)
	}
	if errFlag {
		return false, nil
	}

	return true, nil
}
