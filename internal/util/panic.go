package util

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

func Panic(code int, format string, a ...any) {
	fmt.Fprint(os.Stderr, aurora.Red("panic:"))
	fmt.Fprintf(os.Stderr, " "+format+"\n", a...)

	os.Exit(2)
}
