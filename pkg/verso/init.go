package verso

import (
	"fmt"
	"os"
	"path"
	"verso/utils"
)

// TODO: Change Permission Mode for the directories and files
func InitHandler(args []string) error {

	if err := os.MkdirAll(utils.VersoPath, os.ModePerm); err != nil {
		return fmt.Errorf("error in creating .verso: %v", err)
	}
	subdirs := []string{"objects", "refs"}
	for i := range subdirs {
		if err := os.MkdirAll(path.Join(utils.VersoPath, subdirs[i]), os.ModePerm); err != nil {
			return fmt.Errorf("error in creating %s; %v", subdirs, err)
		}
	}

	headPath := path.Join(utils.VersoPath, "HEAD")
	if err := os.WriteFile(headPath, []byte{}, 0664); err != nil {
		return fmt.Errorf("error in creating HEAD: %v", err)
	}

	fmt.Printf("Initialized empty verso repository in %s\n", utils.VersoPath)
	return nil
}
