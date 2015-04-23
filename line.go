// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

type LineMarker byte

const (
	LineNoneMarker LineMarker = iota
	// Mark the current point of the line as a join to it can draw some specific join Bevel, Miter, Rount
	LineJoinMarker
	// Mark the current point of the line as closed so it draw a line from the current
	// position to the point specified by the last start marker.
	LineCloseMarker
)

type LineBuilder interface {
	NextCommand(cmd LineMarker)
	MoveTo(x, y float64)
	LineTo(x, y float64)
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

func (dc *LineBuilders) End() {
	for _, converter := range dc.builders {
		converter.End()
	}
}
