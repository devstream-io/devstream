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

// Success log success with color and symbol, for a success operation
func Success(format string, args ...interface{}) {
	cligger.Success(color.Success.Render(SUCCESS)+"%s", fmt.Sprint(args...))
}

// Info log info with color and symbol, for a notice
func Info(format string, args ...interface{}) {
	cligger.Info(color.Blue.Render(INFO)+"%s", fmt.Sprint(args...))
}

// Warn log warn with color and symbol, for a warning event
func Warn(format string, args ...interface{}) {
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

//func main() {
//	Info("test %s", "haha")
//	Warn("test %s", "haha")
//	Success("test %s", "haha")
//	Error("test %s", "haha")
//	Fatalf("test %s", "haha")
//}
