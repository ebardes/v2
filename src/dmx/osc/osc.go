package osc

import (
	"v2/config"
	"v2/dmx"

	xosc "github.com/hypebeast/go-osc/osc"
	"github.com/rs/zerolog/log"
)

// OSC implements a NetDMX Listener
type OSC struct {
	dmx.Common
	server *xosc.Server
}

// NewService creates a new instance
func NewService(c *config.Config) (*OSC, error) {
	var o OSC
	addr := "127.0.0.1:8765"
	o.server = &xosc.Server{Addr: addr}

	return &o, nil
}

// Run starts a listening thread
func (osc *OSC) Run() {
	log.Info().Msg("Started goroutine")
	defer log.Info().Msg("Exit goroutine")

	osc.server.Handle("/layer/address", func(msg *xosc.Message) {
		xosc.PrintMessage(msg)
	})
	osc.server.ListenAndServe()
}
