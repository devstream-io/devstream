package githubactions

import (
	"context"
	"github.com/mitchellh/mapstructure"
)

// Reinstall remove and set up GitHub Actions workflows.
func Reinstall(options *map[string]interface{}) (bool, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		return false, err
	}

	for _, pipeline := range workflows {
		param := &Param{
			&ctx,
			getGitHubClient(&ctx),
			&opt,
			&pipeline,
		}
		_, errRemove := removeFile(param)
		if errRemove != nil {
			return false, errRemove
		}

		_, errCreate := createFile(param)
		if errCreate != nil {
			return false, errCreate
		}
	}
	return true, nil
}
