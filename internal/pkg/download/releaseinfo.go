// github release info model
package download

//inner asset model
type Asset struct {
	Url                string `json:"url"`
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browse_download_url"`
}

//release model
type ReleaseInfo struct {
	Id     int     `json:"id"`
	Url    string  `json:"url"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}
