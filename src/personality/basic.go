package personality

// Personality is the base interface for all personalities.
type Personality interface {
	Decode(*Blob)
}

// BasicPersonality is the core media server functionality.  It's a nice
// compact four channel personality.
type BasicPersonality struct {
	start  int
	Level  byte
	Group  byte
	Slot   byte
	Volume byte
}

// MediumPersonality extends the BasePersonality with scaling and pan functions.
type MediumPersonality struct {
	BasicPersonality
	X          int16
	Y          int16
	Xsize      int16
	Ysize      int16
	ZRotate    int16
	Brightness int8
	Contrast   int8
	PlayMode   byte
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
func (me *BasicPersonality) Decode(b *Blob) {
	if !b.Check(4) {
		return
	}

	me.Level = b.Byte()
	me.Group = b.Byte()
	me.Slot = b.Byte()
	me.Volume = b.Byte()
	return
}

// Size returns the number of bytes required
func (*BasicPersonality) Size() int {
	return 4
}

// Decode the medium frame which extends from the Base Decode.
func (me *MediumPersonality) Decode(b *Blob) {
	me.BasicPersonality.Decode(b)
	if !b.Check(13) {
		return
	}
	me.X = b.SignedShort()
	me.Y = b.SignedShort()
	me.Xsize = b.SignedShort()
	me.Ysize = b.SignedShort()
	me.ZRotate = b.SignedShort()
	me.Brightness = b.SignedByte()
	me.Contrast = b.SignedByte()
	me.PlayMode = b.Byte()
	return
}

// Size returns the number of bytes required
func (me *MediumPersonality) Size() int {
	return me.BasicPersonality.Size() + 8
}
