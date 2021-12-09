package download

type DownloadManager interface {
	GetAssetswithretry() error
	GetReleaseDetail() (*[]ReleaseInfo, error)
}
