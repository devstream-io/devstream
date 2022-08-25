package jenkinspipelinekubernetes

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validateAndHandleOptions(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(&opts)
	if err != nil {
		return false, err
	}

	// delete the credentials created by devstream if exists
	if _, err := client.GetCredentialsUsername(jenkinsCredentialID); err == nil {
		if err := client.DeleteCredentialsUsername(jenkinsCredentialID); err != nil {
			return false, err
		}
	}

	// delete the job created by devstream if exists
	if _, err = client.GetJob(context.Background(), opts.JobName); err == nil {
		if _, err := client.DeleteJob(context.Background(), opts.JobName); err != nil {
			return false, err
		}
	}

	return true, nil
}
