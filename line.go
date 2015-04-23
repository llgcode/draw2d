// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

type LineMarker byte

const (
	// Mark the current point of the line as a join to it can draw some specific join Bevel, Miter, Rount
	LineJoinMarker LineMarker = iota
)

type LineBuilder interface {
	NextCommand(cmd LineMarker)
	// MoveTo Start a New line from the point (x, y)
	MoveTo(x, y float64)
	// LineTo Draw a line from the current position to the point (x, y)
	LineTo(x, y float64)
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

func (dc *LineBuilders) NextCommand(cmd LineMarker) {
	for _, converter := range dc.builders {
		converter.NextCommand(cmd)
	}
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
