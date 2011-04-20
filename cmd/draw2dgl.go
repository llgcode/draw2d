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
)

type GLPainter struct {

}

func reshape(w, h int) {
	/* Because Gil specified "screen coordinates" (presumably with an
	   upper-left origin), this short bit of code sets up the coordinate
	   system to correspond to actual window coodrinates.  This code
	   wouldn't be required if you chose a (more typical in 3D) abstract
	   coordinate system. */

	gl.Viewport(0, 0, w, h)                       /* Establish viewing area to cover entire window. */
	gl.MatrixMode(gl.PROJECTION)                  /* Start modifying the projection matrix. */
	gl.LoadIdentity()                             /* Reset project matrix. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1) /* Map abstract coords directly to window coords. */
	gl.Scalef(1, -1, 1)                           /* Invert Y axis so increasing Y goes down. */
	gl.Translatef(0, float32(-h), 0)              /* Shift origin up to upper-left corner. */
}

func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Begin(gl.TRIANGLES)
	gl.Color3f(0.0, 0.0, 1.0) /* blue */
	gl.Vertex2i(0, 0)
	gl.Color3f(0.0, 1.0, 0.0) /* green */
	gl.Vertex2i(200, 200)
	gl.Color3f(1.0, 0.0, 0.0) /* red */
	gl.Vertex2i(20, 200)
	gl.End()
	gl.Flush() /* Single buffered, so needs a flush. */
}

func main() {
	glut.Init()
	glut.CreateWindow("single triangle")
	glut.DisplayFunc(display)
	glut.ReshapeFunc(reshape)
	glut.MainLoop()
}
