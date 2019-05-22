package view

import (
	"net/http"
	"v2/config"
)

func Index(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
	v := InitView(w).Title("Index")
	return v.Render("index.tpl", nil)
}

func Config(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
	cfg.Protocols = []string{
		config.ProtocolSACN,
		config.ProtocolArtNet,
	}
	v := InitView(w).Title("Config")
	return v.Render("config.tpl", cfg)
}
