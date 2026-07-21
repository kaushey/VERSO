package utils

import (
	"os"
	"path"
)

// Directories to ignore
var IgnoreDirs = []string{".verso"}

var (
	VersoPath   string
	WorkindDir string
)

func Init() error {
	if wd, err := os.Getwd(); err != nil {
		return err
	} else {
		WorkindDir = wd
	}
	VersoPath = path.Join(WorkindDir, ".verso")
	return nil
}
