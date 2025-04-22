//go:build windows
// +build windows
package utils

import (
	"os"
	"path/filepath"
)

var (
	// Default base directory for Windows
	baseDir = filepath.Join(os.Getenv("ProgramData"), "eraser", "shared-data")

	// Path constants for Windows
	ScanErasePath            = filepath.Join(baseDir, "scanErase")
	CollectScanPath          = filepath.Join(baseDir, "collectScan")
	EraseCompleteCollectPath = filepath.Join(baseDir, "eraseCompleteCollect")
	EraseCompleteScanPath    = filepath.Join(baseDir, "eraseCompleteScan")

	// Windows uses named pipes with different format
	CRIPath = `\\.\pipe\cri`
)

func init() {
	// Create base directory if it doesn't exist
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		_ = os.MkdirAll(baseDir, 0755)
	}
}
