// +build windows

package log

import (
    "log"
)

func Debug(format string, a ...interface{}) {
    log.Printf(fmt.Sprintf(format, a...))
}

func Err(format string, a ...interface{}) {
    log.Printf(fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
    log.Printf(fmt.Sprintf(format, a...))
}

func Notice(format string, a ...interface{}) {
    log.Printf(fmt.Sprintf(format, a...))
}

func Warning(format string, a ...interface{}) {
    log.Printf(fmt.Sprintf(format, a...))
}
