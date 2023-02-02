package tool

type existsFunc func() bool
type stopedFunc func() bool
type installFunc func() error
type startFunc func() error
type tool struct {
	Name     string
	IfExists existsFunc
	Install  installFunc
	// IfStopped can be nil if it is not needed.
	// if IfStopped != nil -> Start can't be nil too
	IfStopped stopedFunc
	Start     startFunc
}

func GetTools() []tool {
	return []tool{toolDocker, toolMinikube, toolHelm, toolArgocd}
}
