package start

import "os/exec"

type existsFunc func() bool
type installFunc func() error
type tool struct {
	Name    string
	Exists  existsFunc
	Install installFunc
}

var tools = []tool{
	{
		Name: "Docker",
		Exists: func() bool {
			_, err := exec.LookPath("docker")
			return err == nil
		},
		Install: installDocker,
	},
	{
		Name: "Minikube",
		Exists: func() bool {
			_, err := exec.LookPath("minikube")
			return err == nil
		},
		Install: installMinikube,
	},
	{
		Name: "Helm",
		Exists: func() bool {
			_, err := exec.LookPath("helm")
			return err == nil
		},
		Install: installHelm,
	},
	{
		Name: "Argo CD",
		Exists: func() bool {
			// TODO(dh)
			return false
		},
		Install: installArgocd,
	},
}
