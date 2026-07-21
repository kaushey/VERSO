package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"verso/utils"
)

type Commit struct {
	Parent  string // Oid of parent commit
	TreeId  string // Oid of tree that commit points to
	Oid     string // Oid of the commit
	Author  Author
	Message string
	Content string // content to be stored
}

func CreateCommit(parent string, treeOid string, a Author, m string) *Commit {
	c := &Commit{}
	c.Parent = parent
	c.TreeId = treeOid
	c.Author = a
	c.Message = m
	c.setContent()
	c.Oid = utils.CalculateSHA1(c.Content)
	return c
}

// content of commit blob that is to be stored in database
func (c *Commit) setContent() {
	var s string
	if c.Parent != "" {
		s = fmt.Sprintf("parent %s\n", c.Parent)
	}

	c.Content = fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s", c.TreeId, s, c.Author.toStr(), c.Author.toStr(), c.Message)
}

// ParseCommit decodes a decompressed commit object (as produced by setContent)
// back into a Commit struct. Used by the `log` command to walk history.
func ParseCommit(oid string, content string) (*Commit, error) {
	c := &Commit{Oid: oid}

	headerEnd := strings.Index(content, "\n\n")
	if headerEnd == -1 {
		return nil, fmt.Errorf("malformed commit object %s: missing header/message separator", oid)
	}

	header := content[:headerEnd]
	c.Message = content[headerEnd+2:]

	for _, line := range strings.Split(header, "\n") {
		switch {
		case strings.HasPrefix(line, "tree "):
			c.TreeId = strings.TrimPrefix(line, "tree ")
		case strings.HasPrefix(line, "parent "):
			c.Parent = strings.TrimPrefix(line, "parent ")
		case strings.HasPrefix(line, "author "):
			raw := strings.TrimPrefix(line, "author ")
			c.Author = parseAuthorLine(raw)
		}
	}

	return c, nil
}

// parseAuthorLine turns "Name <email> <unixSeconds> <offset>" back into an Author.
func parseAuthorLine(raw string) Author {
	a := Author{}
	nameEnd := strings.Index(raw, "<")
	if nameEnd == -1 {
		a.Name = strings.TrimSpace(raw)
		return a
	}
	a.Name = strings.TrimSpace(raw[:nameEnd])

	rest := raw[nameEnd+1:]
	emailEnd := strings.Index(rest, ">")
	if emailEnd == -1 {
		return a
	}
	a.Email = rest[:emailEnd]

	fields := strings.Fields(rest[emailEnd+1:])
	if len(fields) > 0 {
		if sec, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
			a.T = time.Unix(sec, 0)
		}
	}
	return a
}
