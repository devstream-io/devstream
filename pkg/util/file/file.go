package file

import (
	"io/ioutil"
)

func CopyFile(src, dest string) error {
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, bytesRead, 0644)
}
