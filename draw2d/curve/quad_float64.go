// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff
package curve

import (
	"math"
)

type QuadCurveFloat64 struct {
	X1, Y1, X2, Y2, X3, Y3 float64
}


func (c *QuadCurveFloat64) Subdivide(c1, c2 *QuadCurveFloat64) {
	// Calculate all the mid-points of the line segments
	//----------------------
	c1.X1, c1.Y1 = c.X1, c.Y1
	c2.X3, c2.Y3 = c.X3, c.Y3
	c1.X2 = (c.X1 + c.X2) / 2
	c1.Y2 = (c.Y1 + c.Y2) / 2
	c2.X2 = (c.X2 + c.X3) / 2
	c2.Y2 = (c.Y2 + c.Y3) / 2
	c1.X3 = (c1.X2 + c2.X2) / 2
	c1.Y3 = (c1.Y2 + c2.Y2) / 2
	c2.X1, c2.Y1 = c1.X3, c1.Y3
	return
}


func (curve *QuadCurveFloat64) Segment(t LineTracer, flattening_threshold float64) {
	var curves [CurveRecursionLimit]QuadCurveFloat64
	curves[0] = *curve
	i := 0
	// current curve
	var c *QuadCurveFloat64
	var dx, dy, d float64
	
	for i >= 0 {
		c = &curves[i]
		dx = c.X3 - c.X1
		dy = c.Y3 - c.Y1

		d = math.Fabs(((c.X2-c.X3)*dy - (c.Y2-c.Y3)*dx))

		if (d*d) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c.X3, c.Y3)
			i--
		} else {
			// second half of bezier go lower onto the stack
			c.Subdivide(&curves[i+1], &curves[i])
			i++
		}
	}
}
