package template

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

// ContentGetter gets content from any source
type getFunc func(string) ([]byte, error)

// LocalFileGetter get content bytes from file
func LocalFileGetter(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}

// ContentGetter get content bytes from input content
func ContentGetter(content string) ([]byte, error) {
	return []byte(content), nil
}

// URLGetter get content bytes from remote url
func URLGetter(url string) ([]byte, error) {
	return downloader.FetchContentFromURL(url)
}
