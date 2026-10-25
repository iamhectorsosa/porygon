package mapper

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/iamhectorsosa/porygon/internal/filter"
	"github.com/iamhectorsosa/porygon/internal/ignore"
)

const maxFileSize = 256 << 10

type Options struct {
	Dir     string
	Include []string
	Exclude []string
}

type Entry struct {
	Path  string
	IsDir bool
}

type File struct {
	Path string
	Data []byte
}

type Result struct {
	Entries []Entry
	Files   []File
}

func Map(opts Options) (*Result, error) {
	root, err := filepath.Abs(opts.Dir)
	if err != nil {
		return nil, err
	}
	if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a readable directory", opts.Dir)
	}

	ig, err := ignore.Load(root)
	if err != nil {
		return nil, err
	}
	keep := filter.New(opts.Include, opts.Exclude)

	res := &Result{}
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == root {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		isDir := d.IsDir()
		if d.Name() == ".git" || ig.Match(rel, isDir) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}
		res.Entries = append(res.Entries, Entry{Path: rel, IsDir: isDir})

		if isDir || !keep.Match(d.Name()) {
			return nil
		}
		data, ok, err := readText(p, d)
		if err != nil {
			return err
		}
		if ok {
			res.Files = append(res.Files, File{Path: rel, Data: data})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func readText(p string, d fs.DirEntry) ([]byte, bool, error) {
	info, err := d.Info()
	if err != nil {
		return nil, false, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxFileSize {
		return nil, false, nil
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, false, err
	}
	head := data[:min(len(data), 8000)]
	if bytes.IndexByte(head, 0) >= 0 || !utf8.Valid(data) {
		return nil, false, nil
	}
	return data, true, nil
}
