package view

import (
	"fmt"
	"net/http"
	"strings"
)

type DisplayInfo struct {
	Panel uint
}

func Display(w http.ResponseWriter, r *http.Request) error {
	p := DisplayInfo{}
	path := r.URL.Path
	i := strings.LastIndex(path, "/")
	if i > 0 {
		path = path[i+1:]
	}
	fmt.Sscanf(path, "%d", &p.Panel)

	v := InitView(w).Prefix("../")
	return v.Render("display.tpl", p)
}
