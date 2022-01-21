package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Debug("nice to meet you ", "log")
	Info("nice to meet you ", "log")
	Warn("nice to meet you ", "log")
	Success("nice to meet you ", "log")

	Debugf("nice to meet you %s", "log")
	Infof("nice to meet you %s", "log")
	Warnf("nice to meet you %s", "log")
	Successf("nice to meet you %s", "log")

}
