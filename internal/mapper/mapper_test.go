package mapper

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func writeTestFile(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMap(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, ".git/HEAD", "ref: refs/heads/main\n")
	writeTestFile(t, root, ".gitignore", "ignored.go\nbuild/\n")
	writeTestFile(t, root, "a.go", "package a")
	writeTestFile(t, root, "ignored.go", "package x\n")
	writeTestFile(t, root, "build/x.go", "package x\n")
	writeTestFile(t, root, "sub/b.go", "package b\n")
	writeTestFile(t, root, "sub/c.txt", "text\n")

	res, err := Map(Options{Dir: root, Include: []string{"go"}})
	if err != nil {
		t.Fatal(err)
	}

	wantEntries := []Entry{
		{".gitignore", false},
		{"a.go", false},
		{"sub", true},
		{"sub/b.go", false},
		{"sub/c.txt", false},
	}
	if !slices.Equal(res.Entries, wantEntries) {
		t.Errorf("entries = %v, want %v", res.Entries, wantEntries)
	}

	if len(res.Files) != 2 ||
		res.Files[0].Path != "a.go" || string(res.Files[0].Data) != "package a" ||
		res.Files[1].Path != "sub/b.go" {
		t.Errorf("files = %+v", res.Files)
	}
}

func TestMapNotADirectory(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "f.txt", "x")
	for _, dir := range []string{filepath.Join(root, "missing"), filepath.Join(root, "f.txt")} {
		if _, err := Map(Options{Dir: dir}); err == nil {
			t.Errorf("Dir %q: expected error", dir)
		}
	}
}

func TestReadText(t *testing.T) {
	dir := t.TempDir()
	files := map[string][]byte{
		"text.txt":    []byte("hello\n"),
		"empty.txt":   {},
		"bin.dat":     []byte("ab\x00cd"),
		"invalid.txt": []byte("\xff\xfe"),
		"big.txt":     bytes.Repeat([]byte("a"), maxFileSize+1),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{
		"text.txt": true, "empty.txt": true,
		"bin.dat": false, "invalid.txt": false, "big.txt": false, "subdir": false,
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range entries {
		data, ok, err := readText(filepath.Join(dir, d.Name()), d)
		if err != nil {
			t.Fatalf("%s: %v", d.Name(), err)
		}
		if ok != want[d.Name()] {
			t.Errorf("%s: ok = %v, want %v", d.Name(), ok, want[d.Name()])
		}
		if ok && !bytes.Equal(data, files[d.Name()]) {
			t.Errorf("%s: data = %q", d.Name(), data)
		}
	}
}
