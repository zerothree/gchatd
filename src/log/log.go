// +build linux darwin

package log

import (
    "fmt"
    "log/syslog"
)

var (
    defaultLogger = syslog.New(0, "")
)

func Debug(format string, a ...interface{}) {
    defaultLogger.Debug(fmt.Sprintf(format, a...))
}

func Err(format string, a ...interface{}) {
    defaultLogger.Err(fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
    defaultLogger.Info(fmt.Sprintf(format, a...))
}

func Notice(format string, a ...interface{}) {
    defaultLogger.Notice(fmt.Sprintf(format, a...))
}

func Warning(format string, a ...interface{}) {
    defaultLogger.Warning(fmt.Sprintf(format, a...))
}
