package gui

import (
	"fmt"
	"runtime"
	"v2/content"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/rs/zerolog/log"
)

type layer struct {
	data []byte
}

type G struct {
	win     *glfw.Window
	refresh bool
	layers  []layer
}

type IMG struct {
	height int
	width  int
	name   string
	handle uint32
}

func GUIInit() (g *G, err error) {
	runtime.LockOSThread()
	if err = glfw.Init(); err != nil {
		err = fmt.Errorf("could not initialize glfw: %v", err)
		return
	}

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

	win, err := glfw.CreateWindow(1920, 1080, "Display", nil, nil)
	if err != nil {
		err = fmt.Errorf("could not create opengl renderer: %v", err)
		return
	}

	win.MakeContextCurrent()
	if err = gl.Init(); err != nil {
		return
	}

	g = &G{win: win, refresh: true, layers: nil}
	gl.ClearColor(0, 0.1, 0.8, 1.0)
	return
}

// Run executes the main event loop
func (g *G) Run() {
	for !g.win.ShouldClose() {
		if g.refresh {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			g.win.SwapBuffers()
			g.refresh = false
		}
		glfw.WaitEventsTimeout(0.05)
	}
}

func (g *G) Close() {
	glfw.Terminate()
}

func (g *G) SetLayerImage(n int, data []byte) {
	for len(g.layers) <= n {
		g.layers = append(g.layers, layer{})
	}
	g.layers[n].data = data
}

func (g *G) DrawSlot(slot content.Slot) {
	filename := slot.GetName()
	img := getTexture(filename)

	xx := 0
	yy := 0
	ww := img.width
	hh := img.height
	angle := 0

	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, img.handle)

	gl.LoadIdentity()
	gl.Translatef(float32(xx), float32(yy), 0.0)
	gl.Rotatef(float32(angle), 0.0, 0.0, 1.0)
	gl.Translatef(-float32(xx), -float32(yy), 0.0)

	// Draw a textured quad
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(float32(xx), float32(yy))
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(float32(xx), float32(yy+hh))
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(float32(xx+ww), float32(yy+hh))
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(float32(xx+ww), float32(yy))

	gl.Disable(gl.TEXTURE_2D)
	gl.PopMatrix()

	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()

	gl.MatrixMode(gl.MODELVIEW)
	gl.End()

	// slot.Draw()
}

func getTexture(fn string) (img IMG, err error) {
	file, err := os.Open(fn)
	if err != nil {
		return
	}
	data, _, err := image.Decode(file)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(data.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), data, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	img.width = data.Bounds().Size().X
	img.height = data.Bounds().Size().Y
	img.handle = texture
	img.name = fn
	return
}
