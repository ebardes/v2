package sacn

import (
	"config"
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

const (
	srvAddrTemplate = "239.255.%d.%d:5568"
	maxDatagramSize = 8192
)

// SACN implements a NetDMX Listener
type SACN struct {
	socket *net.UDPConn
	cfg    *config.Config
}

// NewService creates a new instance
func NewService(c *config.Config) (*SACN, error) {
	x := SACN{cfg: c}
	univHigh := c.Universe >> 8
	univLow := c.Universe & 255
	network := fmt.Sprintf(srvAddrTemplate, univHigh, univLow)

	ifi, err := net.InterfaceByName(c.Interface)
	if err != nil {
		log.Println(err)
		ifi = nil
	}
	gaddr, err := net.ResolveUDPAddr("udp", network)
	if err != nil {
		log.Println(err)
	}
	socket, err := net.ListenMulticastUDP("udp", ifi, gaddr)
	if err != nil {
		log.Println(err)
	}

	x.socket = socket
	return &x, err
}

// Run starts a listening thread
func (x *SACN) Run() {
	log.Println("Started goroutine")
	defer log.Println("Exit goroutine")

	b := make([]byte, maxDatagramSize)
	for {
		n, err := x.socket.Read(b)
		if err != nil {
			log.Println(err)
			break
		}

		if n < 0x7d { // Too small
			continue
		}

		// ETC Visualization Mode filter
		if b[0x7d] > 0 {
			continue
		}
		if x.cfg.DebugLevel > 1 {
			fmt.Println(hex.Dump(b[0x7d:n]))
		}
	}
	x.socket.Close()
}

// Stop ends the running thread
func (x *SACN) Stop() {
	x.socket.Close()
}
