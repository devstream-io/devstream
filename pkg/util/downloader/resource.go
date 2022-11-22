package downloader

import (
	"fmt"
	"os"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"

	"github.com/hashicorp/go-cleanhttp"
	getter "github.com/hashicorp/go-getter"
)

type (
	// ResourceLocation represent location of resource, url/localPath/gitPath for example
	ResourceLocation string
	// resourceCache is used to cache ResourceLocation/tempResourcePath map
	resourceCache map[ResourceLocation]string
)

var (
	// locationCache is used to cache resource
	locationCache = resourceCache{}
)

// config detectors|decompressors|getters for resource getter
func (l ResourceLocation) Download() (string, error) {
	// 1. check if has cache for ResourceLocation
	cachedLocation, exist := locationCache[l]
	if exist {
		return cachedLocation, nil
	}
	log.Debugf("start to get resource files [%s]...", l)

	// 2. get template Path where download resource
	resouceClient, err := newResourceClient(l.formatResourceSource())
	if err != nil {
		log.Debugf("create resourceClient failed: %+v", err)
	}
	if resouceClient.checkResourceExist() {
		return resouceClient.Src, nil
	}
	tempResourcePath, err := resouceClient.download()
	if err != nil {
		return "", fmt.Errorf("get resource files %s failed: %w", l, err)
	}

	// 3. set cache for ResourceLocation
	locationCache[l] = tempResourcePath
	return tempResourcePath, nil
}

// for github repo, go-getter only support github.com/owner/repo fomat
func (l ResourceLocation) formatResourceSource() string {
	const githubReplacePrefix = "https://github.com"
	locationString := string(l)
	if strings.HasPrefix(locationString, githubReplacePrefix) {
		return strings.Replace(locationString, "https://", "", 1)
	}
	return locationString
}

// resouceClient is used to download resource
type resouceClient struct {
	*getter.Client
}

// downloadToDst will download any resource from src location into dst location
func (c *resouceClient) download() (string, error) {
	// 1. create template dir by pattern
	const tempDirNamePattern = "download_getter_"
	dst, err := file.CreateTempDir(tempDirNamePattern)
	if err != nil {
		log.Debugf("download getter: get destination failed: %+v", err)
	}

	// 2. download resource to tempDir
	log.Debugf("download getter: start fetching %s to %s", c.Src, dst)
	c.setDst(dst)
	err = c.Get()
	if err != nil {
		log.Debugf("download getter: fetch resource failed: %+v", err)
		return "", err
	}
	log.Debugf("download getter: fetch %s to %s success", c.Src, dst)
	return dst, nil
}

func newResourceClient(source string) (*resouceClient, error) {
	// 1. config detectors/decompressors/getters for client
	var (
		// getterHTTPGetter is used to get resource from http
		getterHTTPGetter = &getter.HttpGetter{
			Client: cleanhttp.DefaultClient(),
			Netrc:  true,
		}
		// detect the type of source to download
		goGetterDetectors = []getter.Detector{
			new(getter.GitHubDetector),
			new(getter.GitLabDetector),
			new(getter.GitDetector),
			new(getter.FileDetector),
			new(getter.BitBucketDetector),
		}
		// decompressors used when encounter compressed file
		goGetterDecompressors = map[string]getter.Decompressor{
			"gz":     new(getter.GzipDecompressor),
			"zip":    new(getter.ZipDecompressor),
			"tar.gz": new(getter.TarGzipDecompressor),
			"tar.xz": new(getter.TarXzDecompressor),
		}
		// these func is used to get resource
		goGetterGetters = map[string]getter.Getter{
			"file":  new(getter.FileGetter),
			"git":   new(getter.GitGetter),
			"http":  getterHTTPGetter,
			"https": getterHTTPGetter,
		}
	)
	// 2. get current work dir
	workDir, err := os.Getwd()
	if err != nil {
		log.Debugf("download getter: get pwd failed: %+v", err)
		return nil, err
	}
	return &resouceClient{
		&getter.Client{
			Pwd:             workDir,
			Detectors:       goGetterDetectors,
			Decompressors:   goGetterDecompressors,
			Getters:         goGetterGetters,
			Insecure:        true,
			DisableSymlinks: true,
			Mode:            getter.ClientModeAny,
			Src:             source,
		},
	}, nil
}

// check source is local files
// if source is just local file or directory, just return source
func (c *resouceClient) checkResourceExist() bool {
	// get full address of resourceAddress
	fullAddress, err := getter.Detect(c.Src, c.Pwd, c.Detectors)
	if err != nil {
		log.Debugf("download getter: detect source failed: %+v", err)
		return false
	}
	log.Debugf("download getter: get %s => %s", c.Src, fullAddress)
	return strings.HasPrefix(fullAddress, "file://")
}

func (c *resouceClient) setDst(dst string) {
	c.Client.Dst = dst
}
