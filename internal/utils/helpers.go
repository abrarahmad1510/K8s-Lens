package utils

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

// PrintInfo prints an info message
func PrintInfo(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\033[36mINFO:\033[0m %s\n", message)
}

// PrintSuccess prints a success message
func PrintSuccess(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\033[32mSUCCESS:\033[0m %s\n", message)
}

// PrintWarning prints a warning message
func PrintWarning(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\033[33mWARNING:\033[0m %s\n", message)
}

// PrintError prints an error message
func PrintError(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("\033[31mERROR:\033[0m %s\n", message)
}

// PrintSection prints a section header
func PrintSection(title string) {
	fmt.Printf("\n\033[1;34m=== %s ===\033[0m\n", title)
}

// GetGoVersion returns the Go version
func GetGoVersion() string {
	return runtime.Version()
}

// GetPlatform returns the platform information
func GetPlatform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

// Colorize returns a colored string
func Colorize(text string, colorCode string) string {
	switch colorCode {
	case "blue":
		return color.New(color.FgBlue).SprintFunc()(text)
	case "green":
		return color.New(color.FgGreen).SprintFunc()(text)
	case "red":
		return color.New(color.FgRed).SprintFunc()(text)
	case "yellow":
		return color.New(color.FgYellow).SprintFunc()(text)
	default:
		return text
	}
}
