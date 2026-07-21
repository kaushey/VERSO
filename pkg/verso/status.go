package verso

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"verso/pkg/model"
	"verso/utils"
)

// StatusHandler compares the working directory against the staged index and
// reports untracked, modified, and deleted files - similar in spirit to
// `git status --short`.
func StatusHandler(args []string) error {
	file, err := os.OpenFile(path.Join(utils.VersoPath, "index"), os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening .verso/index: %v", err)
	}
	defer file.Close()

	index, err := model.ParseIndex(file)
	if err != nil {
		return fmt.Errorf("error parsing .verso/index: %v", err)
	}

	workspaceFiles, err := utils.Ls()
	if err != nil {
		return fmt.Errorf("error listing working directory: %v", err)
	}

	seen := make(map[string]bool)
	var untracked, modified, deleted []string

	for _, relPath := range workspaceFiles {
		seen[relPath] = true
		absPath := filepath.Join(utils.WorkindDir, relPath)

		entry, tracked := index.Entries[relPath]
		if !tracked {
			untracked = append(untracked, relPath)
			continue
		}

		currentOid, err := model.HashFile(absPath)
		if err != nil {
			return fmt.Errorf("error hashing %s: %v", relPath, err)
		}
		if currentOid != hex.EncodeToString([]byte(entry.EntryId())) {
			modified = append(modified, relPath)
		}
	}

	for relPath := range index.Entries {
		if !seen[relPath] {
			deleted = append(deleted, relPath)
		}
	}

	sort.Strings(untracked)
	sort.Strings(modified)
	sort.Strings(deleted)

	for _, p := range modified {
		fmt.Printf(" M %s\n", p)
	}
	for _, p := range deleted {
		fmt.Printf(" D %s\n", p)
	}
	for _, p := range untracked {
		fmt.Printf("?? %s\n", p)
	}

	if len(modified)+len(deleted)+len(untracked) == 0 {
		fmt.Println("nothing to commit, working tree clean")
	}

	return nil
}
