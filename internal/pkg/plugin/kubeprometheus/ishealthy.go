package kubeprometheus

// IsHealthy check the health for kube-prometheus with provided options.
// TODO(daniel-hutao): I'll implement the IsHealthy() function after
// the Install()/Uninstall()/Reinstall() is merged.
func IsHealthy(options *map[string]interface{}) (bool, error) {
	return true, nil
}
