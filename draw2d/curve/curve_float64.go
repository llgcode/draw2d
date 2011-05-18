// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff
package curve

import (
	"math"
)

var (
	flattening_threshold float64 = 0.25
)

type CubicCurveFloat64 struct {
	X1, Y1, X2, Y2, X3, Y3, X4, Y4 float64
}

//mu ranges from 0 to 1, start to end of curve
func (c *CubicCurveFloat64) ArbitraryPoint(mu float64) (x, y float64) {

	mum1 := 1 - mu
	mum13 := mum1 * mum1 * mum1
	mu3 := mu * mu * mu

	x = mum13*c.X1 + 3*mu*mum1*mum1*c.X2 + 3*mu*mu*mum1*c.X3 + mu3*c.X4
	y = mum13*c.Y1 + 3*mu*mum1*mum1*c.Y2 + 3*mu*mu*mum1*c.Y3 + mu3*c.Y4
	return
}

func (c *CubicCurveFloat64) SubdivideAt(c1, c2 *CubicCurveFloat64, t float64) {
	inv_t := (1 - t)
	c1.X1, c1.Y1 = c.X1, c.Y1
	c2.X4, c2.Y4 = c.X4, c.Y4

	c1.X2 = inv_t*c.X1 + t*c.X2
	c1.Y2 = inv_t*c.Y1 + t*c.Y2

	x23 := inv_t*c.X2 + t*c.X3
	y23 := inv_t*c.Y2 + t*c.Y3

	c2.X3 = inv_t*c.X3 + t*c.X4
	c2.Y3 = inv_t*c.Y3 + t*c.Y4

	c1.X3 = inv_t*c1.X2 + t*x23
	c1.Y3 = inv_t*c1.Y2 + t*y23

	c2.X2 = inv_t*x23 + t*c2.X3
	c2.Y2 = inv_t*y23 + t*c2.Y3

	c1.X4 = inv_t*c1.X3 + t*c2.X2
	c1.Y4 = inv_t*c1.Y3 + t*c2.Y2

	c2.X1, c2.Y1 = c1.X4, c1.Y4
}

func (c *CubicCurveFloat64) Subdivide(c1, c2 *CubicCurveFloat64) {
	// Calculate all the mid-points of the line segments
	//----------------------
	c1.X1, c1.Y1 = c.X1, c.Y1
	c2.X4, c2.Y4 = c.X4, c.Y4
	c1.X2 = (c.X1 + c.X2) / 2
	c1.Y2 = (c.Y1 + c.Y2) / 2
	x23 := (c.X2 + c.X3) / 2
	y23 := (c.Y2 + c.Y3) / 2
	c2.X3 = (c.X3 + c.X4) / 2
	c2.Y3 = (c.Y3 + c.Y4) / 2
	c1.X3 = (c1.X2 + x23) / 2
	c1.Y3 = (c1.Y2 + y23) / 2
	c2.X2 = (x23 + c2.X3) / 2
	c2.Y2 = (y23 + c2.Y3) / 2
	c1.X4 = (c1.X3 + c2.X2) / 2
	c1.Y4 = (c1.Y3 + c2.Y2) / 2
	c2.X1, c2.Y1 = c1.X4, c1.Y4
}

func (c *CubicCurveFloat64) EstimateDistance() float64 {
	dx1 := c.X2 - c.X1
	dy1 := c.Y2 - c.Y1
	dx2 := c.X3 - c.X2
	dy2 := c.Y3 - c.Y2
	dx3 := c.X4 - c.X3
	dy3 := c.Y4 - c.Y3
	return math.Sqrt(dx1*dx1+dy1*dy1) + math.Sqrt(dx2*dx2+dy2*dy2) + math.Sqrt(dx3*dx3+dy3*dy3)
}

// subdivide the curve in straight lines using straight line approximation and Casteljau recursive subdivision 
// and computing minimal distance tolerance
func (c *CubicCurveFloat64) SegmentRec(segments []float64) []float64 {
	// reinit segments
	segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = c.X1
	segments[len(segments)-1] = c.Y1
	segments = c.segmentRec(segments)
	segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = c.X4
	segments[len(segments)-1] = c.Y4
	return segments
}

func (c *CubicCurveFloat64) segmentRec(segments []float64) []float64 {
	var c1, c2 CubicCurveFloat64
	c.Subdivide(&c1, &c2)

	// Try to approximate the full cubic curve by a single straight line
	//------------------
	dx := c.X4 - c.X1
	dy := c.Y4 - c.Y1

	d2 := math.Fabs(((c.X2-c.X4)*dy - (c.Y2-c.Y4)*dx))
	d3 := math.Fabs(((c.X3-c.X4)*dy - (c.Y3-c.Y4)*dx))

	if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy) {
		segments = segments[0 : len(segments)+2]
		segments[len(segments)-2] = c2.X4
		segments[len(segments)-1] = c2.Y4
		return segments
	}
	// Continue subdivision
	//----------------------
	segments = c1.segmentRec(segments)
	segments = c2.segmentRec(segments)
	return segments
}

func (curve *CubicCurveFloat64) Segment(segments []float64) []float64 {
	// Add the first point
	segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = curve.X1
	segments[len(segments)-1] = curve.Y1

	var curves [32]CubicCurveFloat64
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
			segments = segments[0 : len(segments)+2]
			segments[len(segments)-2] = c.X4
			segments[len(segments)-1] = c.Y4
			i--
		} else {
			// second half of bezier go lower onto the stack
			c.Subdivide(&curves[i+1], &curves[i])
			i++
		}
	}
	return segments
}
