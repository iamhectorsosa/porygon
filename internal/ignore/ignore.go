package ignore

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type Matcher struct{ m gitignore.Matcher }

func Load(dir string) (*Matcher, error) {
	data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return Parse(string(data)), nil
}

func Parse(content string) *Matcher {
	var ps []gitignore.Pattern
	for _, line := range strings.Split(content, "\n") {
		if line = strings.TrimSpace(line); line != "" && line[0] != '#' {
			ps = append(ps, gitignore.ParsePattern(line, nil))
		}
	}
	return &Matcher{gitignore.NewMatcher(ps)}
}

func (m *Matcher) Match(rel string, isDir bool) bool {
	return m.m.Match(strings.Split(rel, "/"), isDir)
}
