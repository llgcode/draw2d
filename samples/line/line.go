// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff, Stani Michiels

// Package line draws vertically spaced lines.
package line

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/samples"
)

// Main draws vertically spaced lines and returns the filename.
// This should only be used during testing.
func Main(gc draw2d.GraphicContext, ext string) (string, error) {
	// Draw the line
	for x := 5.0; x < 297; x += 10 {
		Draw(gc, x, 0, x, 210)
	}

	// Return the output filename
	return samples.Output("line", ext), nil
}

// Draw vertically spaced lines
func Draw(gc draw2d.GraphicContext, x0, y0, x1, y1 float64) {
	// Draw a line
	gc.MoveTo(x0, y0)
	gc.LineTo(x1, y1)
	gc.Stroke()
}
