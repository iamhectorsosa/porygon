package filter

import (
	"path/filepath"
	"strings"
)

type Filter struct {
	include, exclude map[string]bool
}

func New(include, exclude []string) *Filter {
	return &Filter{include: extSet(include), exclude: extSet(exclude)}
}

func (f *Filter) Match(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if f.exclude[ext] {
		return false
	}
	return len(f.include) == 0 || f.include[ext]
}

func extSet(xs []string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		if x = strings.ToLower(strings.TrimSpace(x)); x != "" {
			m["."+strings.TrimPrefix(x, ".")] = true
		}
	}
	return m
}
