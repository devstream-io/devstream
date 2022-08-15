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
	return &localGetter{filepath: filepath}
}

func FromContent(content string) ContentGetter {
	return &contentGetter{content: content}
}

func FromURL(url string) ContentGetter {
	return &urlGetter{url: url}
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

// Getters definition

type localGetter struct {
	filepath string
}

func (g *localGetter) GetContent() ([]byte, error) {
	return os.ReadFile(g.filepath)
}

type contentGetter struct {
	content string
}

func (g *contentGetter) GetContent() ([]byte, error) {
	return []byte(g.content), nil
}

type urlGetter struct {
	url string
}

func (g *urlGetter) GetContent() ([]byte, error) {
	resp, err := http.Get(g.url)
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
