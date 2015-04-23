// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff

package curve

import (
	"math"
)

// x1, y1, cpx1, cpy2, x2, y2 float64
// type Quad [6]float64

// Subdivide a Bezier quad curve in 2 equivalents Bezier quad curves.
// c1 and c2 parameters are the resulting curves
func SubdivideQuad(c, c1, c2 []float64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[4], c2[5] = c[4], c[5]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2
	c2[2] = (c[2] + c[4]) / 2
	c2[3] = (c[3] + c[5]) / 2
	c1[4] = (c1[2] + c2[2]) / 2
	c1[5] = (c1[3] + c2[3]) / 2
	c2[0], c2[1] = c1[4], c1[5]
	return
}

// Trace generate lines subdividing the curve using a LineBuilder
// flattening_threshold helps determines the flattening expectation of the curve
func TraceQuad(t LineBuilder, quad []float64, flattening_threshold float64) {
	// Allocates curves stack
	var curves [CurveRecursionLimit * 6]float64
	copy(curves[0:6], quad[0:6])
	i := 0
	// current curve
	var c []float64
	var dx, dy, d float64

	for i >= 0 {
		c = curves[i*6:]
		dx = c[4] - c[0]
		dy = c[5] - c[1]

		d = math.Abs(((c[2]-c[4])*dy - (c[3]-c[5])*dx))

		// if it's flat then trace a line
		if (d*d) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c[4], c[5])
			i--
		} else {
			// second half of bezier go lower onto the stack
			SubdivideQuad(c, curves[(i+1)*6:], curves[i*6:])
			i++
		}
	}
}
