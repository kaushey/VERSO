package verso

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"

	"verso/pkg/model"
	"verso/utils"
)

// LogHandler walks the commit graph starting at HEAD, following each
// commit's parent pointer back to the root commit, and prints a short
// history similar to `git log --oneline`.
func LogHandler(args []string) error {
	refs := model.Refs{Path: utils.VersoPath}
	oid, err := refs.ReadHead()
	if err != nil {
		return fmt.Errorf("error reading HEAD: %v", err)
	}
	if oid == "" {
		fmt.Println("fatal: no commits yet")
		return nil
	}

	database, err := model.CreateDatabase(utils.WorkindDir)
	if err != nil {
		return fmt.Errorf("error in creating database: %v", err)
	}

	for oid != "" {
		raw, err := database.Read(oid)
		if err != nil {
			return fmt.Errorf("error reading object %s: %v", oid, err)
		}

		r, err := zlib.NewReader(bytes.NewReader(raw))
		if err != nil {
			return fmt.Errorf("error decompressing object %s: %v", oid, err)
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return fmt.Errorf("error reading object %s: %v", oid, err)
		}
		r.Close()

		// Strip the "commit <size>\x00" object header before parsing.
		content := buf.String()
		if idx := strings.IndexByte(content, 0); idx != -1 {
			content = content[idx+1:]
		}

		commit, err := model.ParseCommit(oid, content)
		if err != nil {
			return err
		}

		fmt.Printf("commit %s\n", commit.Oid)
		fmt.Printf("Author: %s <%s>\n", commit.Author.Name, commit.Author.Email)
		fmt.Printf("Date:   %s\n\n", commit.Author.T.Format("Mon Jan 2 15:04:05 2006 -0700"))
		for _, line := range strings.Split(strings.TrimRight(commit.Message, "\n"), "\n") {
			fmt.Printf("    %s\n", line)
		}
		fmt.Println()

		oid = commit.Parent
	}

	return nil
}
