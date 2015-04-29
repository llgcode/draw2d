// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2dbase

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"github.com/llgcode/draw2d"
	"math"
)

func toFtCap(c draw2d.LineCap) raster.Capper {
	switch c {
	case draw2d.RoundCap:
		return raster.RoundCapper
	case draw2d.ButtCap:
		return raster.ButtCapper
	case draw2d.SquareCap:
		return raster.SquareCapper
	}
	return raster.RoundCapper
}

func toFtJoin(j draw2d.LineJoin) raster.Joiner {
	switch j {
	case draw2d.RoundJoin:
		return raster.RoundJoiner
	case draw2d.BevelJoin:
		return raster.BevelJoiner
	}
	return raster.RoundJoiner
}

type LineStroker struct {
	Next          draw2d.Flattener
	HalfLineWidth float64
	Cap           draw2d.LineCap
	Join          draw2d.LineJoin
	vertices      []float64
	rewind        []float64
	x, y, nx, ny  float64
}

func NewLineStroker(c draw2d.LineCap, j draw2d.LineJoin, flattener draw2d.Flattener) *LineStroker {
	l := new(LineStroker)
	l.Next = flattener
	l.HalfLineWidth = 0.5
	l.Cap = c
	l.Join = j
	return l
}

func (l *LineStroker) MoveTo(x, y float64) {
	l.x, l.y = x, y
}

func (l *LineStroker) LineTo(x, y float64) {
	l.line(l.x, l.y, x, y)
}

func (l *LineStroker) LineJoin() {

}

func (l *LineStroker) line(x1, y1, x2, y2 float64) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)
	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		l.appendVertex(x1+nx, y1+ny, x2+nx, y2+ny, x1-nx, y1-ny, x2-nx, y2-ny)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}

func (l *LineStroker) Close() {
	if len(l.vertices) > 1 {
		l.appendVertex(l.vertices[0], l.vertices[1], l.rewind[0], l.rewind[1])
	}
}

func (l *LineStroker) End() {
	if len(l.vertices) > 1 {
		l.Next.MoveTo(l.vertices[0], l.vertices[1])
		for i, j := 2, 3; j < len(l.vertices); i, j = i+2, j+2 {
			l.Next.LineTo(l.vertices[i], l.vertices[j])
		}
	}
	for i, j := len(l.rewind)-2, len(l.rewind)-1; j > 0; i, j = i-2, j-2 {
		l.Next.LineTo(l.rewind[i], l.rewind[j])
	}
	if len(l.vertices) > 1 {
		l.Next.LineTo(l.vertices[0], l.vertices[1])
	}
	l.Next.End()
	// reinit vertices
	l.vertices = l.vertices[0:0]
	l.rewind = l.rewind[0:0]
	l.x, l.y, l.nx, l.ny = 0, 0, 0, 0

}

func (l *LineStroker) appendVertex(vertices ...float64) {
	s := len(vertices) / 2
	l.vertices = append(l.vertices, vertices[:s]...)
	l.rewind = append(l.rewind, vertices[s:]...)
}

func vectorDistance(dx, dy float64) float64 {
	return float64(math.Sqrt(dx*dx + dy*dy))
}
