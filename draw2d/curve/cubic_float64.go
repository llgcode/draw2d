// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff
package curve

import (
	"math"
)

const (
	CurveRecursionLimit = 32
)

type CubicCurveFloat64 struct {
	X1, Y1, X2, Y2, X3, Y3, X4, Y4 float64
}

type LineTracer interface {
	LineTo(x, y float64)
}

func (c *CubicCurveFloat64) Subdivide(c1, c2 *CubicCurveFloat64) (x23, y23 float64) {
	// Calculate all the mid-points of the line segments
	//----------------------
	c1.X1, c1.Y1 = c.X1, c.Y1
	c2.X4, c2.Y4 = c.X4, c.Y4
	c1.X2 = (c.X1 + c.X2) / 2
	c1.Y2 = (c.Y1 + c.Y2) / 2
	x23 = (c.X2 + c.X3) / 2
	y23 = (c.Y2 + c.Y3) / 2
	c2.X3 = (c.X3 + c.X4) / 2
	c2.Y3 = (c.Y3 + c.Y4) / 2
	c1.X3 = (c1.X2 + x23) / 2
	c1.Y3 = (c1.Y2 + y23) / 2
	c2.X2 = (x23 + c2.X3) / 2
	c2.Y2 = (y23 + c2.Y3) / 2
	c1.X4 = (c1.X3 + c2.X2) / 2
	c1.Y4 = (c1.Y3 + c2.Y2) / 2
	c2.X1, c2.Y1 = c1.X4, c1.Y4
	return
}

func (curve *CubicCurveFloat64) Segment(t LineTracer, flattening_threshold float64) {
	// Add the first point
	t.LineTo(curve.X1, curve.Y1)

	var curves [CurveRecursionLimit]CubicCurveFloat64
	curves[0] = *curve
	i := 0
	// current curve
	var c *CubicCurveFloat64
	var dx, dy, d2, d3 float64
	for i >= 0 {
		c = &curves[i]
		dx = c.X4 - c.X1
		dy = c.Y4 - c.Y1

		d2 = math.Fabs(((c.X2-c.X4)*dy - (c.Y2-c.Y4)*dx))
		d3 = math.Fabs(((c.X3-c.X4)*dy - (c.Y3-c.Y4)*dx))

		if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c.X4, c.Y4)
			i--
		} else {
			// second half of bezier go lower onto the stack
			c.Subdivide(&curves[i+1], &curves[i])
			i++
		}
	}
}
