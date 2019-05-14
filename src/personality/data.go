package personality

type Blob struct {
	Data []byte
	Pos  int
}

func (b *Blob) SignedShort() int16 {
	high := b.Data[b.Pos]
	b.Pos++
	low := b.Data[b.Pos]
	b.Pos++
	return ss(low, high)
}

func (b *Blob) UnsignedShort() uint16 {
	high := b.Data[b.Pos]
	b.Pos++
	low := b.Data[b.Pos]
	b.Pos++
	return us(low, high)
}

func (b *Blob) Check(n int) bool {
	return len(b.Data)-b.Pos > n
}

func (b *Blob) Byte() byte {
	x := b.Data[b.Pos]
	b.Pos++
	return x
}

func (b *Blob) SignedByte() int8 {
	return int8(int(b.Byte()) + 128)
}

func us(low byte, high byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}
func ss(low byte, high byte) int16 {
	return int16(int(uint16(high)<<8|uint16(low)) + 32768)
}
