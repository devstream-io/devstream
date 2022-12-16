package scm

import (
	"crypto/md5"
	"encoding/hex"
)

// calculateSHA is used to calculate file content's md5
func calculateSHA(fileContent []byte) string {
	h := md5.New()
	h.Write(fileContent)
	return hex.EncodeToString(h.Sum(nil))
}
