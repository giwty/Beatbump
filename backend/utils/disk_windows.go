package utils

import (
	"golang.org/x/sys/windows"
)

// GetFreeDiskSpace returns the available free space in bytes for the given path.
func GetFreeDiskSpace(path string) (uint64, error) {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64

	// GetDiskFreeSpaceEx requires a pointer to a UTF-16 string
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	err = windows.GetDiskFreeSpaceEx(pathPtr, &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return 0, err
	}

	return freeBytesAvailable, nil
}
