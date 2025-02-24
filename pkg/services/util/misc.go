package util

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// GenerateGUID generates a random GUID
func GenerateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetExecutablePath returns the path of the current executable
func GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// GetOSInfo returns information about the operating system
func GetOSInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

// GetTimestamp returns the current Unix timestamp in milliseconds
func GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Sleep sleeps for the specified duration in milliseconds
func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
