// 1. init download
// 2. get github releases
// 3. get assets .so

package download

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultReleaseUrl = "https://api.github.com/repos/merico-dev/stream/releases"
)

type DownloadClient struct {
	Filepath   string
	AssetName  string
	Version    string
	ReleaseUrl string
	client     *resty.Client
	RetryCount int
	WaitTime   int64
	TimeOut    int64
}

func NewDownloadClient() *DownloadClient {
	return &DownloadClient{
		ReleaseUrl: defaultReleaseUrl,
		client:     resty.New(),
		RetryCount: 3,
		WaitTime:   1000,
		TimeOut:    1000,
	}
}

func (dc *DownloadClient) GetReleaseDetail() ([]ReleaseInfo, error) {
	//github api get release details
	details, err := dc.client.R().Get(dc.ReleaseUrl)
	if err != nil {
		return nil, err
	}

	//unmarshal to json
	jsonDetails := []ReleaseInfo{}
	errums := json.Unmarshal(details.Body(), &jsonDetails)
	if errums != nil {
		return nil, errums
	}

	log.Print(jsonDetails)
	return jsonDetails, nil
}

func (dc *DownloadClient) GetAssetswithretry() error {

	//github api, to get release-info
	releasesInfo, err := dc.GetReleaseDetail()
	if err != nil {
		return err
	}

	//get the assert url from release-info
	githubApiUrl, _ := findAssetId(releasesInfo, dc.AssetName)

	fmt.Println(githubApiUrl)

	//download .so files
	realName := dc.Filepath
	tmpName := dc.Filepath + ".tmp"
	_, errifno := dc.client.R().SetOutput(tmpName).
		SetHeader("Accept", "application/octet-stream").
		Get(githubApiUrl)

	if errifno != nil {
		return errifno
	}

	//rename, tmp file to real file
	errename := os.Rename(tmpName, realName)
	if errename != nil {
		return errename
	}

	return nil
}

func findAssetId(releases []ReleaseInfo, assetName string) (assetUrl string, browserUrl string) {
	for _, releaseInfo := range releases {
		for _, asset := range releaseInfo.Assets {
			if asset.Name == assetName {
				return asset.Url, asset.BrowserDownloadUrl
			}
		}
	}
	return "", ""
}

// set resty http client param
func (dc *DownloadClient) SetRestyClientParam() {
	dc.client.SetRetryCount(dc.RetryCount).
		SetRetryWaitTime(time.Duration(dc.WaitTime) * time.Second).
		SetRetryMaxWaitTime(time.Duration(dc.TimeOut) * time.Second).
		AddRetryCondition(func(resp *resty.Response, e error) bool {
			if resp.RawResponse == nil || resp.StatusCode() == 0 ||
				(resp.StatusCode() >= http.StatusLocked && resp.StatusCode() < http.StatusNotExtended) {
				return true
			} else {
				return false
			}
		},
		)
}
