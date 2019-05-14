package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"v2/config"
	"v2/personality"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type DisplayLayer struct {
	P     personality.Personality
	di    *DisplayInfo
	layer uint
}

type DisplayInfo struct {
	Display uint
	Layers  map[uint]DisplayLayer
	Conns   map[string]*websocket.Conn
}

var displays map[uint]*DisplayInfo

func init() {
	displays = make(map[uint]*DisplayInfo)
}

func Display(w http.ResponseWriter, r *http.Request, cfg *config.Config) error {
	var n uint
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
	ra := conn.RemoteAddr().String()
	if me.Conns == nil {
		me.Conns = make(map[string]*websocket.Conn)
	}
	me.Conns[ra] = conn
}

func (me *DisplayInfo) RemoveConnection(conn *websocket.Conn) {
	if me.Conns != nil {
		ra := conn.RemoteAddr().String()
		delete(me.Conns, ra)
	}
}

func (me *DisplayInfo) Send(msg personality.Message) error {
	if me.Conns == nil {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for ra, c := range me.Conns {
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Error().Err(err).Str("Remote Address", ra).Msg("Send error")
		}
	}
	return nil
}

func (me *DisplayInfo) AddLayer(start uint, p personality.Personality) *DisplayLayer {
	dl := DisplayLayer{P: p, di: me, layer: start}
	me.Layers[start] = dl
	return &dl
}

func (me *DisplayLayer) OnFrame(b []byte) {
	blob := personality.Blob{
		Data: b,
		Pos:  0,
	}
	me.P.Decode(&blob)
	msg := me.P.DataMessage()
	msg.Layer = me.layer
	msg.Display = me.di.Display
	me.di.Send(msg)
}

func FindDisplay(cd *config.Display) (di *DisplayInfo) {
	var ok bool

	n := cd.ID
	di, ok = displays[n]
	if !ok {
		di = &DisplayInfo{}
		di.Layers = make(map[uint]DisplayLayer)
		di.Display = n
		displays[n] = di
	}
	return
}

func RemoveConnection(conn *websocket.Conn) {
	for _, di := range displays {
		di.RemoveConnection(conn)
	}
}
