// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

// Path describes the interface for path drawing.
type Path interface {
	// LastPoint returns the current point of the path
	LastPoint() (x, y float64)
	// MoveTo creates a new subpath that start at the specified point
	MoveTo(x, y float64)
	// RMoveTo creates a new subpath that start at the specified point
	// relative to the current point
	RMoveTo(dx, dy float64)
	// LineTo adds a line to the current subpath
	LineTo(x, y float64)
	// RLineTo adds a line to the current subpath
	// relative to the current point
	RLineTo(dx, dy float64)
	// QuadCurveTo adds a quadratic Bézier curve to the current subpath
	QuadCurveTo(cx, cy, x, y float64)
	// QuadCurveTo adds a quadratic Bézier curve to the current subpath
	// relative to the current point
	RQuadCurveTo(dcx, dcy, dx, dy float64)
	// CubicCurveTo adds a cubic Bézier curve to the current subpath
	CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64)
	// RCubicCurveTo adds a cubic Bézier curve to the current subpath
	// relative to the current point
	RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float64)
	// ArcTo adds an arc to the current subpath
	ArcTo(cx, cy, rx, ry, startAngle, angle float64)
	// RArcTo adds an arc to the current subpath
	// relative to the current point
	RArcTo(dcx, dcy, rx, ry, startAngle, angle float64)
	// Close creates a line from the current point to the last MoveTo
	// point (if not the same) and mark the path as closed so the
	// first and last lines join nicely.
	Close()
}
