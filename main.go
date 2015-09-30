package hello

import (
	"appengine"
	"bytes"
	"cmds"
	"fmt"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/api/message", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	// token := r.Form.Get("token")
	text := r.Form.Get("text")
	fields := strings.Fields(text)

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "{\"text\":\"")

	if fields[1] == "payup" {
		cmds.Payup(appengine.NewContext(r), fields[2:], &buf)
	} else if fields[1] == "imageme" {
		cmds.Imageme(fields[2:], &buf)
	}

	fmt.Fprintf(&buf, "\"}")

	w.Write(buf.Bytes())
}
