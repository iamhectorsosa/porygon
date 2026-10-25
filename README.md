# Porygon 🔺

Porygon is a dead simple CLI tool that maps a codebase into a single, LLM-friendly text output: a directory tree followed by the contents of your source files.

Ideal for pasting a project into an AI chat, or piping it into other tools.

## Quick Overview

- **Tree + Contents**: A `<tree>` of your project, followed by each file wrapped in a `<file path="...">` tag.
- **Respects `.gitignore`**: Ignored files and the `.git` directory are skipped automatically.
- **Text Only**: Binary files, invalid UTF-8, and files over 256 KB are left out.
- **Filter by Extension**: Include only the file types you care about.
- **Pipe-friendly**: Output goes to `stdout`.

## Requirements

- Go 1.21 or later

## Installation

```
go install github.com/iamhectorsosa/porygon@latest
```

Or from source:

```
git clone https://github.com/iamhectorsosa/porygon.git
cd porygon
go build -o porygon
```

## Usage

### Basic Usage

```
porygon
```

Maps the current directory.

### Examples

```
# Map a specific directory
porygon --dir ~/projects/my-app

# Include only Go files
porygon --include go

# Include multiple file types (leading dot is optional)
porygon -i .go -i .md

# Copy the result to your clipboard (macOS)
porygon -d ./src -i ts,tsx | pbcopy

# Save to a file
porygon -i go > codebase.txt
```

### Output

```
<tree>
go.mod
main.go
internal/
  filter/
    filter.go
</tree>
<file path="main.go">
package main
...
</file>
```

The tree lists every non-ignored file and directory. File contents are included only for files matching `--include` (or all text files if no filter is set).

### Flags

- `--dir`, `-d`: Directory to map (default: `.`)
- `--include`, `-i`: File extensions to include. Repeatable or comma-separated.

## Development

```
go mod download
go test ./...
go build
```
# porygon
