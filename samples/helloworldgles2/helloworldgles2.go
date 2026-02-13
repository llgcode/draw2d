// Open an OpenGL window and display graphics using the modern OpenGL ES 2 backend
package main

import (
	"image/color"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgles2"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	width, height = 800, 600
	rotate        int
	redraw        = true
)

func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Viewport(0, 0, int32(w), int32(h))
	
	// Enable blending for transparency
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	
	width, height = w, h
	redraw = true
}

func display(gc *draw2dgles2.GraphicContext) {
	// Clear screen
	gc.Clear()

	// Draw filled rectangle
	gc.SetFillColor(color.RGBA{200, 50, 50, 255})
	gc.BeginPath()
	draw2dkit.Rectangle(gc, 50, 50, 250, 250)
	gc.Fill()

	// Draw stroked rounded rectangle
	gc.SetStrokeColor(color.RGBA{50, 50, 200, 255})
	gc.SetLineWidth(5)
	gc.BeginPath()
	draw2dkit.RoundedRectangle(gc, 300, 50, 500, 250, 20, 20)
	gc.Stroke()

	// Draw filled circle
	gc.SetFillColor(color.RGBA{50, 200, 50, 255})
	gc.BeginPath()
	draw2dkit.Circle(gc, 400, 400, 80)
	gc.Fill()

	// Draw filled and stroked ellipse
	gc.SetFillColor(color.RGBA{200, 200, 50, 200})
	gc.SetStrokeColor(color.RGBA{100, 100, 100, 255})
	gc.SetLineWidth(3)
	gc.BeginPath()
	draw2dkit.Ellipse(gc, 150, 450, 100, 60)
	gc.FillStroke()

	// Flush all batched drawing commands to GPU
	gc.Flush()

	gl.Flush()
}

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// Request OpenGL 3.2 core profile (minimum for modern shaders)
	// Note: Can also use OpenGL ES 2.0+ on mobile/embedded
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "draw2d OpenGL ES 2 Example", nil, nil)
	if err != nil {
		// Fall back to default context if core profile fails
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		window, err = glfw.CreateWindow(width, height, "draw2d OpenGL ES 2 Example", nil, nil)
		if err != nil {
			panic(err)
		}
	}

	window.MakeContextCurrent()
	window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	log.Printf("OpenGL version: %s", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Printf("GLSL version: %s", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	// Create graphics context
	gc, err := draw2dgles2.NewGraphicContext(width, height)
	if err != nil {
		panic(err)
	}
	defer gc.Destroy()

	// Setup font
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic,
	})

	reshape(window, width, height)

	for !window.ShouldClose() {
		if redraw {
			display(gc)
			window.SwapBuffers()
			redraw = false
		}
		glfw.PollEvents()
	}
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	case key == glfw.KeySpace && action == glfw.Press:
		redraw = true
	}
}
