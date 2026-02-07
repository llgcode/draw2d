// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2dbase

import (
	"math"

	"github.com/llgcode/draw2d"
)

type LineStroker struct {
	Flattener     Flattener
	HalfLineWidth float64
	Cap           draw2d.LineCap
	Join          draw2d.LineJoin
	vertices      []float64
	rewind        []float64
	center        []float64 // Store centerline points for cap calculations
	x, y, nx, ny  float64
}

func NewLineStroker(c draw2d.LineCap, j draw2d.LineJoin, flattener Flattener) *LineStroker {
	l := new(LineStroker)
	l.Flattener = flattener
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
		// Store centerline points for cap calculations
		l.center = append(l.center, x1, y1, x2, y2)
		l.x, l.y, l.nx, l.ny = x2, y2, nx, ny
	}
}

func (l *LineStroker) Close() {
	if len(l.vertices) > 1 {
		l.appendVertex(l.vertices[0], l.vertices[1], l.rewind[0], l.rewind[1])
	}
}

func (l *LineStroker) End() {
	if len(l.vertices) < 2 {
		l.Flattener.End()
		l.vertices = l.vertices[0:0]
		l.rewind = l.rewind[0:0]
		l.center = l.center[0:0]
		l.x, l.y, l.nx, l.ny = 0, 0, 0, 0
		return
	}
	
	// Start the stroke outline
	l.Flattener.MoveTo(l.vertices[0], l.vertices[1])
	
	// Draw the first edge (vertices side)
	for i, j := 2, 3; j < len(l.vertices); i, j = i+2, j+2 {
		l.Flattener.LineTo(l.vertices[i], l.vertices[j])
	}
	
	// Apply cap at the end of the stroke
	lastIdx := len(l.vertices) - 2
	lastRewindIdx := len(l.rewind) - 2
	l.applyEndCap(l.vertices[lastIdx], l.vertices[lastIdx+1], l.rewind[lastRewindIdx], l.rewind[lastRewindIdx+1])
	
	// Draw the second edge (rewind side) in reverse
	for i, j := len(l.rewind)-2, len(l.rewind)-1; j > 0; i, j = i-2, j-2 {
		l.Flattener.LineTo(l.rewind[i], l.rewind[j])
	}
	
	// Apply cap at the start of the stroke
	l.applyStartCap(l.vertices[0], l.vertices[1], l.rewind[0], l.rewind[1])
	
	// Close the path
	l.Flattener.LineTo(l.vertices[0], l.vertices[1])
	
	l.Flattener.End()
	
	// reinit vertices
	l.vertices = l.vertices[0:0]
	l.rewind = l.rewind[0:0]
	l.center = l.center[0:0]
	l.x, l.y, l.nx, l.ny = 0, 0, 0, 0
}

// applyStartCap applies the appropriate line cap at the start of a stroke
// v1x, v1y: point on the "vertices" side (outer edge)
// v2x, v2y: point on the "rewind" side (inner edge)
func (l *LineStroker) applyStartCap(v1x, v1y, v2x, v2y float64) {
	if len(l.center) < 4 {
		return
	}
	
	// Get centerline point and direction at the start
	cx, cy := l.center[0], l.center[1]
	dx := l.center[2] - l.center[0]
	dy := l.center[3] - l.center[1]
	
	// Normalize direction
	d := vectorDistance(dx, dy)
	if d == 0 {
		return
	}
	dx /= d
	dy /= d
	
	switch l.Cap {
	case draw2d.ButtCap:
		// ButtCap: just connect the edges with a straight line
		// This is handled by the final LineTo(vertices[0], vertices[1])
		
	case draw2d.SquareCap:
		// SquareCap: extend backwards by HalfLineWidth
		// Add a small epsilon to ensure the edge pixel is included in rasterization
		extX := -dx * (l.HalfLineWidth + 0.5)
		extY := -dy * (l.HalfLineWidth + 0.5)
		
		// Draw the square cap
		l.Flattener.LineTo(v2x+extX, v2y+extY)
		l.Flattener.LineTo(v1x+extX, v1y+extY)
		
	case draw2d.RoundCap:
		// RoundCap: draw a semicircular arc from v2 back to v1
		// The arc should wrap around the start point
		angle1 := math.Atan2(v2y-cy, v2x-cx)
		angle2 := math.Atan2(v1y-cy, v1x-cx)
		
		// Ensure we go the short way around
		if angle2-angle1 > math.Pi {
			angle2 -= 2 * math.Pi
		} else if angle1-angle2 > math.Pi {
			angle2 += 2 * math.Pi
		}
		
		// Draw semicircle with 8 segments
		numSegments := 8
		for i := 1; i <= numSegments; i++ {
			t := float64(i) / float64(numSegments)
			angle := angle1 + t*(angle2-angle1)
			x := cx + l.HalfLineWidth*math.Cos(angle)
			y := cy + l.HalfLineWidth*math.Sin(angle)
			l.Flattener.LineTo(x, y)
		}
	}
}

// applyEndCap applies the appropriate line cap at the end of a stroke
// v1x, v1y: point on the "vertices" side (outer edge)
// v2x, v2y: point on the "rewind" side (inner edge)
func (l *LineStroker) applyEndCap(v1x, v1y, v2x, v2y float64) {
	if len(l.center) < 4 {
		return
	}
	
	// Get centerline point and direction at the end
	lastIdx := len(l.center) - 2
	cx, cy := l.center[lastIdx], l.center[lastIdx+1]
	dx := l.center[lastIdx] - l.center[lastIdx-2]
	dy := l.center[lastIdx+1] - l.center[lastIdx-1]
	
	// Normalize direction
	d := vectorDistance(dx, dy)
	if d == 0 {
		return
	}
	dx /= d
	dy /= d
	
	switch l.Cap {
	case draw2d.ButtCap:
		// ButtCap: just connect the edges with a straight line
		l.Flattener.LineTo(v2x, v2y)
		
	case draw2d.SquareCap:
		// SquareCap: extend forwards by HalfLineWidth
		// Add a small epsilon to ensure the edge pixel is included in rasterization
		extX := dx * (l.HalfLineWidth + 0.5)
		extY := dy * (l.HalfLineWidth + 0.5)
		
		// Draw the square cap
		l.Flattener.LineTo(v1x+extX, v1y+extY)
		l.Flattener.LineTo(v2x+extX, v2y+extY)
		l.Flattener.LineTo(v2x, v2y)
		
	case draw2d.RoundCap:
		// RoundCap: draw a semicircular arc from v1 to v2
		angle1 := math.Atan2(v1y-cy, v1x-cx)
		angle2 := math.Atan2(v2y-cy, v2x-cx)
		
		// Ensure we go the short way around
		if angle2-angle1 > math.Pi {
			angle2 -= 2 * math.Pi
		} else if angle1-angle2 > math.Pi {
			angle2 += 2 * math.Pi
		}
		
		// Draw semicircle with 8 segments
		numSegments := 8
		for i := 1; i <= numSegments; i++ {
			t := float64(i) / float64(numSegments)
			angle := angle1 + t*(angle2-angle1)
			x := cx + l.HalfLineWidth*math.Cos(angle)
			y := cy + l.HalfLineWidth*math.Sin(angle)
			l.Flattener.LineTo(x, y)
		}
	}
}

func (l *LineStroker) appendVertex(vertices ...float64) {
	s := len(vertices) / 2
	l.vertices = append(l.vertices, vertices[:s]...)
	l.rewind = append(l.rewind, vertices[s:]...)
}

func vectorDistance(dx, dy float64) float64 {
	return float64(math.Sqrt(dx*dx + dy*dy))
}
