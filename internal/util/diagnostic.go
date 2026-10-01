package util

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)



type Diagnostic struct {
	Filename string
	Pos, Loc, Col int
	Line string
	Msg string
}

func (e Diagnostic) String() string {
	var sb strings.Builder

	gutterWidth := 4

	fmt.Fprintf(&sb, "%s at %s\n", aurora.Red("error"),
		aurora.Sprintf(aurora.Underline("%s:%d:%d"), e.Filename, e.Loc, e.Col))

	fmt.Fprintf(&sb, "%*s |\n", gutterWidth, "")
	fmt.Fprintf(&sb, "%*d | %s\n", gutterWidth, e.Loc, e.Line)
	fmt.Fprintf(&sb, "%*s | %*s\n", gutterWidth, "", e.Col, "^")
	fmt.Fprintf(&sb, "%*s = %s\n", gutterWidth, "", e.Msg)

	return sb.String()
}

func (e Diagnostic) Error() string {
	return e.String()
}

