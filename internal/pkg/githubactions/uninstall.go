package githubactions

import (
	"context"
	"github.com/mitchellh/mapstructure"
)

// Uninstall remove GitHub Actions workflows.
func Uninstall(options *map[string]interface{}) (bool, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		return false, err
	}

	for _, pipeline := range workflows {
		_, errRemove := removeFile(&Param{
			&ctx,
			getGitHubClient(&ctx),
			&opt,
			&pipeline,
		})
		if errRemove != nil {
			return false, errRemove
		}
	}
	return true, nil
}
