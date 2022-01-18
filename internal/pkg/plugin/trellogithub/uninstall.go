package trellogithub

import (
	"log"
)

// Uninstall remove GitHub Actions workflows.
func Uninstall(options *map[string]interface{}) (bool, error) {
	gis, err := NewTrelloGithub(options)
	if err != nil {
		return false, err
	}

	api := gis.GetApi()
	log.Printf("api is %s", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		err := gis.DeleteWorkflow(w)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
