package gui

import "testing"

func TestImage(t *testing.T) {

	g, err := GUIInit()
	if err != nil {
		t.Error(err)
		return
	}
	go insertImage(g)
	g.Run()
}
