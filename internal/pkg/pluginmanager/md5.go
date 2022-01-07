package pluginmanager

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
)

// ConntentMD5 gen from file
func LocalContentMD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
