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
	// VerbRefresh requests a current packet be sent
	VerbRefresh Verb = "ref"
)

// Message is a shared object between the server and the JavaScript browser side.
type Message struct {
	Verb    Verb             `json:"verb"`
	Display uint             `json:"display"`
	Layer   uint             `json:"layer"`
	Layers  []int            `json:"layers,omitempty"`
	Packet  *RootPersonality `json:"packet,omitempty"`
}
