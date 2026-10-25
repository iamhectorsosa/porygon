// Package render formats a mapper.Result as the XML-ish text the tool emits.
package render

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/iamhectorsosa/porygon/internal/mapper"
)

func Write(w io.Writer, res *mapper.Result) error {
	var b bytes.Buffer
	b.WriteString("<tree>\n")
	for _, e := range res.Entries {
		b.WriteString(treeLine(e))
	}
	b.WriteString("</tree>\n")

	for _, f := range res.Files {
		writeFile(&b, f)
	}
	_, err := b.WriteTo(w)
	return err
}

func treeLine(e mapper.Entry) string {
	line := strings.Repeat("  ", strings.Count(e.Path, "/")) + path.Base(e.Path)
	if e.IsDir {
		line += "/"
	}
	return line + "\n"
}

func writeFile(b *bytes.Buffer, f mapper.File) {
	fmt.Fprintf(b, "<file path=\"%s\">\n", f.Path)
	b.Write(f.Data)
	if len(f.Data) > 0 && f.Data[len(f.Data)-1] != '\n' {
		b.WriteByte('\n')
	}
	b.WriteString("</file>\n")
}
