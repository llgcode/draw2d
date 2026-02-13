// Open an OpenGL window and display graphics using the modern OpenGL ES 2 backend
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgles2"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	width, height  = 800, 600
	rotate         int
	redraw         = true
	wantScreenshot = false
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

	// Request OpenGL 3.2 core profile for modern shader support.
	// The gles2 backend uses VAO+VBO+EBO so it works in core profile.
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4) // Enable 4x MSAA for antialiasing

	window, err := glfw.CreateWindow(width, height, "draw2d OpenGL ES 2 Example", nil, nil)
	if err != nil {
		// Fall back to compatibility profile if core profile fails
		glfw.DefaultWindowHints()
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
			if wantScreenshot {
				saveScreenshot("output/samples/helloworldgles2.png")
				wantScreenshot = false
			}
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
	case key == glfw.KeyS && action == glfw.Press:
		wantScreenshot = true
		redraw = true
	}
}

// saveScreenshot reads back the framebuffer and saves it as PNG
func saveScreenshot(filename string) {
	pixels := make([]uint8, width*height*4)
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// OpenGL reads from bottom-up, flip rows
	for y := 0; y < height; y++ {
		glY := height - 1 - y
		copy(img.Pix[y*width*4:(y+1)*width*4], pixels[glY*width*4:(glY+1)*width*4])
	}

	if err := os.MkdirAll("output/samples", 0755); err != nil {
		log.Printf("Failed to create output dir: %v", err)
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to save screenshot: %v", err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Printf("Failed to encode PNG: %v", err)
		return
	}
	log.Printf("Screenshot saved to %s", filename)
}
