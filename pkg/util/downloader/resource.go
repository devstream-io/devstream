package downloader

import (
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"
)

const (
	getterTempDirName   = "download_getter_"
	githubReplacePrefix = "https://github.com"
)

// detect the type of source to download
var goGetterDetectors = []getter.Detector{
	new(getter.GitHubDetector),
	new(getter.GitLabDetector),
	new(getter.GitDetector),
	new(getter.FileDetector),
	new(getter.BitBucketDetector),
}

// decompressors used when encounter compressed file
var goGetterDecompressors = map[string]getter.Decompressor{
	"gz":     new(getter.GzipDecompressor),
	"zip":    new(getter.ZipDecompressor),
	"tar.gz": new(getter.TarGzipDecompressor),
	"tar.xz": new(getter.TarXzDecompressor),
}

// these func is used to get resource
var goGetterGetters = map[string]getter.Getter{
	"file":  new(getter.FileGetter),
	"git":   new(getter.GitGetter),
	"http":  getterHTTPGetter,
	"https": getterHTTPGetter,
}

var getterHTTPClient = cleanhttp.DefaultClient()

var getterHTTPGetter = &getter.HttpGetter{
	Client: getterHTTPClient,
	Netrc:  true,
}

type ResourceClient struct {
	Source string
	// the destination of resource
	Destination string
	// this variable is used for save downloadPath for delete later
	tempCreatePath string
}

// GetWithGoGetter will download any resource from resourceAddress into localPath
// if ResourceClient.Source is just local file or directory, just return source
// if ResourceClient.Destination is not set, this func will create a temporary directory for download
func (c *ResourceClient) GetWithGoGetter() (string, error) {
	// 1. get current pwd
	workDir, err := os.Getwd()
	if err != nil {
		log.Debugf("download getter: get pwd failed: %+v", err)
		return "", err
	}

	source := c.getResourceSource()
	destination, err := c.getDestination()
	if err != nil {
		log.Debugf("download getter: get destination failed: %+v", err)
	}

	// 2. if source is localPath, just return source path
	if isLocalPath(source, workDir) {
		log.Debugf("download getter: source %s is local path, just return", source)
		return source, nil
	}

	// 3. create temp dir for download
	log.Debugf("download getter: start fetching %s to %s", source, destination)
	client := getter.Client{
		Src: source,
		Dst: destination,
		Pwd: workDir,

		Mode: getter.ClientModeAny,

		Detectors:       goGetterDetectors,
		Decompressors:   goGetterDecompressors,
		Getters:         goGetterGetters,
		Insecure:        true,
		DisableSymlinks: true,
	}
	err = client.Get()
	if err != nil {
		log.Debugf("download getter: fetch resource failed: %+v", err)
		return "", err
	}
	log.Debugf("download getter: fetch %s to %s success", source, destination)
	return destination, nil
}

// if c.Destination is set, just return c.Destination
// else create a tempDir for destination
func (c *ResourceClient) getDestination() (string, error) {
	if c.Destination != "" {
		return c.Destination, nil
	}
	tempDir, err := os.MkdirTemp("", getterTempDirName)
	if err != nil {
		log.Debugf("ci create tempDir for get files failed: %+v", err)
		return "", err
	}
	c.tempCreatePath = tempDir
	return tempDir, nil
}

// CleanUp is used to delete tempDir created before
func (c *ResourceClient) CleanUp() {
	if c.tempCreatePath != "" {
		err := os.RemoveAll(c.tempCreatePath)
		if err != nil {
			log.Debugf("download getter: delete temp dir failed: %+v", err)
		}
	}
}

// check source is local files
func isLocalPath(path, workDir string) bool {
	// get full address of resourceAddress
	fullAddress, err := getter.Detect(path, workDir, goGetterDetectors)
	if err != nil {
		log.Debugf("download getter: detect source failed: %+v", err)
		return false
	}
	log.Debugf("download getter: get %s => %s", path, fullAddress)
	return strings.HasPrefix(fullAddress, "file://")
}

// for github repo, go-getter only support github.com/owner/repo fomat
func (c *ResourceClient) getResourceSource() string {
	if strings.HasPrefix(c.Source, githubReplacePrefix) {
		return strings.Replace(c.Source, "https://", "", 1)
	}
	return c.Source
}
