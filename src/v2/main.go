package main

import (
	"config"
	"fmt"
	"log"
	"net/artnet"
	"net/http"
	"net/sacn"
	"text/template"
)

var cfg config.Config

type NetDMX interface {
	Run()
	Stop()
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("view/config.html"))
	x := t.Execute(w, cfg) // merge.
	if x != nil {
		log.Println(x)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	err := config.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(cfg)
	var x NetDMX

	switch cfg.Protocol {
	default:
		cfg.Protocol = "sacn"
		fallthrough

	case "sacn":
		x, err = sacn.NewService(&cfg)
		if err != nil {
			log.Fatal(err)
		}

	case "artnet":
		x, err = artnet.NewService(&cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	cfg.Save()

	go x.Run()

	http.HandleFunc("/app/", handler)
	http.HandleFunc("/view/config.html", handleConfig)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Printf("Listening... port %d", cfg.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebPort), nil)
}
