package s3

import (
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func validate(bucket, region, key string) error {
	var errs []error
	if bucket == "" {
		errs = append(errs, fmt.Errorf("state s3 Bucket is empty"))
	}
	if region == "" {
		errs = append(errs, fmt.Errorf("state s3 Region is empty"))
	}
	if key == "" {
		errs = append(errs, fmt.Errorf("state s3 Key is empty"))
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("s3 Backend config error")
	}

	return nil
}
