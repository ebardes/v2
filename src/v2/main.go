package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"v2/config"
	"v2/dmx"
	"v2/dmx/artnet"
	"v2/dmx/sacn"
	"v2/personality"
)

var cfg config.Config

// DMX is the current DMX receiver
var DMX dmx.NetDMX

var done chan bool

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
	switch cfg.Protocol {
	default:
		cfg.Protocol = "sacn"
		fallthrough

	case "sacn":
		DMX, err = sacn.NewService(&cfg)
		if err != nil {
			log.Fatal(err)
		}

	case "artnet":
		DMX, err = artnet.NewService(&cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	cfg.Save()

	for _, pane := range cfg.Panes {
		p := personality.NewPersonality(pane)
		DMX.AddPersonality(&p)
	}

	go DMX.Run()

	// go func() {
	http.HandleFunc("/app/", handler)
	http.HandleFunc("/view/config.html", handleConfig)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Printf("Listening... port %d", cfg.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.WebPort), nil)
	// }()
	//
	// g := gui.Open("V2", make(chan []byte))
	// g.Run()

	<-done
}
