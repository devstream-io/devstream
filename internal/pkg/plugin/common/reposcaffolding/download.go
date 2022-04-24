package reposcaffolding

import "github.com/devstream-io/devstream/pkg/util/github"

func download(owner, repo, workpath string) error {
	ghOption := &github.Option{
		Owner:    owner,
		Repo:     repo,
		NeedAuth: false,
		WorkPath: workpath,
	}
	ghClient, err := github.NewClient(ghOption)
	if err != nil {
		return err
	}

	if err = ghClient.DownloadLatestCodeAsZipFile(); err != nil {
		return err
	}

	return nil
}
