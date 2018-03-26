package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"v2/personality"
	"v2/view"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WS is the Websocket entrypoint from the browser
func WS(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	log.Printf("Received connection from %v\n", conn.RemoteAddr())
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}

		var m personality.Message
		err = json.Unmarshal(p, &m)
		if err != nil {
			log.Println(err)
			return err
		}

		switch m.Verb {
		case personality.VerbRegister:
			di := view.GetDisplay(m.Display)
			di.SetConnection(conn)
			m.Verb = personality.VerbAck

			di.Send(m)
		}
	}
}
