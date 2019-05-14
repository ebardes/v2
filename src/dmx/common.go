package dmx

import (
	"bytes"
	"log"
	"net"
	"sync"
	"v2/config"
	"v2/view"
)

// NetDMX is the base interface for DMX
type NetDMX interface {
	Run()
	Stop()
	AddLayer(start int, layer *view.DisplayLayer)
}

// Common stores aspects common to all Networked DMX implementations
type Common struct {
	NetDMX
	Frame  []byte
	Cfg    *config.Config
	Layers []DMX2Layer
	Sync   sync.Mutex
}

type DMX2Layer struct {
	start int
	dl    *view.DisplayLayer
}

// OnFrame is the main event listener for when DMX packets arrive
func (me *Common) OnFrame(addr net.Addr, universe int, b []byte) {
	me.Sync.Lock()
	defer me.Sync.Unlock()

	if !bytes.Equal(b, me.Frame) {
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
				layer.dl.OnFrame(b[layer.start:])
			}
		}
	}
}

func (me *Common) AddLayer(start int, dl *view.DisplayLayer) {
	if me.Layers == nil {
		me.Layers = []DMX2Layer{}
	}
	me.Layers = append(me.Layers, DMX2Layer{start - 1, dl})
}
