package statemanager

type ComponentAction string

const (
	ActionInstall   ComponentAction = "Install"
	ActionReinstall ComponentAction = "Reinstall"
	ActionUninstall ComponentAction = "Uninstall"
)
