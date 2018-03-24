package personality

type Verb string

const (
	VerbRegister Verb = "reg"
	VerbAck      Verb = "ack"
	VerbPacket   Verb = "pkt"
)

type Message struct {
	Verb    Verb        `json:"verb"`
	Display uint        `json:"display"`
	Layer   uint        `json:"layer"`
	Packet  Personality `json:"packet,optional"`
}
