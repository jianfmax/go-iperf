package iperf

import (
	"path"
	"runtime"
)

var (
	binaryLocation = ""
)

func init() {
	// Extract the binaries
	if runtime.GOOS == "windows" {
		InitWindowsPath("./")
	} else if runtime.GOOS == "darwin" {
		binaryLocation = "iperf3"
	} else {
		binaryLocation = "iperf3"
	}
}

func InitWindowsPath(filePath string) {
	binaryLocation = path.Join(filePath, "iperf3.exe")
}
