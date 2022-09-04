package gitlab

import (
	"errors"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/pkgerror"
)

var (
	errRepoNotFound   pkgerror.ErrorMessage = "Project Not Found"
	errRepoExist      pkgerror.ErrorMessage = "name already exists"
	errWebHookInvalid pkgerror.ErrorMessage = "invlid url given"
	errFileExist      pkgerror.ErrorMessage = "A file with this name already exists"
)

var errorMsgMap = map[pkgerror.ErrorMessage]string{
	errWebHookInvalid: "webhook config doesn't support local network, should config gitlab or change jenkinsURL config",
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
