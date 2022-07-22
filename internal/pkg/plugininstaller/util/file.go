package util

import (
	"io/ioutil"
)

func CopyFile(src, dest string) error {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dest, bytesRead, 0644); err != nil {
		return err
	}
	return nil
}
