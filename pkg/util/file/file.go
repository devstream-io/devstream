package file

import (
	"io"
	"os"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func CopyFile(srcFile, dstFile string) (err error) {
	// prepare source file
	sFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", srcFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// create destination file
	dFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", dstFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// copy and sync
	if _, err = io.Copy(dFile, sFile); err != nil {
		return nil
	}
	return dFile.Sync()
}
