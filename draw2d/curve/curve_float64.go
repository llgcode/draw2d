// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff
package curve

import (
	"math"
)

var (
	m_distance_tolerance float64 = 0.25
)

type CubicCurveFloat64 struct {
	x1, y1, x2, y2, x3, y3, x4, y4 float64
	segments                       []float64
}

func NewCubicCurveFloat64(x1, y1, x2, y2, x3, y3, x4, y4 float64) (*CubicCurveFloat64){
	return &CubicCurveFloat64{x1, y1, x2, y2, x3, y3, x4, y4, make([]float64, 0, 2)}
}

func (c *CubicCurveFloat64) addPoint(x, y float64) {
	c.segments = append(c.segments, x, y)
}

func (c *CubicCurveFloat64) Subdivide() (c1, c2 *CubicCurveFloat64) {
	// Calculate all the mid-points of the line segments
	//----------------------
	x12 := (c.x1 + c.x2) / 2
	y12 := (c.y1 + c.y2) / 2
	x23 := (c.x2 + c.x3) / 2
	y23 := (c.y2 + c.y3) / 2
	x34 := (c.x3 + c.x4) / 2
	y34 := (c.y3 + c.y4) / 2
	x123 := (x12 + x23) / 2
	y123 := (y12 + y23) / 2
	x234 := (x23 + x34) / 2
	y234 := (y23 + y34) / 2
	x1234 := (x123 + x234) / 2
	y1234 := (y123 + y234) / 2
	c1 = &CubicCurveFloat64{c.x1, c.y1, x12, y12, x123, y123, x1234, y1234, c.segments}
	c2 = &CubicCurveFloat64{x1234, y1234, x234, y234, x34, y34, c.x4, c.y4, c.segments}
	return
}

// subdivide the curve in straight lines using Casteljau subdivision 
// and computing minimal distance tolerance
func (c *CubicCurveFloat64) SegmentCasteljau() []float64{
	// reinit segments
	c.segments = make([]float64, 0, 2)
	c.addPoint(c.x1, c.y1)
	c.segmentCasteljauRec()
	c.addPoint(c.x4, c.y4)
	return c.segments
}

func (c *CubicCurveFloat64) segmentCasteljauRec() {

	c1, c2 := c.Subdivide()

	// Try to approximate the full cubic curve by a single straight line
	//------------------
	dx := c.x4 - c.x1
	dy := c.y4 - c.y1

	d2 := math.Fabs(((c.x2-c.x4)*dy - (c.y2-c.y4)*dx))
	d3 := math.Fabs(((c.x3-c.x4)*dy - (c.y3-c.y4)*dx))

	if (d2+d3)*(d2+d3) < m_distance_tolerance*(dx*dx+dy*dy) {
		c.addPoint(c2.x1, c2.y1)
		return
	} else {
		// Continue subdivision
		//----------------------
		c1.segmentCasteljauRec()
		c2.segments = c1.segments 
		c2.segmentCasteljauRec()
		c.segments = c2.segments 
	}
}
