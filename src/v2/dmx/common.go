package dmx

import (
	"encoding/hex"
	"fmt"
	"v2/config"
	"v2/personality"
)

// NetDMX is the base interface for DMX
type NetDMX interface {
	Run()
	Stop()
	AddPersonality(p *personality.Personality)
}

// Common stores aspects common to all Networked DMX implementations
type Common struct {
	NetDMX
	Frame    []byte
	Cfg      *config.Config
	Bindings []personality.Personality
}

func changed(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return true
	}
	for i, v := range a {
		if v != b[i] {
			return true
		}
	}
	return false
}

// OnFrame is the main event listener for when DMX packets arrive
func (me *Common) OnFrame(universe int, b []byte) {
	if changed(b, me.Frame) {
		if len(me.Frame) != len(b) {
			me.Frame = make([]byte, len(b))
		}
		copy(me.Frame, b)

		if me.Cfg != nil && me.Cfg.DebugLevel > 1 {
			fmt.Printf("Universe: %d\n", universe)
			fmt.Print(hex.Dump(b))
		}

		if me.Bindings != nil {
			for _, pers := range me.Bindings {
				if pers != nil {
					pers.Decode(b)
				}
			}
		}
	}
}

// AddPersonality binds a display to the listener
func (me *Common) AddPersonality(p *personality.Personality) {
	if me.Bindings == nil {
		me.Bindings = make([]personality.Personality, 0)
	}
	me.Bindings = append(me.Bindings, *p)
}
