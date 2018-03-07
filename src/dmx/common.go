package dmx

import (
	"config"
	"encoding/hex"
	"fmt"
	"personality"
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
func (me *Common) OnFrame(b []byte) {
	if changed(b, me.Frame) {
		me.Frame = b

		if me.Cfg.DebugLevel > 1 {
			fmt.Println(hex.Dump(b))
		}

		if me.Bindings != nil {
			for _, pers := range me.Bindings {
				pers.Decode(b)
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
