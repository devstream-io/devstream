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

func Success(format string, args ...interface{}) {
	cligger.Success(color.Success.Render(SUCCESS)+"%s", fmt.Sprint(args...))
}

func Info(format string, args ...interface{}) {
	cligger.Info(color.Blue.Render(INFO)+"%s", fmt.Sprint(args...))
}

func Warn(format string, args ...interface{}) {
	cligger.Warning(color.Warn.Render(WARN)+"%s", fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	cligger.Error(color.Error.Render(ERROR)+format, args...)
}

func Error(args ...interface{}) {
	cligger.Error(color.Error.Render(ERROR)+"%s", fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	cligger.Fatal(color.BgRed.Render(FATAL)+format, args...)
}

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
