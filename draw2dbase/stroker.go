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
	pendingJoin   bool // Flag to indicate if we need to process a join
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
	if l.pendingJoin && len(l.vertices) >= 4 {
		// Process the join before adding the new line segment
		l.processJoin(x, y)
	}
	l.line(l.x, l.y, x, y)
	l.pendingJoin = false
}

func (l *LineStroker) LineJoin() {
	// Mark that a join is needed before the next segment
	l.pendingJoin = true
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

// processJoin handles the join between the current segment and the next segment
func (l *LineStroker) processJoin(nextX, nextY float64) {
	// Get the current position and normal
	prevX, prevY := l.x, l.y
	prevNX, prevNY := l.nx, l.ny
	
	// Calculate the normal for the next segment
	dx := nextX - prevX
	dy := nextY - prevY
	d := vectorDistance(dx, dy)
	if d == 0 {
		return
	}
	nextNX := dy * l.HalfLineWidth / d
	nextNY := -(dx * l.HalfLineWidth / d)
	
	// The join point is at (prevX, prevY)
	// We need to connect the offset edges from the previous segment to the next segment
	
	// Previous segment ends at:
	// - outer edge: (prevX + prevNX, prevY + prevNY)
	// - inner edge: (prevX - prevNX, prevY - prevNY)
	
	// Next segment starts at:
	// - outer edge: (prevX + nextNX, prevY + nextNY)
	// - inner edge: (prevX - nextNX, prevY - nextNY)
	
	// Determine which side needs the join (outer or inner)
	// This is determined by the turn direction (cross product of the two direction vectors)
	
	// Get the direction of the previous segment (from the last two centerline points)
	if len(l.center) < 4 {
		return
	}
	
	lastCenterIdx := len(l.center) - 2
	prevDX := prevX - l.center[lastCenterIdx-2]
	prevDY := prevY - l.center[lastCenterIdx-1]
	
	// Cross product to determine turn direction
	// positive = left turn (counterclockwise), negative = right turn (clockwise)
	cross := prevDX*dy - prevDY*dx
	
	// For the outer edge (vertices side), we need a join if the normals don't align
	// For simplicity, we'll apply the join style based on the angle between segments
	
	switch l.Join {
	case draw2d.BevelJoin:
		// Bevel join: simply connect the two edges with a straight line
		// This is implicitly handled by the vertices already, no extra work needed
		
	case draw2d.RoundJoin:
		// Round join: add an arc between the two edges
		// We need to add vertices along the outer edge
		if cross != 0 {
			// Determine which edge is the outer edge based on turn direction
			var centerX, centerY, startAngle, endAngle float64
			centerX, centerY = prevX, prevY
			
			if cross > 0 {
				// Left turn - outer edge is on the vertices side
				startAngle = math.Atan2(prevNY, prevNX)
				endAngle = math.Atan2(nextNY, nextNX)
			} else {
				// Right turn - outer edge is on the rewind side
				startAngle = math.Atan2(-prevNY, -prevNX)
				endAngle = math.Atan2(-nextNY, -nextNX)
			}
			
			// Normalize angle difference
			angleDiff := endAngle - startAngle
			if angleDiff > math.Pi {
				endAngle -= 2 * math.Pi
			} else if angleDiff < -math.Pi {
				endAngle += 2 * math.Pi
			}
			
			// Add arc vertices for the round join
			// We'll add them directly to the vertices or rewind array as needed
			// For now, let's just add intermediate points
			numSegments := 4
			for i := 1; i < numSegments; i++ {
				t := float64(i) / float64(numSegments)
				angle := startAngle + t*(endAngle-startAngle)
				vx := centerX + l.HalfLineWidth*math.Cos(angle)
				vy := centerY + l.HalfLineWidth*math.Sin(angle)
				
				if cross > 0 {
					l.vertices = append(l.vertices, vx, vy)
				} else {
					// For inner edge, we prepend to rewind (will be reversed later)
					l.rewind = append(l.rewind, vx, vy)
				}
			}
		}
		
	case draw2d.MiterJoin:
		// Miter join: extend the two edges until they meet
		// This can create very long spikes at sharp angles, so we may need a miter limit
		// For now, we'll implement a simple miter
		
		// Calculate the miter point where the two extended edges meet
		// This requires finding the intersection of two lines
		
		// Line 1: prevX + prevNX + t1 * prevDX, prevY + prevNY + t1 * prevDY
		// Line 2: prevX + nextNX + t2 * dx, prevY + nextNY + t2 * dy
		
		// For simplicity, we'll skip the miter calculation for now
		// and fall back to bevel behavior
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
