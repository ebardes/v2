package webserver

import (
	"fmt"
	"log"
	"net/http"
	"v2/config"
)

// FN is a generic handler type
type FN func(http.ResponseWriter, *http.Request, *config.Config) error

var cfg *config.Config

// Register connections an HTTP path to a Function
func Register(pattern string, fn FN) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r, cfg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}
	})
}

// Run launches the webserver and runs "forever"
func Run(config *config.Config) {
	cfg = config
	http.Handle("/", http.FileServer(http.Dir(cfg.Static)))

	addr := fmt.Sprintf(":%d", cfg.WebPort)
	log.Println("Listening to " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
