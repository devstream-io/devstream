package scm

import (
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type MockScmClient struct {
	InitRaiseError         error
	PushRaiseError         error
	GetPathInfoError       error
	DownloadRepoError      error
	AddRepoSecretError     error
	DownloadRepoValue      string
	NeedRollBack           bool
	DeleteFuncIsRun        bool
	GetPathInfoReturnValue []*git.RepoFileStatus
}

func (m *MockScmClient) InitRepo() error {
	if m.InitRaiseError != nil {
		return m.InitRaiseError
	}
	return nil
}
func (m *MockScmClient) PushFiles(commitInfo *git.CommitInfo, checkUpdate bool) (bool, error) {
	if m.PushRaiseError != nil {
		return m.NeedRollBack, m.PushRaiseError
	}
	return m.NeedRollBack, nil
}
func (m *MockScmClient) DeleteRepo() error {
	m.DeleteFuncIsRun = true
	return nil
}
func (m *MockScmClient) GetPathInfo(path string) ([]*git.RepoFileStatus, error) {
	if m.GetPathInfoError != nil {
		return nil, m.GetPathInfoError
	}
	if m.GetPathInfoReturnValue != nil {
		return m.GetPathInfoReturnValue, nil
	}
	return nil, nil
}
func (m *MockScmClient) DeleteFiles(commitInfo *git.CommitInfo) error {
	return nil
}
func (m *MockScmClient) AddWebhook(webhookConfig *git.WebhookConfig) error {
	return nil
}
func (m *MockScmClient) DeleteWebhook(webhookConfig *git.WebhookConfig) error {
	return nil
}
func (m *MockScmClient) DownloadRepo() (string, error) {
	if m.DownloadRepoError != nil {
		return "", m.DownloadRepoError
	}
	return m.DownloadRepoValue, nil
}

func (m *MockScmClient) DescribeRepo() (*git.RepoInfo, error) {
	return nil, nil
}

func (m *MockScmClient) AddRepoSecret(secretKey, secretValue string) error {
	return m.AddRepoSecretError
}
