package utils

import (
	"fmt"
	"runtime"
)

// PrintInfo displays informational messages
func PrintInfo(format string, a ...interface{}) {
	fmt.Printf("INFO: "+format+"\n", a...)
}

// PrintSuccess displays success messages
func PrintSuccess(format string, a ...interface{}) {
	fmt.Printf("SUCCESS: "+format+"\n", a...)
}

// PrintWarning displays warning messages
func PrintWarning(format string, a ...interface{}) {
	fmt.Printf("WARNING: "+format+"\n", a...)
}

// PrintError displays error messages
func PrintError(format string, a ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", a...)
}

// GetGoVersion returns the current Go version
func GetGoVersion() string {
	return runtime.Version()
}

// GetPlatform returns the current platform
func GetPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// CheckError checks if error exists and prints it
func CheckError(err error, message string) {
	if err != nil {
		PrintError("%s: %s", message, err)
	}
}
