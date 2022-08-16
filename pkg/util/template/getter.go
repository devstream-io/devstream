package template

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/devstream-io/devstream/pkg/util/log"
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
		return getContentFromURL(url)
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

func getContentFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Close response body failed: %s", err)
		}
	}(resp.Body)

	// check response error
	if err != nil {
		log.Debugf("Download file from url failed: %s", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Download file from url failed: %+v", resp)
	}

	return io.ReadAll(resp.Body)
}
