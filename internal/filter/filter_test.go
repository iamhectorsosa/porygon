package filter

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		include []string
		exclude []string
		file    string
		want    bool
	}{
		{"include match", []string{"go"}, nil, "a.go", true},
		{"case insensitive", []string{".GO"}, nil, "a.go", true},
		{"include miss", []string{"go"}, nil, "a.py", false},
		{"no include allows all", nil, nil, "a.py", true},
		{"blank include allows all", []string{" ", ""}, nil, "a.py", true},
		{"exclude wins", []string{"go"}, []string{"go"}, "a.go", false},
		{"exclude only", nil, []string{".md"}, "README.md", false},
		{"no extension, include set", []string{"go"}, nil, "Makefile", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := New(tc.include, tc.exclude).Match(tc.file); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
