package main

import (
	"log"
	"v2/config"
	"v2/dmx"
	"v2/dmx/artnet"
	"v2/dmx/sacn"
	"v2/personality"
	"v2/view"
	"v2/web"
)

// DMX is the current DMX receiver
var DMX dmx.NetDMX

var done chan bool

func main() {
	var cfg config.Config

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

	if len(cfg.Displays) <= 0 {
		p := config.Display{
			ID: 1,
			Layers: []config.Layer{
				config.Layer{
					Personality:  "basic",
					StartAddress: 1,
				},
			},
		}
		cfg.Displays = append(cfg.Displays, p)
	}

	cfg.Save()

	for _, display := range cfg.Displays {
		d := view.AddDisplay(display)
		for _, layer := range display.Layers {
			p := personality.NewPersonality(layer)
			DMX.AddPersonality(&p)
			d.AddLayer(p)
		}
	}

	go DMX.Run()

	webserver.Register("/index.html", view.Index)
	webserver.Register("/config.html", view.Config)
	webserver.Register("/display/", view.Display)
	webserver.Run(cfg)
}
