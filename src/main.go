package main

import (
	"fmt"
	"os"
	"v2/config"
	"v2/dmx"
	"v2/dmx/artnet"
	"v2/dmx/osc"
	"v2/dmx/sacn"
	"v2/personality"
	"v2/view"
	"v2/webserver"

	_ "image/jpeg"
	_ "image/png"

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

	cfg := &config.GlobalConfig
	cfg.TemplateDir = opts.TemplateDir
	cfg.StaticDir = opts.Static

	err = cfg.Load()
	if err != nil {
		log.Error().Err(err).Msg("Unable to load config file")
		return
	}
	cfg.Normalize()

	switch cfg.Protocol {
	default:
		cfg.Protocol = "sacn"
		fallthrough

	case "sacn":
		DMX, err = sacn.NewService(cfg)
		if err != nil {
			log.Error().Err(err).Msg("Error starting DMX listener")
		}

	case "artnet":
		DMX, err = artnet.NewService(cfg)
		if err != nil {
			log.Error().Err(err).Msg("Error starting DMX listener")
		}

	case "osc":
		DMX, err = osc.NewService(cfg)
	}

	if cfg.Displays == nil || len(cfg.Displays) <= 0 {
		p := config.Display{
			ID: 1,
			Layers: []config.Layer{
				config.Layer{
					Personality:  "regular",
					StartAddress: 1,
				},
			},
		}
		cfg.Displays = append(cfg.Displays, p)
	}

	// Construct the runtime layers
	for _, display := range cfg.Displays {
		log.Info().Msg(fmt.Sprintf("%v", display))
		di := view.FindDisplay(&display)
		for _, layer := range display.Layers {
			p := personality.NewPersonality(layer.Personality)
			dl := di.AddLayer(layer.StartAddress, p)
			DMX.AddLayer(layer.StartAddress, dl)
		}
	}

	go DMX.Run()
	ws(cfg)
}

func ws(cfg *config.Config) {
	view.Init(cfg)
	webserver.Register("/index.go", view.Index)
	webserver.Register("/config.go", view.Config)
	webserver.Register("/display/", view.Display)
	webserver.Register("/ws/", webserver.WS)
	webserver.Run(cfg)
}
