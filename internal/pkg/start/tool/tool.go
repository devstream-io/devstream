package tool

type existsFunc func() bool
type stopedFunc func() bool
type installFunc func() error
type startFunc func() error
type tool struct {
	Name   string
	Exists existsFunc
	// Stopped can be nil if it is not needed.
	// if Stopped != nil -> Start can't be nil too
	Stopped stopedFunc
	Install installFunc
	Start   startFunc
}

var Tools []tool

func init() {
	Tools = []tool{toolDocker, toolMinikube, toolHelm, toolArgocd}
}
