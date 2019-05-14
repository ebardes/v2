package personality

// Personality is the base interface for all personalities.
type Personality interface {
	Decode(*Blob)
	DataMessage() Message
}

type RootPersonality struct {
	display uint

	Level      byte
	Group      byte
	Slot       byte
	Volume     byte
	X          int16
	Y          int16
	Xsize      int16
	Ysize      int16
	ZRotate    int16
	Brightness int8
	Contrast   int8
	PlayMode   byte
	Changed    bool
	URL        string
}

// BasicPersonality is the core media server functionality.  It's a nice
// compact four channel personality.
type BasicPersonality struct {
	RootPersonality
}

// MediumPersonality extends the BasePersonality with scaling and pan functions.
type MediumPersonality struct {
	BasicPersonality
}

type ExtendedPersonality struct {
	MediumPersonality
}

// NewPersonality instantiates the appropriate fixture definitions
func NewPersonality(personalityName string) (bp Personality) {
	switch personalityName {
	default:
	case "basic":
		bp = &BasicPersonality{}

	case "medium":
		bp = &MediumPersonality{}

	case "extended":
		bp = &ExtendedPersonality{}
	}

	return
}

// Decode reads the DMX packet and sets parameters accordingly
func (me *BasicPersonality) Decode(b *Blob) {
	if !b.Check(4) {
		return
	}

	me.Level = b.Byte()
	group := b.Byte()
	slot := b.Byte()
	if group != me.Group || slot != me.Slot {
		me.Group = group
		me.Slot = slot
		me.Changed = true
	} else {
		me.Changed = false
	}
	me.Volume = b.Byte()
	return
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

func (me *ExtendedPersonality) Decode(b *Blob) {
	me.MediumPersonality.Decode(b)
	return
}

func (me *RootPersonality) DataMessage() Message {
	return Message{
		Verb:   VerbPacket,
		Packet: me,
	}
}
