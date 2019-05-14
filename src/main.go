package main

import (
	"fmt"
	"os"
	"v2/config"
	"v2/dmx"
	"v2/dmx/artnet"
	"v2/dmx/sacn"
	"v2/personality"
	"v2/view"
	"v2/webserver"

	flags "github.com/jessevdk/go-flags"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type options struct {
	TemplateDir string `short:"t" long:"template" description:"Location of template directory" default:"view"`
	Static      string `short:"s" long:"static" description:"Location of static content directory" default:"static"`
}

// DMX is the current DMX receiver
var DMX dmx.NetDMX
var opts options

func main() {
	tty := isatty.IsTerminal(os.Stderr.Fd())
	if tty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Error().Err(err).Msg("Parse error")
		os.Exit(1)
	}

	var cfg config.Config
	cfg.TemplateDir = opts.TemplateDir
	cfg.Static = opts.Static

	err = cfg.Load()
	if err != nil {
		log.Error().Err(err).Msg("Unable to load config file")
		return
	}

	// log.Println(cfg)
	switch cfg.Protocol {
	default:
		cfg.Protocol = "sacn"
		fallthrough

	case "sacn":
		DMX, err = sacn.NewService(&cfg)
		if err != nil {
			log.Error().Err(err).Msg("Error starting DMX listener")
		}

	case "artnet":
		DMX, err = artnet.NewService(&cfg)
		if err != nil {
			log.Error().Err(err).Msg("Error starting DMX listener")
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

	/*
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
	*/

	// Construct the runtime layers
	for _, display := range cfg.Displays {
		log.Info().Msg(fmt.Sprintf("%v", display))
		di := view.AddDisplay(&display)
		for _, layer := range display.Layers {
			p := personality.NewPersonality(layer.Personality)
			dl := di.AddLayer(layer.StartAddress, p)
			DMX.AddLayer(layer.StartAddress, dl)
		}
	}

	go DMX.Run()

	view.Init(&cfg)
	webserver.Register("/index.go", view.Index)
	webserver.Register("/config.go", view.Config)
	webserver.Register("/display/", view.Display)
	webserver.Register("/ws/", webserver.WS)
	webserver.Run(&cfg)
}
