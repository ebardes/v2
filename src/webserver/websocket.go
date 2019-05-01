package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"v2/config"
	"v2/personality"
	"v2/view"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WS is the Websocket entrypoint from the browser
func WS(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
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
			d, err := cfg.GetDisplay(m.Display)
			if err != nil {
				return err
			}

			di := view.DisplayInfo{}
			for i := range d.Layers {
				dl := view.DisplayLayer{}
				di.Layers[i] = dl
			}
			di.SetConnection(conn)
			m.Verb = personality.VerbAck

			di.Send(m)
		}
	}
}
