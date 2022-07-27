package jenkinsgithub

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validateAndHandleOptions(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(&opts)
	if err != nil {
		return nil, err
	}

	res := &resource{}

	if _, err = client.GetCredentialsUsername(jenkinsCredentialID); err == nil {
		res.CredentialsCreated = true
	}

	if _, err = client.GetJob(context.Background(), opts.J.JobName); err == nil {
		res.JobCreated = true
	}

	pluginMap, err := getPluginExistsMap(client)
	if err != nil {
		return nil, err
	}

	return buildResource(res, pluginMap), nil
}
