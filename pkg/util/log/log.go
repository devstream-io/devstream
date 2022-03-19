package log

import (
	"github.com/sirupsen/logrus"
)

var (
	debugLog     = &CliLoggerFormatter{showType: "debug"}
	infoLog      = &CliLoggerFormatter{showType: "info"}
	warnLog      = &CliLoggerFormatter{showType: "warn"}
	errorLog     = &CliLoggerFormatter{showType: "error"}
	fatalLog     = &CliLoggerFormatter{showType: "fatal"}
	successLog   = &CliLoggerFormatter{showType: "success"}
	separatorLog = &SeparatorFormatter{}
)

// Debugf log info with color,symbol and format for a notice
func Debugf(format string, args ...interface{}) {
	logrus.SetFormatter(debugLog)
	logrus.Debugf(format, args...)
}

// Debug log info with color and symbol, for a notice
func Debug(args ...interface{}) {
	logrus.SetFormatter(debugLog)
	logrus.Debug(args...)
}

// Infof log info with color,symbol and format for a notice
func Infof(format string, args ...interface{}) {
	logrus.SetFormatter(infoLog)
	logrus.Infof(format, args...)
}

// Info log info with color and symbol, for a notice
func Info(args ...interface{}) {
	logrus.SetFormatter(infoLog)
	logrus.Info(args...)
}

// Warnf log warn with color,symbol and format for a warning event
func Warnf(format string, args ...interface{}) {
	logrus.SetFormatter(warnLog)
	logrus.Warnf(format, args...)
}

// Warn log warn with color and symbol, for a warning event
func Warn(args ...interface{}) {
	logrus.SetFormatter(warnLog)
	logrus.Warn(args...)
}

// Errorf log error with color,symbol and format for a warning event
func Errorf(format string, args ...interface{}) {
	logrus.SetFormatter(errorLog)
	logrus.Errorf(format, args...)
}

// Error log error with color adn symbol for a warning event
func Error(args ...interface{}) {
	logrus.SetFormatter(errorLog)
	logrus.Error(args...)
}

// Fatalf log fatal with color,symbol and format for a fatal event
func Fatalf(format string, args ...interface{}) {
	logrus.SetFormatter(fatalLog)
	logrus.Fatalf(format, args...)
}

// Fatal log fatal with color and symbol for a fatal event
func Fatal(args ...interface{}) {
	logrus.SetFormatter(fatalLog)
	logrus.Fatal(args...)
}

// Successf log success with color,symbol and format for a success operation
func Successf(format string, args ...interface{}) {
	logrus.SetFormatter(successLog)
	logrus.Infof(format, args...)
}

// Success log success with color and symbol, for a success operation
func Success(args ...interface{}) {
	logrus.SetFormatter(successLog)
	logrus.Info(args...)
}

// Separatorf prints a line for separating as green color
func Separatorf(format string, args ...interface{}) {
	logrus.SetFormatter(separatorLog)
	logrus.Infof(format, args...)
}

// Separator prints a line for separating as green color
func Separator(args ...interface{}) {
	logrus.SetFormatter(separatorLog)
	logrus.Info(args...)
}
