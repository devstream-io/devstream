package util

import "github.com/devstream-io/devstream/pkg/util/log"

func LogAWSError(err error) {
	if err == nil {
		return
	}
	log.Errorf("AWS error: %s", err)
}
