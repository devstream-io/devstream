package docker

import (
	"fmt"
	"os"
)

// RemoveDirs removes the all the directories in the given list recursively
func RemoveDirs(dirs []string) []error {
	var errs []error
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove data %v: %v", dir, err))
		}
	}

	return errs
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
