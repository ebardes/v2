package view

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) error {
	v := InitView(w).Title("Index")
	return v.Render("index.tpl", nil)
}

func Config(w http.ResponseWriter, r *http.Request) error {
	v := InitView(w).Title("Config")
	return v.Render("config.tpl", nil)
}
