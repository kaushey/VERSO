package model

import (
	"fmt"
	"os"
	"verso/utils"
)

type Blob struct {
	Path    string
	Oid     string
	Mode    string
	Content string
	Name    string
}

func CreateBlob(path string, file os.DirEntry) *Blob {
	o := &Blob{}
	o.Name = file.Name()
	o.Path = path
	o.Mode = utils.GetFileMode(file)
	o.setContent()
	o.Oid = utils.CalculateSHA1(o.Content)
	return o
}

func (o *Blob) setContent() {
	content, err := os.ReadFile(o.Path)
	if err != nil {
		fmt.Print("Error Occured while reading", o.Path, err)
	}
	o.Content = fmt.Sprintf("%s %d\x00%s", "blob", len(content), string(content))
}

// HashFile computes the object id a file would have if staged right now,
// without needing an os.DirEntry. Used by `status` to detect modifications.
func HashFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	blobContent := fmt.Sprintf("%s %d\x00%s", "blob", len(content), string(content))
	return utils.CalculateSHA1(blobContent), nil
}
