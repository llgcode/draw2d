// Ported from GLUT's samples.  Original copyright below applies.

/* Copyright (c) Mark J. Kilgard, 1996. */

/* This program is freely distributable without licensing fees 
   and is provided without guarantee or warrantee expressed or 
   implied. This program is -not- in the public domain. */

/* This program is a response to a question posed by Gil Colgate
   <gcolgate@sirius.com> about how lengthy a program is required using
   OpenGL compared to using  Direct3D immediate mode to "draw a
   triangle at screen coordinates 0,0, to 200,200 to 20,200, and I
   want it to be blue at the top vertex, red at the left vertex, and
   green at the right vertex".  I'm not sure how long the Direct3D
   program is; Gil has used Direct3D and his guess is "about 3000
   lines of code". */

package main

import (
	"gl"
	"glut"
	"exp/draw"
	"image"
	"freetype-go.googlecode.com/hg/freetype/raster"
	"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
	"postscript-go.googlecode.com/svn/trunk/postscript-go/src/pkg/postscript"
	"fmt"
)

type GLPainter struct {
	// The Porter-Duff composition operator.
	Op draw.Op
	// The 16-bit color to paint the spans.
	cr, cg, cb uint8
	ca         uint32
}

const M16 uint32 = 1<<16 - 1
const M32 uint32 = 1<<32 - 1

// Paint satisfies the Painter interface by painting ss onto an image.RGBA.
func (p *GLPainter) Paint(ss []raster.Span, done bool) {
	gl.Begin(gl.LINES)
	for _, s := range ss {
		ma := s.A >> 16
		a := ma * p.ca / M16
		gl.Color4ub(p.cr, p.cg, p.cb, uint8(a>>8))
		gl.Vertex2i(s.X0, s.Y)
		gl.Vertex2i(s.X1, s.Y)
	}
	gl.End()
}

// SetColor sets the color to paint the spans.
func (p *GLPainter) SetColor(c image.Color) {
	r, g, b, a := c.RGBA()
	if a == 0 {
		p.cr = 0
		p.cg = 0
		p.cb = 0
		p.ca = a
	} else {
		p.cr = uint8((r * M16 / a) >> 8)
		p.cg = uint8((g * M16 / a) >> 8)
		p.cb = uint8((b * M16 / a) >> 8)
		p.ca = a
	}
}

// NewRGBAPainter creates a new RGBAPainter for the given image.
func NewGLPainter() *GLPainter {
	return &GLPainter{}
}


func TestDrawCubicCurve(gc draw2d.GraphicContext) {
	// draw a cubic curve
	x, y := 25.6, 128.0
	x1, y1 := 102.4, 230.4
	x2, y2 := 153.6, 25.6
	x3, y3 := 230.4, 128.0

	gc.SetStrokeColor(image.NRGBAColor{0, 0, 0, 0xff})
	gc.SetLineWidth(10)
	gc.MoveTo(x, y)
	gc.CubicCurveTo(x1, y1, x2, y2, x3, y3)
	gc.Stroke()

	gc.SetStrokeColor(image.NRGBAColor{0xFF, 0x33, 0x33, 0x99})

	gc.SetLineWidth(6)
	// draw segment of curve
	gc.MoveTo(x, y)
	gc.LineTo(x1, y1)
	gc.MoveTo(x2, y2)
	gc.LineTo(x3, y3)
	gc.Stroke()
}

var (
	width, height int
)

func reshape(w, h int) {
	/* Because Gil specified "screen coordinates" (presumably with an
	   upper-left origin), this short bit of code sets up the coordinate
	   system to correspond to actual window coodrinates.  This code
	   wouldn't be required if you chose a (more typical in 3D) abstract
	   coordinate system. */
	gl.ClearColor(1, 1, 1, 1)
	//fmt.Println(gl.GetString(gl.EXTENSIONS))
	gl.Viewport(0, 0, w, h)                       /* Establish viewing area to cover entire window. */
	gl.MatrixMode(gl.PROJECTION)                  /* Start modifying the projection matrix. */
	gl.LoadIdentity()                             /* Reset project matrix. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1) /* Map abstract coords directly to window coords. */
	gl.Scalef(1, -1, 1)                           /* Invert Y axis so increasing Y goes down. */
	gl.Translatef(0, float32(-h), 0)              /* Shift origin up to upper-left corner. */
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	width, height = w, h
}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LineWidth(1)
	p := NewGLPainter()
	fmt.Println("draw")
	gc := draw2d.NewImageGraphicContextFromPainter(p, image.Rect(0, 0, width, height))
	gc.Translate(0, 380)
	gc.Scale(1, -1)
	gc.Translate(0, -380)
	interpreter := postscript.NewInterpreter(gc)
	interpreter.ExecuteFile("../../tiger.ps")
	gl.Flush() /* Single buffered, so needs a flush. */
}

func main() {
	glut.Init()
	glut.InitWindowSize(800, 800)
	glut.CreateWindow("single triangle")

	glut.DisplayFunc(display)
	glut.ReshapeFunc(reshape)
	glut.MainLoop()
}
