package pluginengine

type CommandType string

const (
	CommandApply  CommandType = "apply"
	CommandDelete CommandType = "delete"

	askUserIfContinue string = "Continue? [y/n]"
)
