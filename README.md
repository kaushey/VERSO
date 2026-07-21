# verso

A lightweight, Git-compatible version control system implemented from scratch in Go.

`verso` reimplements the core of Git's plumbing layer — content-addressable storage,
the staging index, tree/commit objects, and history traversal — to understand (and
demonstrate) how Git actually works under the hood, rather than just calling `git`.

## Why

Git's UI is simple, but its internals are a genuinely interesting piece of systems
design: a content-addressed object store, a binary staging index with a checksum
footer, and an immutable commit DAG. Building it from the ground up is a good way to
learn how hashing, compression, and binary file formats come together in a real
tool people use every day.

## Features

| Command    | What it does                                                              |
|------------|----------------------------------------------------------------------------|
| `init`     | Creates a new `.verso` repository (objects store, refs, HEAD)             |
| `add`      | Stages one or more files/directories into the binary index                |
| `commit`   | Snapshots the staged tree into an immutable commit object                 |
| `log`      | Walks the commit graph from HEAD and prints history                       |
| `status`   | Diffs the working directory against the index (untracked/modified/deleted)|
| `cat-file` | Reads and decompresses a raw object by its SHA-1 id                       |

### How objects are stored

Every blob, tree, and commit is:
1. Serialized into Git's plain-text object format (`<type> <size>\0<content>`)
2. Hashed with SHA-1 to get its object id (oid)
3. Compressed with zlib
4. Written to `.verso/objects/<first 2 chars of oid>/<remaining 38 chars>`

This is the same layout Git itself uses, which is what makes the on-disk objects
interoperable at the storage level.

### The index

`.verso/index` is a binary file (not JSON/text) with a 12-byte header, one
fixed-size metadata + variable-length-name entry per staged file, 8-byte alignment
padding, and a SHA-1 checksum trailer — mirroring Git's real index format rather
than a simplified stand-in.

## Building

Requires Go 1.22+. Builds natively on Linux, macOS, and Windows.

```bash
git clone https://github.com/<your-username>/verso.git
cd verso
make build
```

This produces a `bin/verso` binary (`bin/verso.exe` on Windows).

**On Windows**, `make` isn't available by default. Either install it (via
[Chocolatey](https://chocolatey.org/): `choco install make`, or use WSL/Git Bash), or
just run the build command directly:
```powershell
go build -o bin\verso.exe .\cmd\verso\main.go
```
Windows PowerShell users should set author env vars with `$env:` instead of `export`:
```powershell
$env:VERSO_AUTHOR_NAME = "Your Name"
$env:VERSO_AUTHOR_EMAIL = "you@example.com"
```

Cross-compiling from any OS to any target also works out of the box, e.g. to build a
Windows binary from Linux/macOS:
```bash
GOOS=windows GOARCH=amd64 go build -o bin/verso.exe cmd/verso/main.go
```

## Trying it out

```bash
export VERSO_AUTHOR_NAME="Your Name"
export VERSO_AUTHOR_EMAIL="you@example.com"

mkdir /tmp/demo && cd /tmp/demo
/path/to/bin/verso init

echo "hello" > file1.txt
/path/to/bin/verso status        # -> ?? file1.txt (untracked)

/path/to/bin/verso add file1.txt
/path/to/bin/verso commit -m "Initial commit"

/path/to/bin/verso log           # -> full commit history

echo "hello again" > file1.txt
/path/to/bin/verso status        # -> M file1.txt (modified)
```

Add `bin/` to your `PATH` (see `scripts/README.md`) to run `verso` from anywhere,
and source `scripts/autocomplete.sh` for tab-completion.

## Project layout

```
cmd/verso/       entry point, command dispatch
pkg/verso/       command handlers (init, add, commit, log, status, cat-file)
pkg/model/       core data model: Blob, Tree, Commit, Index, Database, Refs
utils/           SHA-1 hashing, zlib compression, file metadata (syscall-backed)
```

## Known limitations

This implements Git's core object model and staging flow, not the full command
surface. There's no branching, merging, diffing, or remote support (yet) —
contributions welcome.


