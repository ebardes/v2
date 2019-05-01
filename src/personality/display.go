package personality

// Verb is shared between things
type Verb string

const (
	// VerbRegister Register a client
	VerbRegister Verb = "reg"
	// VerbAck Acknowledge the client registration
	VerbAck Verb = "ack"
	// VerbPacket Data from DMX
	VerbPacket Verb = "pkt"
)

// Message is a shared object between the server and the JavaScript browser side.
type Message struct {
	Verb    Verb   `json:"verb"`
	Display int    `json:"display"`
	Packet  []byte `json:"packet"`
	Start   int    `json:"start"`
	Lenght  int    `json:"end"`
	// Layer   uint   `json:"layer"`
	// Packet  Personality `json:"packet,optional"`
}
