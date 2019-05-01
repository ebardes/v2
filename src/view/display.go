package view

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"v2/config"
	"v2/personality"

	"github.com/gorilla/websocket"
)

type DisplayLayer struct {
	P    personality.Personality
	Chan chan bool
}

type DisplayInfo struct {
	Panel  uint
	Layers map[int]DisplayLayer
	Conn   *websocket.Conn
}

func Display(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
	var n int
	path := r.URL.Path
	i := strings.LastIndex(path, "/")
	if i > 0 {
		path = path[i+1:]
	}
	fmt.Sscanf(path, "%d", &n)
	d, err := cfg.GetDisplay(n)
	if err != nil {
		return err
	}
	v := InitView(w).Prefix("../")
	return v.Render("display.tpl", d)
}

func (me *DisplayInfo) SetConnection(conn *websocket.Conn) {
	me.Conn = conn
}

func (me *DisplayInfo) Send(msg personality.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return me.Conn.WriteMessage(websocket.TextMessage, data)
}

func (me *DisplayInfo) AddLayer(start int, p personality.Personality) *DisplayLayer {
	dl := DisplayLayer{P: p, Chan: make(chan bool, 25)}
	me.Layers[start] = dl
	return &dl
}

func (me *DisplayLayer) OnFrame(b []byte) {
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &me.P)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(json.Marshal(me.P))
}
