package infra

import (
	"fmt"
	"os"
	"runtime"
)

const (
	// InfoColor console color
	InfoColor = "\033[1;34m%s\033[0m"
	// ErrorColor console color
	ErrorColor = "\033[1;31m%s\033[0m"
	// WarningColor cosole color
	WarningColor = "\033[1;33m%s\033[0m"
)

// Warn logs to stdout in a yellow color if the terminal supports it
func Warn(format string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	userMsg := fmt.Sprintf(format, args...)
	if ok && details != nil {
		warning := fmt.Sprintf(WarningColor, fmt.Sprintf("Warning from Func %s and Line %d", details.Name(), line))
		fmt.Printf("%s : %s\n", warning, userMsg)
	} else {
		fmt.Print(userMsg)
	}
}

// Error logs to strerr in a red color if the terminal supports it
func Error(format string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	userMsg := fmt.Sprintf(format, args...)
	if ok && details != nil {
		err := fmt.Sprintf(ErrorColor, fmt.Sprintf("Error from Func %s and Line %d", details.Name(), line))
		fmt.Fprintf(os.Stderr, "%s : %s\n", err, userMsg)
	} else {
		fmt.Print(userMsg)
	}
}

// Info logs to stdout in a blue color if the terminal supports it
func Info(format string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	userMsg := fmt.Sprintf(format, args...)
	if ok && details != nil {
		info := fmt.Sprintf(InfoColor, fmt.Sprintf("Info from Func %s and Line %d", details.Name(), line))
		fmt.Printf("%s : %s\n", info, userMsg)
	} else {
		fmt.Print(userMsg)
	}
}
