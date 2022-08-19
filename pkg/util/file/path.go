package file

import (
	"regexp"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func ReplaceAppNameInPathStr(filePath, appNamePlaceHolder, appName string) string {
	if !strings.Contains(filePath, appNamePlaceHolder) {
		return filePath
	}
	newFilePath := regexp.MustCompile(appNamePlaceHolder).ReplaceAllString(filePath, appName)
	log.Debugf("Replace file path place holder. Before: %s, after: %s.", filePath, newFilePath)
	return newFilePath
}
