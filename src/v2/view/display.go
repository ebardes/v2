package view

import (
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

var displays map[uint]*DisplayInfo

func init() {
	displays = make(map[uint]*DisplayInfo)
}

func Display(w http.ResponseWriter, r *http.Request) error {
	var n uint
	path := r.URL.Path
	i := strings.LastIndex(path, "/")
	if i > 0 {
		path = path[i+1:]
	}
	fmt.Sscanf(path, "%d", &n)
	d, ok := displays[n]
	if !ok {
		return fmt.Errorf("Display %v not found", n)
	}

	v := InitView(w).Prefix("../")
	return v.Render("display.tpl", d)
}

func AddDisplay(display config.Display) *DisplayInfo {
	d := DisplayInfo{
		Panel:  display.ID,
		Layers: make(map[int]DisplayLayer),
	}
	displays[display.ID] = &d
	return &d
}

func GetDisplay(id uint) *DisplayInfo {

	d, ok := displays[id]
	if ok {
		return d
	}
	return nil
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

func (me *DisplayInfo) AddLayer(start int, p personality.Personality) {
	me.Layers[start] = DisplayLayer{P: p, Chan: make(chan bool, 25)}
}
