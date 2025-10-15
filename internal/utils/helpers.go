package utils

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var (
	Info    = color.New(color.FgBlue).PrintfFunc()
	Success = color.New(color.FgGreen).PrintfFunc()
	Warning = color.New(color.FgYellow).PrintfFunc()
	Error   = color.New(color.FgRed).PrintfFunc()
)

// PrintInfo Displays Informational Messages
func PrintInfo(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Info("INFO: %s\n", capitalizeFirst(message))
}

// PrintSuccess Displays Success Messages
func PrintSuccess(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Success("SUCCESS: %s\n", capitalizeFirst(message))
}

// PrintWarning Displays Warning Messages
func PrintWarning(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Warning("WARNING: %s\n", capitalizeFirst(message))
}

// PrintError Displays Error Messages
func PrintError(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	Error("ERROR: %s\n", capitalizeFirst(message))
}

// PrintSection Creates A Section Header
func PrintSection(title string) {
	cyan := color.New(color.FgCyan).Add(color.Bold)
	cyan.Printf("\n%s\n", strings.ToUpper(title))
	cyan.Println(strings.Repeat("=", len(title)))
}

// GetGoVersion Returns The Current Go Version
func GetGoVersion() string {
	return runtime.Version()
}

// GetPlatform Returns The Current Platform
func GetPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// CheckError Simplifies Error Handling
func CheckError(err error, message string) {
	if err != nil {
		PrintError("%s: %s", message, err)
	}
}

// FatalError Prints Error And Exits
func FatalError(err error, message string) {
	if err != nil {
		PrintError("%s: %s", message, err)
		panic(err)
	}
}

// CapitalizeFirst Capitalizes The First Letter Of A String
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FormatDuration Formats Duration In Human Readable Format
func FormatDuration(duration string) string {
	return strings.ToLower(duration)
}
