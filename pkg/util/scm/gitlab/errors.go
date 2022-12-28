package gitlab

import (
	"errors"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
)

var (
	errRepoNotFound   pkgerror.ErrorMessage = "Project Not Found"
	errRepoExist      pkgerror.ErrorMessage = "{name: [has already been taken]}"
	errWebHookInvalid pkgerror.ErrorMessage = "invlid url given"
	errFileExist      pkgerror.ErrorMessage = "A file with this name already exists"
	errVariableExist  pkgerror.ErrorMessage = "has already been taken"
	errFileNotExist   pkgerror.ErrorMessage = "file with this name doesn't exist"
)

var errorMsgMap = map[pkgerror.ErrorMessage]string{
	errWebHookInvalid: "webhook config doesn't support local networks, and you should config gitlab or change jenkinsURL config. For more info, you can refer to https://docs.gitlab.com/ee/security/webhooks.html#allow-webhook-and-service-requests-to-local-network",
	errFileExist:      "file already exist",
}

func (c *Client) newModuleError(err error) error {
	var newError = err
	for k, v := range errorMsgMap {
		if strings.Contains(err.Error(), string(k)) {
			newError = errors.New(v)
		}
	}
	return pkgerror.NewErrorFromPlugin("gitlab", c.GetRepoPath(), newError)
}
