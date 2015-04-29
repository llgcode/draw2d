// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 06/12/2010 by Laurent Le Goff

package path

// LineBuilder defines drawing line methods
type LineBuilder interface {
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

type LineBuilders struct {
	builders []LineBuilder
}

func NewLineBuilders(builders ...LineBuilder) *LineBuilders {
	return &LineBuilders{builders}
}

func (dc *LineBuilders) MoveTo(x, y float64) {
	for _, converter := range dc.builders {
		converter.MoveTo(x, y)
	}
}

func (dc *LineBuilders) LineTo(x, y float64) {
	for _, converter := range dc.builders {
		converter.LineTo(x, y)
	}
}

func (dc *LineBuilders) LineJoin() {
	for _, converter := range dc.builders {
		converter.LineJoin()
	}
}

func (dc *LineBuilders) Close() {
	for _, converter := range dc.builders {
		converter.Close()
	}
}

func (dc *LineBuilders) End() {
	for _, converter := range dc.builders {
		converter.End()
	}
}
