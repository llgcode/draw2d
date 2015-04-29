// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 06/12/2010 by Laurent Le Goff

package draw2d

// Flattener receive segment definition
type Flattener interface {
	// MoveTo Start a New line from the point (x, y)
	MoveTo(x, y float64)
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
	// LineJoin add the most recent starting point to close the path to create a polygon
	LineJoin()
	// Close add the most recent starting point to close the path to create a polygon
	Close()
	// End mark the current line as finished so we can draw caps
	End()
}

type SegmentedPath struct {
	Points []float64
}

func (p *SegmentedPath) MoveTo(x, y float64) {
	p.Points = append(p.Points, x, y)
	// TODO need to mark this point as moveto
}

func (p *SegmentedPath) LineTo(x, y float64) {
	p.Points = append(p.Points, x, y)
}

func (p *SegmentedPath) LineJoin() {
	// TODO need to mark the current point as linejoin
}

func (p *SegmentedPath) Close() {
	// TODO Close
}

func (p *SegmentedPath) End() {
	// Nothing to do
}

type LineCap int

const (
	RoundCap LineCap = iota
	ButtCap
	SquareCap
)

type LineJoin int

const (
	BevelJoin LineJoin = iota
	RoundJoin
	MiterJoin
)
