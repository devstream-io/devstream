package version

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// CalcFileMD5 calculate file md5
func calcFileMD5(filename string) (string, error) {
	f, err := os.Open(filename)
	if nil != err {
		return "", err
	}
	defer f.Close()

	md5Handle := md5.New()
	if _, err := io.Copy(md5Handle, f); nil != err {
		return "", err
	}
	md := md5Handle.Sum(nil)
	md5str := fmt.Sprintf("%x", md)
	return md5str, nil
}