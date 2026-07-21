package verso

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"verso/pkg/model"
	"verso/utils"
)

func CatFileHandler(args []string) error {
	obId := args[0]
	database, err := model.CreateDatabase(utils.WorkindDir)
	if err != nil {
		return err
	}
	if content, err := database.Read(obId); err != nil {
		return err
	} else {
		buffer := bytes.NewBuffer(content)
		r, err := zlib.NewReader(buffer)
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, r)
	}
	return nil
}
