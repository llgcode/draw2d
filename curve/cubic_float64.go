// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff

// Package curve implements Bezier Curve Subdivision using De Casteljau's algorithm
package curve

import (
	"math"
)

const (
	CurveRecursionLimit = 32
)

//	x1, y1, cpx1, cpx2, cpx2, cpy2, x2, y2 float64
type CubicCurveFloat64 [8]float64

// Subdivide a Bezier cubic curve in 2 equivalents Bezier cubic curves.
// c1 and c2 parameters are the resulting curves
func (c *CubicCurveFloat64) Subdivide(c1, c2 *CubicCurveFloat64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[6], c2[7] = c[6], c[7]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2

	midX := (c[2] + c[4]) / 2
	midY := (c[3] + c[5]) / 2

	c2[4] = (c[4] + c[6]) / 2
	c2[5] = (c[5] + c[7]) / 2

	c1[4] = (c1[2] + midX) / 2
	c1[5] = (c1[3] + midY) / 2

	c2[2] = (midX + c2[4]) / 2
	c2[3] = (midY + c2[5]) / 2

	c1[6] = (c1[4] + c2[2]) / 2
	c1[7] = (c1[5] + c2[3]) / 2

	// Last Point of c1 is equal to the first point of c2
	c2[0], c2[1] = c1[6], c1[7]
}

// Trace generate lines subdividing the curve using a LineTracer
// flattening_threshold helps determines the flattening expectation of the curve
func (curve *CubicCurveFloat64) Trace(t LineTracer, flattening_threshold float64) {
	// Allocation curves
	var curves [CurveRecursionLimit]CubicCurveFloat64
	curves[0] = *curve
	i := 0

	// current curve
	var c *CubicCurveFloat64

	var dx, dy, d2, d3 float64

	for i >= 0 {
		c = &curves[i]
		dx = c[6] - c[0]
		dy = c[7] - c[1]

		d2 = math.Abs((c[2]-c[6])*dy - (c[3]-c[7])*dx)
		d3 = math.Abs((c[4]-c[6])*dy - (c[5]-c[7])*dx)

		// if it's flat then trace a line
		if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c[6], c[7])
			i--
		} else {
			// second half of bezier go lower onto the stack
			c.Subdivide(&curves[i+1], &curves[i])
			i++
		}
	}
}
