package webserver

import (
	"fmt"
	"log"
	"net/http"
	"v2/config"
)

type FN func(http.ResponseWriter, *http.Request) error

func Register(pattern string, fn FN) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}
	})
}

func Run(cfg config.Config) {
	addr := fmt.Sprintf(":%d", cfg.WebPort)
	log.Println("Listening to " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
