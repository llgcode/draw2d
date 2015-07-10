// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels

package draw2dpdf

import (
	"math"

	"github.com/llgcode/draw2d"
)

const deg = 180 / math.Pi

// PathConverter converts the paths to the pdf api
type PathConverter struct {
	pdf Vectorizer
}

// NewPathConverter constructs a PathConverter from a pdf vectorizer
func NewPathConverter(pdf Vectorizer) *PathConverter {
	return &PathConverter{pdf: pdf}
}

// Convert converts the paths to the pdf api
func (c *PathConverter) Convert(paths ...*draw2d.PathStorage) {
	for _, path := range paths {
		j := 0
		for _, cmd := range path.Commands {
			j = j + c.ConvertCommand(cmd, path.Vertices[j:]...)
		}
	}
}

// ConvertCommand converts a single path segment to the pdf api
func (c *PathConverter) ConvertCommand(cmd draw2d.PathCmd, vertices ...float64) int {
	switch cmd {
	case draw2d.MoveTo:
		c.pdf.MoveTo(vertices[0], vertices[1])
		return 2
	case draw2d.LineTo:
		c.pdf.LineTo(vertices[0], vertices[1])
		return 2
	case draw2d.QuadCurveTo:
		c.pdf.CurveTo(vertices[0], vertices[1], vertices[2], vertices[3])
		return 4
	case draw2d.CubicCurveTo:
		c.pdf.CurveBezierCubicTo(vertices[0], vertices[1], vertices[2], vertices[3], vertices[4], vertices[5])
		return 6
	case draw2d.ArcTo:
		c.pdf.ArcTo(vertices[0], vertices[1], vertices[2], vertices[3],
			0,                             // degRotate
			vertices[4]*deg,               // degStart = startAngle
			(vertices[4]-vertices[5])*deg) // degEnd = startAngle-angle
		return 6
	default: // case draw2d.Close:
		c.pdf.ClosePath()
		return 0
	}
}
