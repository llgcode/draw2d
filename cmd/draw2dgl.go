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
	"os"
	"math"
	"io/ioutil"
	"strings"
	"gl"
	"glut"
	"image"
	"draw2d.googlecode.com/hg/draw2d"
	"draw2d.googlecode.com/hg/draw2dgl"
	"draw2d.googlecode.com/hg/postscript"
	"log"
	"time"
)

var postscriptContent string

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
	rotate        int
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
	gc := draw2dgl.NewGraphicContext(width, height)

	gc.Translate(380, 400)
	gc.Scale(1, -1)
	rotate = (rotate + 10) % 360
	gc.Rotate(float64(rotate) * math.Pi / 180)
	gc.Translate(-380, -400)
	interpreter := postscript.NewInterpreter(gc)
	reader := strings.NewReader(postscriptContent)
	lastTime := time.Nanoseconds()
	interpreter.Execute(reader)
	dt := time.Nanoseconds() - lastTime
	log.Printf("Redraw in : %f ms\n", float64(dt)*1e-6)
	gl.Flush() /* Single buffered, so needs a flush. */
	glut.PostRedisplay()
}

func main() {
	src, err := os.OpenFile("../resource/postscript/tiger.ps", 0, 0)
	if err != nil {
		log.Println("can't find postscript file.")
		return
	}
	defer src.Close()
	bytes, err := ioutil.ReadAll(src)
	postscriptContent = string(bytes)
	glut.Init()
	glut.InitWindowSize(800, 800)
	glut.CreateWindow("single triangle")

	glut.DisplayFunc(display)
	glut.ReshapeFunc(reshape)
	glut.MainLoop()
}
