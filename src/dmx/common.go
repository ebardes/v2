package dmx

import (
	"bytes"
	"net"
	"sync"
	"time"
	"v2/config"
	"v2/view"

	"github.com/rs/zerolog/log"
)

// NetDMX is the base interface for DMX
type NetDMX interface {
	Run()
	Stop()
	Refresh()
	AddLayer(start uint, layer *view.DisplayLayer)
}

// Common stores aspects common to all Networked DMX implementations
type Common struct {
	NetDMX
	Frame  []byte
	Cfg    *config.Config
	Layers []DMX2Layer
	Sync   sync.Mutex
	last   time.Time
}

type DMX2Layer struct {
	start uint
	dl    *view.DisplayLayer
}

var universes map[uint][]NetDMX

func init() {
	universes = make(map[uint][]NetDMX)
}

// OnFrame is the main event listener for when DMX packets arrive
func (me *Common) OnFrame(addr net.Addr, b []byte) {
	me.Sync.Lock()
	defer me.Sync.Unlock()

	now := time.Now()

	if !bytes.Equal(b, me.Frame) || now.Sub(me.last) > (10*time.Second) {
		me.last = now
		if len(me.Frame) != len(b) {
			me.Frame = make([]byte, len(b))
		}
		copy(me.Frame, b)
		if me.Cfg.DebugLevel > 4 {
			log.Printf("Packet from %v size %d\n", addr, len(b))
		}

		/*
			if me.Cfg != nil && me.Cfg.DebugLevel > 1 {
				fmt.Printf("Universe: %d\n", universe)
				fmt.Print(hex.Dump(b))
			}
		*/

		if me.Layers != nil {
			for _, layer := range me.Layers {
				if layer.dl != nil {
					layer.dl.OnFrame(b[layer.start:])
				}
			}
		}
	}
}

func (me *Common) Refresh() {
	me.Frame = []byte{}
}

func (me *Common) AddLayer(start uint, dl *view.DisplayLayer) {
	if me.Layers == nil {
		me.Layers = []DMX2Layer{}
	}
	me.Layers = append(me.Layers, DMX2Layer{start - 1, dl})
}

func Refresh() {
	for _, u := range universes {
		for _, d := range u {
			d.Refresh()
		}
	}
}
