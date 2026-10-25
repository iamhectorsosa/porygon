package render

import (
	"bytes"
	"testing"

	"github.com/iamhectorsosa/porygon/internal/mapper"
)

func TestWrite(t *testing.T) {
	res := &mapper.Result{
		Entries: []mapper.Entry{
			{Path: "a.go"},
			{Path: "sub", IsDir: true},
			{Path: "sub/deep", IsDir: true},
			{Path: "sub/deep/f.go"},
		},
		Files: []mapper.File{
			{Path: "a.go", Data: []byte("package a")}, // no trailing newline
			{Path: "sub/deep/f.go", Data: nil},        // empty file
		},
	}
	want := "<tree>\na.go\nsub/\n  deep/\n    f.go\n</tree>\n" +
		"<file path=\"a.go\">\npackage a\n</file>\n" +
		"<file path=\"sub/deep/f.go\">\n</file>\n"

	var b bytes.Buffer
	if err := Write(&b, res); err != nil {
		t.Fatal(err)
	}
	if got := b.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
