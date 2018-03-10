package gui

import (
	"fmt"

	"github.com/shibukawa/glfw"
	"github.com/shibukawa/gui4go"
)

type G struct {
	Nano  *nanogui.Screen
	Queue chan []byte
}

func init() {
	fmt.Println("Nanogui init")
	nanogui.Init()

}

func Open(title string, q chan []byte) (me *G) {
	glfw.WindowHint(glfw.Samples, 4)

	me = &G{
		Nano:  nanogui.NewScreen(1280, 720, title, true, false),
		Queue: q,
	}
	me.Nano.SetDrawContentsCallback(func() {})
	return
}

func (me *G) Run() {
	nanogui.MainLoop()
}

func (me *G) Close() {

}
