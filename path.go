// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

// PathBuilder define method that create path
type PathBuilder interface {
	// Return the current point of the current path
	LastPoint() (x, y float64)

	// MoveTo start a new path at (x, y) position
	MoveTo(x, y float64)

	// LineTo add a line to the current path
	LineTo(x, y float64)

	// QuadCurveTo add a quadratic curve to the current path
	QuadCurveTo(cx, cy, x, y float64)

	// CubicCurveTo add a cubic bezier curve to the current path
	CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64)

	// ArcTo add an arc to the path
	ArcTo(cx, cy, rx, ry, startAngle, angle float64)

	// Close the current path
	Close()
}
