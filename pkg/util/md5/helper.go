package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// CalcFileMD5 calculate file md5
func CalcFileMD5(filename string) (string, error) {
	f, err := os.Open(filename)
	if nil != err {
		return "", err
	}
	defer f.Close()

	return CalcMD5(f)
}

func CalcMD5(r io.Reader) (string, error) {
	if r == nil {
		return "", fmt.Errorf("reader is nil")
	}
	md5Handle := md5.New()
	if _, err := io.Copy(md5Handle, r); nil != err {
		return "", err
	}
	md := md5Handle.Sum(nil)
	md5str := fmt.Sprintf("%x", md)
	return md5str, nil
}
