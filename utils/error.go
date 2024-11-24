package utils

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

func VErrorAt(code string, pos scanner.Position, format string, arg ...interface{}) {
	start := pos.Offset - (pos.Column - 1)
	line := code[start : pos.Offset+1]
	fmt.Fprintln(os.Stderr, line)
	index := pos.Column - 1
	fmt.Fprintf(os.Stderr, "%s", strings.Repeat(" ", index))
	fmt.Fprintf(os.Stderr, "^")

	posDesc := fmt.Sprintf("(%d,%d):", pos.Line, pos.Column)
	fmt.Fprintf(os.Stderr, posDesc+format, arg...)
	fmt.Fprintf(os.Stderr, "\n")
}
