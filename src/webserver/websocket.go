package webserver

import (
	"encoding/json"
	"net/http"
	"v2/config"
	"v2/personality"
	"v2/view"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WS is the Websocket entrypoint from the browser
func WS(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return err
	}
	defer conn.Close()

	log.Printf("Received websocket connection from %v\n", conn.RemoteAddr())
	for {
		x, p, err := conn.ReadMessage()
		if x == -1 {
			break // end of stream
		}

		if err != nil {
			log.Error().Err(err).Int("x", x).Msg("read error")
			return err
		}

		var m personality.Message
		err = json.Unmarshal(p, &m)
		if err != nil {
			log.Print(err)
			return err
		}

		switch m.Verb {
		case personality.VerbRegister:
			d, err := cfg.GetDisplay(m.Display)
			if err != nil {
				return err
			}

			di := view.FindDisplay(d)
			di.SetConnection(conn)

			layers := []uint{}
			for k := range di.Layers {
				layers = append(layers, k)
			}

			m.Verb = personality.VerbAck
			m.Layers = layers
			di.Send(m)
		}
	}

	view.RemoveConnection(conn)
	log.Print("Websocket shutdown")
	return nil
}
