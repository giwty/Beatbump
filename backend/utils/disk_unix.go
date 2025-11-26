//go:build !windows
// +build !windows

package utils

import (
	"syscall"
)

// GetFreeDiskSpace returns the available free space in bytes for the given path.
func GetFreeDiskSpace(path string) (uint64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}

	// Available blocks * block size
	return uint64(stat.Bavail) * uint64(stat.Bsize), nil
}
