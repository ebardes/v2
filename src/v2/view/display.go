package view

import (
	"fmt"
	"net/http"
	"strings"
	"v2/config"
	"v2/personality"
)

type DisplayLayer struct {
	P personality.Personality
}

type DisplayInfo struct {
	Panel  uint
	Layers map[int]DisplayLayer
}

var displays map[uint]DisplayInfo

func init() {
	displays = make(map[uint]DisplayInfo)
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

func AddDisplay(pane config.Display) *DisplayInfo {
	d := DisplayInfo{
		Panel:  pane.ID,
		Layers: make(map[int]DisplayLayer),
	}
	displays[pane.ID] = d
	return &d
}

func (me *DisplayInfo) AddLayer(p personality.Personality) {
	me.Layers[p.Start()] = DisplayLayer{P: p}
}
