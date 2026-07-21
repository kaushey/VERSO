package utils

import "os"

// GetFileMetadata delegates to getPlatformMetadata, which has exactly one
// implementation compiled in per build (see sysCall_unix.go, sysCall_darwin.go,
// sysCall_windows.go). Each is guarded by a `//go:build` tag so the Go
// toolchain only ever compiles the one matching GOOS - this avoids the
// classic mistake of a runtime.GOOS switch calling into an OS-specific
// function that doesn't even exist in that build.
func GetFileMetadata(d os.DirEntry) (Metadata, error) {
	fileInfo, err := d.Info()
	if err != nil {
		return Metadata{}, err
	}
	return getPlatformMetadata(fileInfo)
}
