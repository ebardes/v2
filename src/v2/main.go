package main

import (
	"log"
	"v2/config"
	"v2/dmx"
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
	/*
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

		if cfg.Displays == nil || len(cfg.Displays) <= 0 {
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

		if cfg.Content == nil || len(cfg.Content) <= 0 {
			cfg.Content = map[int]content.Group{
				1: content.Group{
					Slots: map[int]content.Slot{
						1: content.Slot{
							Name: "Test Pattern",
							Type: "image",
							URL:  "static/TestBars.png",
						},
					},
				},
			}
		}

		cfg.Save()

		// Construct the runtime layers
		for _, display := range cfg.Displays {
			di := view.AddDisplay(display)
			for _, layer := range display.Layers {
				p := personality.NewPersonality(layer)
				dl := di.AddLayer(layer.StartAddress, p)
				DMX.AddLayer(layer.StartAddress, dl)
			}
		}

		go DMX.Run()
	*/
	webserver.Register("/index.go", view.Index)
	webserver.Register("/config.go", view.Config)
	webserver.Register("/display/", view.Display)
	webserver.Register("/ws/", webserver.WS)
	webserver.Run(cfg)
}
