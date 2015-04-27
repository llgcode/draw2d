// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 06/12/2010 by Laurent Le Goff

package path

import (
	"github.com/llgcode/draw2d/curve"
)

type PathConverter struct {
	converter          LineBuilder
	ApproximationScale float64
}

func NewPathConverter(converter LineBuilder) *PathConverter {
	return &PathConverter{converter, 1, 0, 0, 0, 0}
}

// may not been in path instead put it in a troke package thing
func (c *PathConverter) Interpret(liner LineBuilder, scale float64, paths ...*Path) {
	// First Point
	var startX, startY float64 = 0, 0
	// Current Point
	var x, y float64 = 0, 0
	for _, path := range paths {
		i := 0
		for _, cmd := range path.Components {
			switch cmd {
			case MoveToCmp:
				x, y = path.Points[i], path.Points[i+1]
				startX, startY = x, y
				if i != 0 {
					liner.End()
				}
				liner.MoveTo(x, y)
				i += 2
			case LineToCmp:
				x, y = path.Points[i], path.Points[i+1]
				liner.LineTo(x, y)
				liner.LineJoin()
				i += 2
			case QuadCurveToCmp:
				curve.TraceQuad(liner, path.Points[i-2:], 0.5)
				x, y = path.Points[i+2], path.Points[i+3]
				liner.LineTo(x, y)
				i += 4
			case CubicCurveToCmp:
				curve.TraceCubic(liner, path.Points[i-2:], 0.5)
				x, y = path.Points[i+4], path.Points[i+5]
				liner.LineTo(x, y)
				i += 6
			case ArcToCmp:
				x, y = arc(liner, path.Points[i], path.Points[i+1], path.Points[i+2], path.Points[i+3], path.Points[i+4], path.Points[i+5], scale)
				liner.LineTo(x, y)
				i += 6
			case CloseCmp:
				liner.LineTo(startX, startY)
				liner.Close()
			}
		}
		liner.End()
	}
}
