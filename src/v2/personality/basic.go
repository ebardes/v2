package personality

import (
	"v2/config"
)

type Personality interface {
	Decode([]byte) int
}

type BasicPersonality struct {
	start  int
	level  byte
	page   byte
	slot   byte
	volume byte
}

type MediumPersonality struct {
	BasicPersonality
	x     int16
	y     int16
	xsize int16
	ysize int16
}

func NewPersonality(c config.Layer) (bp Personality) {
	switch c.Personality {
	default:
	case "basic":
		bp = &BasicPersonality{}

	case "medium":
		bp = &MediumPersonality{}
	}

	return
}

func (me *BasicPersonality) Decode(b []byte) (i int) {
	if len(b)-me.Size() < me.start {
		return 0
	}

	i = me.start
	me.level = b[i]
	i++
	me.page = b[i]
	i++
	me.slot = b[i]
	i++
	me.volume = b[i]
	i++
	return
}

func (*BasicPersonality) Size() int {
	return 4
}

func us(low byte, high byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}
func ss(low byte, high byte) int16 {
	return int16(uint16(high)<<8 | uint16(low))
}

func (me *MediumPersonality) Decode(b []byte) (i int) {
	i = me.BasicPersonality.Decode(b)
	me.x = ss(b[i], b[i+1])
	i += 2
	me.y = ss(b[i], b[i+1])
	i += 2
	me.xsize = ss(b[i], b[i+1])
	i += 2
	me.ysize = ss(b[i], b[i+1])
	i += 2
	return
}

func (me *MediumPersonality) Size() int {
	return me.BasicPersonality.Size() + 8
}
