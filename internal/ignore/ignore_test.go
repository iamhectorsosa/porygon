package ignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMatch(t *testing.T) {
	m := Parse("# comment\n\n*.log\n!keep.log\nbuild/\n/root.txt\n**/gen/*.go\n")

	tests := []struct {
		rel   string
		isDir bool
		want  bool
	}{
		{"sub/a.log", false, true},
		{"keep.log", false, false}, // negated
		{"build", true, true},
		{"build", false, false}, // dir-only
		{"root.txt", false, true},
		{"sub/root.txt", false, false}, // anchored
		{"a/b/gen/x.go", false, true},  // **
		{"main.go", false, false},
	}
	for _, tc := range tests {
		if got := m.Match(tc.rel, tc.isDir); got != tc.want {
			t.Errorf("Match(%q, dir=%v) = %v, want %v", tc.rel, tc.isDir, got, tc.want)
		}
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()

	m, err := Load(dir) // no .gitignore
	if err != nil {
		t.Fatal(err)
	}
	if m.Match("a.tmp", false) {
		t.Error("missing .gitignore should ignore nothing")
	}

	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.tmp\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if m, err = Load(dir); err != nil {
		t.Fatal(err)
	}
	if !m.Match("a.tmp", false) {
		t.Error("a.tmp should be ignored")
	}
}
