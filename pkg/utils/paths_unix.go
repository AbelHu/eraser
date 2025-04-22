//go:build !windows
// +build !windows
package utils

const (
	ScanErasePath            = "/run/eraser.sh/shared-data/scanErase"
	CollectScanPath          = "/run/eraser.sh/shared-data/collectScan"
	EraseCompleteCollectPath = "/run/eraser.sh/shared-data/eraseCompleteCollect"
	EraseCompleteScanPath    = "/run/eraser.sh/shared-data/eraseCompleteScan"

	CRIPath = "/run/cri/cri.sock"
)
