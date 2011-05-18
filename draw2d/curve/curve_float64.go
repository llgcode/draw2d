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
	x1, y1, x2, y2, x3, y3, x4, y4 float64
}

func NewCubicCurveFloat64(x1, y1, x2, y2, x3, y3, x4, y4 float64) *CubicCurveFloat64 {
	return &CubicCurveFloat64{x1, y1, x2, y2, x3, y3, x4, y4}
}

//mu ranges from 0 to 1, start to end of curve
func (c *CubicCurveFloat64) ArbitraryPoint(mu float64) (x, y float64) {

	mum1 := 1 - mu
	mum13 := mum1 * mum1 * mum1
	mu3 := mu * mu * mu

	x = mum13*c.x1 + 3*mu*mum1*mum1*c.x2 + 3*mu*mu*mum1*c.x3 + mu3*c.x4
	y = mum13*c.y1 + 3*mu*mum1*mum1*c.y2 + 3*mu*mu*mum1*c.y3 + mu3*c.y4
	return
}

func (c *CubicCurveFloat64) SubdivideAt(c1, c2 *CubicCurveFloat64, t float64) {
  	inv_t := (1 - t) 
  	c1.x1, c1.y1 = c.x1, c.y1
  	c2.x4, c2.y4 = c.x4, c.y4
  	
    c1.x2 = inv_t * c.x1 + t * c.x2
	c1.y2 = inv_t * c.y1 + t * c.y2
	
    x23 := inv_t * c.x2 + t * c.x3
	y23 := inv_t * c.y2 + t * c.y3
	
    c2.x3 = inv_t * c.x3 + t * c.x4
	c2.y3 = inv_t * c.y3 + t * c.y4

    c1.x3 = inv_t * c1.x2 + t * x23
	c1.y3 = inv_t * c1.y2 + t * y23
	
    c2.x2 = inv_t * x23 + t * c2.x3
	c2.y2 = inv_t * y23 + t * c2.y3
	
    c1.x4 = inv_t * c1.x3 + t * c2.x2
	c1.y4 =  inv_t * c1.y3 + t * c2.y2
	
    c2.x1, c2.y1 = c1.x4, c1.y4
}

func (c *CubicCurveFloat64) Subdivide(c1, c2 *CubicCurveFloat64) {
	// Calculate all the mid-points of the line segments
	//----------------------
	c1.x1, c1.y1 = c.x1, c.y1
	c2.x4, c2.y4 = c.x4, c.y4
	c1.x2 = (c.x1 + c.x2) / 2
	c1.y2 = (c.y1 + c.y2) / 2
	x23 := (c.x2 + c.x3) / 2
	y23 := (c.y2 + c.y3) / 2
	c2.x3 = (c.x3 + c.x4) / 2
	c2.y3 = (c.y3 + c.y4) / 2
	c1.x3 = (c1.x2 + x23) / 2
	c1.y3 = (c1.y2 + y23) / 2
	c2.x2 = (x23 + c2.x3) / 2
	c2.y2 = (y23 + c2.y3) / 2
	c1.x4 = (c1.x3 + c2.x2) / 2
	c1.y4 = (c1.y3 + c2.y2) / 2
	c2.x1, c2.y1 = c1.x4, c1.y4
}

func (c *CubicCurveFloat64) EstimateDistance() float64 {
	dx1 := c.x2 - c.x1
	dy1 := c.y2 - c.y1
	dx2 := c.x3 - c.x2
	dy2 := c.y3 - c.y2
	dx3 := c.x4 - c.x3
	dy3 := c.y4 - c.y3
	return math.Sqrt(dx1*dx1+dy1*dy1) + math.Sqrt(dx2*dx2+dy2*dy2) + math.Sqrt(dx3*dx3+dy3*dy3)
}

// subdivide the curve in straight lines using Casteljau subdivision 
// and computing minimal distance tolerance
func (c *CubicCurveFloat64) SegmentCasteljauRec(segments []float64) []float64 {
	// reinit segments
	segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = c.x1
	segments[len(segments)-1] = c.y1
	segments = c.segmentCasteljauRec(segments)
	segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = c.x4
	segments[len(segments)-1] = c.y4
	return segments
}

func (c *CubicCurveFloat64) segmentCasteljauRec(segments []float64) []float64 {
	var c1, c2 CubicCurveFloat64
	c.Subdivide(&c1, &c2)

	// Try to approximate the full cubic curve by a single straight line
	//------------------
	dx := c.x4 - c.x1
	dy := c.y4 - c.y1

	d2 := math.Fabs(((c.x2-c.x4)*dy - (c.y2-c.y4)*dx))
	d3 := math.Fabs(((c.x3-c.x4)*dy - (c.y3-c.y4)*dx))

	if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy) {
		segments = segments[0 : len(segments)+2]
		segments[len(segments)-2] = c2.x1
		segments[len(segments)-1] = c2.y1
		return segments
	}
	// Continue subdivision
	//----------------------
	segments = c1.segmentCasteljauRec(segments)
	segments = c2.segmentCasteljauRec(segments)
	return segments
}

func (curve *CubicCurveFloat64) SegmentCasteljau(segments []float64) ([]float64) {
	var curves [32]CubicCurveFloat64
	curves[0] = *curve
	i := 0
	// current curve
	var c *CubicCurveFloat64
	var dx, dy, d2, d3 float64
	for i >= 0 {
		c = &curves[i]
		dx = c.x4 - c.x1
		dy = c.y4 - c.y1
	
		d2 = math.Fabs(((c.x2-c.x4)*dy - (c.y2-c.y4)*dx))
		d3 = math.Fabs(((c.x3-c.x4)*dy - (c.y3-c.y4)*dx))

		if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy)  || i == len(curves) - 1 {
	        segments = segments[0 : len(segments)+2]
			segments[len(segments)-2] = c.x1
			segments[len(segments)-1] = c.y1
	        i--;
	    } else {
	    	// second half of bezier go lower onto the stack
	   		c.Subdivide(&curves[i+1], &curves[i])
	        i++;
	    }
    }
    segments = segments[0 : len(segments)+2]
	segments[len(segments)-2] = curve.x1
	segments[len(segments)-1] = curve.y1
    return segments
}