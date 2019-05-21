package gui

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/rs/zerolog/log"
)

type G struct {
	win     *glfw.Window
	refresh bool
}

func GUIInit() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	defer glfw.Terminate()

	monitors := glfw.GetMonitors()
	for _, m := range monitors {
		vid := m.GetVideoMode()
		log.Info().Str("name", m.GetName()).Int("x", vid.Width).Int("y", vid.Height).Msg("monitor")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(1920, 1080, "Hello world", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(0, 0.1, 0.8, 1.0)
	g := G{win: win, refresh: true}
	g.run()
}

func (g *G) run() {
	for !g.win.ShouldClose() {
		if g.refresh {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			g.win.SwapBuffers()
			g.refresh = false
		}
		glfw.WaitEventsTimeout(0.05)
	}
}
