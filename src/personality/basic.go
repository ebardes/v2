package personality

import "config"

type Personality interface {
	Decode([]byte) int
	Size() int
}

type BasicPersonality struct {
	Start  int
	level  byte
	page   byte
	slot   byte
	volume byte
}

func NewPersonality(c config.Pane) Personality {
	switch c.Personality {
	default:
	case "basic":
		bp := BasicPersonality{Start: c.StartAddress}
		return &bp
	}
	return nil
}

func (me *BasicPersonality) Decode(b []byte) int {
	if len(b)-me.Size() < me.Start {
		return 0
	}

	i := me.Start
	me.level = b[i]
	i++
	me.page = b[i]
	i++
	me.slot = b[i]
	i++
	me.volume = b[i]
	i++
	return i
}

func (*BasicPersonality) Size() int {
	return 4
}
