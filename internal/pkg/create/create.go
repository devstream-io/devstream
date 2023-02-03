package create

import (
	"fmt"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/create/param"
	"github.com/devstream-io/devstream/internal/pkg/plugin/argocdapp"
	general "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	"github.com/devstream-io/devstream/internal/pkg/plugin/reposcaffolding"
	"github.com/devstream-io/devstream/pkg/util/cli"
)

func Create() error {
	params, err := param.GetParams()
	if err != nil {
		return err
	}
	fmt.Printf("Start create app %s's stream...\n", params.GitHubRepo)
	return create(params)
}

// create will do following three things:
// 1. create repo by repoScaffolding
// 2. config GitHub actions for this repo
// 3. create Argo CD application for this repo
func create(params *param.Param) error {
	if err := createRepo(params); err != nil {
		return err
	}

	if err := createApp(params); err != nil {
		return err
	}

	// finalMessage is used to help user to vist this app
	finalMessage := `You can now connect to you app with:

  kubectl port-forward service/%s 8080:8080 -n default

Then you can visit this app by http://127.0.0.1:8080 in your browser.

Happy Hacking! 😊
`
	fmt.Printf(finalMessage, params.GitHubRepo)
	return nil
}

func createRepo(params *param.Param) error {
	repoOptions := configmanager.RawOptions{
		"owner":   params.GithubUsername,
		"name":    params.GitHubRepo,
		"scmType": "github",
		"token":   params.GithubToken,
	}

	// 1.create repo
	status := cli.StatusForPlugin()
	repoScaffoldingOptions := configmanager.RawOptions{
		"destinationRepo": repoOptions,
		"sourceRepo": configmanager.RawOptions{
			"url": params.RepoScaffoldingURL,
		},
	}
	status.Start("Creating repo from scaffolding 🖼")
	_, err := reposcaffolding.Create(repoScaffoldingOptions)
	status.End(err)
	if err != nil {
		return err
	}

	// 2.config ci
	ciOptions := configmanager.RawOptions{
		"scm": repoOptions,
		"pipeline": configmanager.RawOptions{
			"language": configmanager.RawOptions{
				"name":      params.Language,
				"framework": params.Framework,
			},
			"imageRepo": configmanager.RawOptions{
				"user":     params.DockerhubUsername,
				"password": params.DockerhubToken,
			},
		},
	}
	status.Start("Writing github action configuration ✍️ ")
	_, err = general.Create(ciOptions)
	status.End(err)
	status.Start("Waiting for github action finished 🐎")

	// 3.wait repo ci finished
	waitCIFinished()
	status.End(err)
	return err
}

func createApp(params *param.Param) error {
	status := cli.StatusForPlugin()
	argocdAppOption := configmanager.RawOptions{
		"app": configmanager.RawOptions{
			"name":      params.GitHubRepo,
			"namespace": "argocd",
		},
		"destination": configmanager.RawOptions{
			"server":    "https://kubernetes.default.svc",
			"namespace": "default",
		},
		"source": configmanager.RawOptions{
			"valuefile":  "values.yaml",
			"path":       fmt.Sprintf("helm/%s", params.GitHubRepo),
			"repoURL":    fmt.Sprintf("https://github.com/%s/%s", params.GithubUsername, params.GitHubRepo),
			"repoBranch": "main",
			"token":      params.GithubToken,
		},
		"imageRepo": configmanager.RawOptions{
			"user": params.DockerhubUsername,
		},
	}
	status.Start("Creating argocd app 🕹️")
	_, err := argocdapp.Create(argocdAppOption)
	status.End(err)
	status.Start("Waiting for app to running 🚀")
	// wait argocd app status to running
	waitAppUp()
	status.End(nil)
	return err
}

// TODO(steinliber): add logic to wait for ci finished
func waitCIFinished() {
	time.Sleep(70 * time.Second) // current github actions takes 62 seconds for finished
}

// TODO(steinliber): add logic to wait for pod start running
func waitAppUp() {
	time.Sleep(30 * time.Second)
}
