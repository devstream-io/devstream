package file

import (
	"os"
)

const (
	defaultTempName    = "pkg-util-file-create_"
	appNamePlaceHolder = "_app_name_"
)

// CopyFile will copy file content from src to dst
func CopyFile(src, dest string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytesRead, 0644)
}
