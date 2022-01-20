package log

import (
	"fmt"

	"github.com/guumaster/cligger"
	"gopkg.in/gookit/color.v1"
)

const (
	WARN    = "[WARN] "
	INFO    = "[INFO] "
	SUCCESS = "[SUCCESS] "
	ERROR   = "[ERROR] "
	FATAL   = "[FATAL] "
)

// Successf log success with color,symbol and format for a success operation
func Successf(format string, args ...interface{}) {
	cligger.Success(color.Success.Render(SUCCESS)+format, args...)
}

// Success log success with color and symbol, for a success operation
func Success(args ...interface{}) {
	cligger.Success(color.Success.Render(SUCCESS)+"%s", fmt.Sprint(args...))
}

// Infof log info with color,symbol and format for a notice
func Infof(format string, args ...interface{}) {
	cligger.Info(color.Blue.Render(INFO)+format, args...)
}

// Info log info with color and symbol, for a notice
func Info(args ...interface{}) {
	cligger.Info(color.Blue.Render(INFO)+"%s", fmt.Sprint(args...))
}

// Warnf log warn with color,symbol and format for a warning event
func Warnf(format string, args ...interface{}) {
	cligger.Warning(color.Warn.Render(WARN)+format, args...)
}

// Warn log warn with color and symbol, for a warning event
func Warn(args ...interface{}) {
	cligger.Warning(color.Warn.Render(WARN)+"%s", fmt.Sprint(args...))
}

// Errorf log error with color,symbol and format for a warning event
func Errorf(format string, args ...interface{}) {
	cligger.Error(color.Error.Render(ERROR)+format, args...)
}

// Error log error with color adn symbol for a warning event
func Error(args ...interface{}) {
	cligger.Error(color.Error.Render(ERROR)+"%s", fmt.Sprint(args...))
}

// Fatalf log fatal with color,symbol and format for a fatal event
func Fatalf(format string, args ...interface{}) {
	cligger.Fatal(color.BgRed.Render(FATAL)+format, args...)
}

// Fatal log fatal with color and symbol for a fatal event
func Fatal(args ...interface{}) {
	cligger.Fatal(color.BgRed.Render(FATAL)+"%s", fmt.Sprint(args...))
}
