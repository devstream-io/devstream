package jenkins

import (
	"context"

	"github.com/bndr/gojenkins"

	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

type MockClient struct {
	ExecuteScriptError             error
	GetFolderJobError              error
	GetFolderJobValue              *gojenkins.Job
	InstallPluginsIfNotExistsError error
	ConfigCascForRepoError         error
	DeleteJobError                 error
	ApplyDingTalkBotError          error
	CreatePasswordCredentialError  error
	CreateSSHKeyCredentialError    error
	CreateGiltabCredentialError    error
	CreateSecretCredentialError    error
	BasicInfo                      *JenkinsConfigOption
}

func (m *MockClient) ExecuteScript(string) (string, error) {
	if m.ExecuteScriptError != nil {
		return "", m.ExecuteScriptError
	}
	return "", nil
}
func (m *MockClient) GetFolderJob(string, string) (*gojenkins.Job, error) {
	if m.GetFolderJobError != nil {
		return nil, m.GetFolderJobError
	}
	if m.GetFolderJobValue != nil {
		return m.GetFolderJobValue, nil
	}
	return nil, nil
}
func (m *MockClient) DeleteJob(context.Context, string) (bool, error) {
	if m.DeleteJobError != nil {
		return false, m.DeleteJobError
	}
	return true, nil
}
func (m *MockClient) InstallPluginsIfNotExists([]*JenkinsPlugin) error {
	if m.InstallPluginsIfNotExistsError != nil {
		return m.InstallPluginsIfNotExistsError
	}
	return nil
}
func (m *MockClient) CreateGiltabCredential(string, string) error {
	if m.CreateGiltabCredentialError != nil {
		return m.CreateGiltabCredentialError
	}
	return nil
}
func (m *MockClient) CreateSecretCredential(string, string) error {
	if m.CreateSecretCredentialError != nil {
		return m.CreateSecretCredentialError
	}
	return nil
}
func (m *MockClient) ConfigCascForRepo(*RepoCascConfig) error {
	if m.ConfigCascForRepoError != nil {
		return m.ConfigCascForRepoError
	}
	return nil
}
func (m *MockClient) ApplyDingTalkBot(dingtalk.BotConfig) error {
	if m.ApplyDingTalkBotError != nil {
		return m.ApplyDingTalkBotError
	}
	return nil
}
func (m *MockClient) CreateSSHKeyCredential(id, userName, privateKey string) error {
	if m.CreateSSHKeyCredentialError != nil {
		return m.CreateSSHKeyCredentialError
	}
	return nil
}
func (m *MockClient) CreatePasswordCredential(id, userName, privateKey string) error {
	if m.CreatePasswordCredentialError != nil {
		return m.CreatePasswordCredentialError
	}
	return nil
}

func (m *MockClient) GetBasicInfo() *JenkinsConfigOption {
	return &JenkinsConfigOption{
		URL: "http://mock.exmaple.com",
	}
}
