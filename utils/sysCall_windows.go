//go:build windows

package utils

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// getWindowsMetadata extracts creation time, last-write time, size, and a
// synthesized mode from Windows' file attribute data. Windows exposes file
// timestamps and attributes very differently from POSIX (no inode, no
// separate "ctime" in the Unix sense, no direct executable bit) so this is
// a distinct implementation rather than a shared code path with Unix.
func getPlatformMetadata(fileInfo os.FileInfo) (Metadata, error) {
	stat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return Metadata{}, fmt.Errorf("unable to read windows file attributes")
	}

	ctime := time.Unix(0, stat.CreationTime.Nanoseconds())
	mtime := time.Unix(0, stat.LastWriteTime.Nanoseconds())

	// Windows has no POSIX-style permission bits. Synthesize a mode value
	// using Go's portable os.FileMode so downstream code (GetFileMode) still
	// has something consistent to inspect.
	mode := uint32(fileInfo.Mode())

	return Metadata{
		Ctime: ctime,
		Mtime: mtime,
		Mode:  mode,
		Size:  fileInfo.Size(),
	}, nil
}
