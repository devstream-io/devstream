package template

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

// Getters

func FromLocalFile(filepath string) ContentGetter {
	return func() ([]byte, error) {
		return os.ReadFile(filepath)
	}
}

func FromContent(content string) ContentGetter {
	return func() ([]byte, error) {
		return []byte(content), nil
	}
}

func FromURL(url string) ContentGetter {
	return func() ([]byte, error) {
		return downloader.FetchContentFromURL(url)
	}
}

// Quick Calls

func (r *render) FromLocalFile(filepath string) *rendererWithGetter {
	return r.SetContentGetter(FromLocalFile(filepath))
}

func (r *render) FromContent(content string) *rendererWithGetter {
	return r.SetContentGetter(FromContent(content))
}

func (r *render) FromURL(url string) *rendererWithGetter {
	return r.SetContentGetter(FromURL(url))
}
