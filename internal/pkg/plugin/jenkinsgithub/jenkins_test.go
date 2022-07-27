package jenkinsgithub

import (
	_ "embed"
)

//func initJenkins() *jenkins.Jenkins {
//	pwd := "L4VlROiaotUCRI4WwkrdR6"
//	client, err := jenkins.NewJenkins("http://localhost:32000", "admin", pwd)
//	if err != nil {
//		panic(err)
//	}
//
//	return client
//}

//func TestCreateCredential(t *testing.T) {
//	c := initJenkins()
//	err := c.CreateCredentialsUsername(jenkinsCredentialUsername, "", jenkinsCredentialID, jenkinsCredentialDesc)
//	if err != nil {
//		t.Errorf("CreateCredentialsUsername failed: %v", err)
//	}
//}
//
//func TestGetCredential(t *testing.T) {
//	c := initJenkins()
//	_, err := c.GetCredentialsUsername(jenkinsCredentialID)
//	if err != nil {
//		t.Errorf("GetCredentialsUsername failed: %v", err)
//	}
//}
//
//func TestHasPlugin(t *testing.T) {
//	j := initJenkins()
//	fmt.Println(j.HasPlugin(context.Background(), "ghprb"))
//}
//
//func TestInstallPlugin(t *testing.T) {
//	j := initJenkins()
//	fmt.Println(j.InstallPlugin(context.Background(), "antisamy-markup-formatter", "latest"))
//}
//
//func TestPlugins(t *testing.T) {
//	j := initJenkins()
//	fmt.Println(getPluginExistsMap(j))
//	fmt.Println(installPluginsIfNotExists(j))
//	fmt.Println(getPluginExistsMap(j))
//}
//
//func TestCreateConfigMapForJCasC(t *testing.T) {
//	opts := &GitHubIntegOptions{
//		AdminList:          []string{"aFlyBird0", "Bird"},
//		CredentialsID:      jenkinsCredentialID,
//		GithubAuthID:       githubAuthID,
//		JenkinsURLOverride: "https://891e-125-111-206-162.ap.ngrok.io/",
//	}
//
//	content, err := renderGitHubInteg(opts)
//	if err != nil {
//		t.Errorf("renderGitHubInteg failed: %v", err)
//	}
//
//	fmt.Println(applyJCasC("jenkins", "dev", githubIntegName, content))
//}
