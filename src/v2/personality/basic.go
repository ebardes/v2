package personality

// Personality is the base interface for all personalities.
type Personality interface {
	Decode([]byte) int
}

// BasicPersonality is the core media server functionality.  It's a nice
// compact four channel personality.
type BasicPersonality struct {
	start  int
	level  byte
	page   byte
	slot   byte
	volume byte
}

// MediumPersonality extends the BasePersonality with scaling and pan functions.
type MediumPersonality struct {
	BasicPersonality
	x     int16
	y     int16
	xsize int16
	ysize int16
}

// NewPersonality instantiates the appropriate fixture definitions
func NewPersonality(personalityName string) (bp Personality) {
	switch personalityName {
	default:
	case "basic":
		bp = &BasicPersonality{}

	case "medium":
		bp = &MediumPersonality{}
	}

	return
}

// Decode reads the DMX packet and sets parameters accordingly
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

// Size returns the number of bytes required
func (*BasicPersonality) Size() int {
	return 4
}

func us(low byte, high byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}
func ss(low byte, high byte) int16 {
	return int16(uint16(high)<<8 | uint16(low))
}

// Decode the medium frame which extends from the Base Decode.
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

// Size returns the number of bytes required
func (me *MediumPersonality) Size() int {
	return me.BasicPersonality.Size() + 8
}
