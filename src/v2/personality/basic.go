package personality

import (
	"bytes"
	"encoding/binary"
	"v2/config"
)

type Personality interface {
	Decode([]byte)
}

type BasicPersonality struct {
	// start  int
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

func (me *BasicPersonality) Decode(b []byte) {
	binary.Read(bytes.NewReader(b), binary.LittleEndian, me)
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

func (me *MediumPersonality) Decode(b []byte) {
	binary.Read(bytes.NewReader(b), binary.LittleEndian, me)
}

func (me *MediumPersonality) Size() int {
	return me.BasicPersonality.Size() + 8
}
