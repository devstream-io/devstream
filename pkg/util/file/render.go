package file

import (
	"regexp"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func replaceAppNameInPathStr(filePath, appNamePlaceHolder, appName string) (string, error) {
	if !strings.Contains(filePath, appNamePlaceHolder) {
		return filePath, nil
	}
	newFilePath := regexp.MustCompile(appNamePlaceHolder).ReplaceAllString(filePath, appName)
	log.Debugf("Replace file path place holder. Before: %s, after: %s.", filePath, newFilePath)
	return newFilePath, nil
}
