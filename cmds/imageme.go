package cmds

import (
	"fmt"
	"io"
)

func Imageme(p []string, w io.Writer) {
	fmt.Fprintf(w, "%v", p)
}
