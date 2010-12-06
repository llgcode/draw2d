// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"fmt"
	"math"
)

type PathCmd int

const (
	MoveTo PathCmd = iota
	LineTo
	QuadCurveTo
	CubicCurveTo
	ArcTo
	Close
)

type PathStorage struct {
	commands []PathCmd
	vertices []float
	x, y     float
}

func (p *PathStorage) appendToPath(cmd PathCmd, vertices ...float) {
	p.commands = append(p.commands, cmd)
	p.vertices = append(p.vertices, vertices...)
}

func (src *PathStorage) Copy() (dest *PathStorage) {
	dest = new(PathStorage)
	dest.commands = make([]PathCmd, len(src.commands))
	copy(dest.commands, src.commands)
	dest.vertices = make([]float, len(src.vertices))
	copy(dest.vertices, src.vertices)
	return dest
}
func (p *PathStorage) LastPoint() (x, y float) {
	return p.x, p.y
}

func (p *PathStorage) isEmpty() bool {
	return len(p.commands) == 0
}

func (p *PathStorage) Close() *PathStorage {
	p.appendToPath(Close)
	return p
}

func (p *PathStorage) MoveTo(x, y float) *PathStorage {
	p.appendToPath(MoveTo, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RMoveTo(dx, dy float) *PathStorage {
	x, y := p.LastPoint()
	p.MoveTo(x+dx, y+dy)
	return p
}

func (p *PathStorage) LineTo(x, y float) *PathStorage {
	p.appendToPath(LineTo, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RLineTo(dx, dy float) *PathStorage {
	x, y := p.LastPoint()
	p.LineTo(x+dx, y+dy)
	return p
}

func (p *PathStorage) QuadCurveTo(cx, cy, x, y float) *PathStorage {
	p.appendToPath(QuadCurveTo, cx, cy, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RQuadCurveTo(dcx, dcy, dx, dy float) *PathStorage {
	x, y := p.LastPoint()
	p.RQuadCurveTo(x+dcx, y+dcy, x+dx, y+dy)
	return p
}

func (p *PathStorage) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float) *PathStorage {
	p.appendToPath(CubicCurveTo, cx1, cy1, cx2, cy2, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float) *PathStorage {
	x, y := p.LastPoint()
	p.RCubicCurveTo(x+dcx1, y+dcy1, x+dcx2, y+dcy2, x+dx, y+dy)
	return p
}

func (p *PathStorage) ArcTo(cx, cy, rx, ry, startAngle, angle float) *PathStorage {
	endAngle := startAngle + angle
	clockWise := true
	if angle < 0 {
		clockWise = false
	}
	// normalize
	if clockWise {
		for endAngle < startAngle {
			endAngle += math.Pi * 2.0
		}
	} else {
		for startAngle < endAngle {
			startAngle += math.Pi * 2.0
		}
	}
	startX := cx + cos(startAngle)*rx
	startY := cy + sin(startAngle)*ry
	if len(p.commands) > 0 {
		p.LineTo(startX, startY)
	} else {
		p.MoveTo(startX, startY)
	}
	p.appendToPath(ArcTo, cx, cy, rx, ry, startAngle, angle)
	p.x = cx + cos(endAngle)*rx
	p.y = cy + sin(endAngle)*ry
	return p
}

func (p *PathStorage) RArcTo(dcx, dcy, rx, ry, startAngle, angle float) *PathStorage {
	x, y := p.LastPoint()
	p.RArcTo(x+dcx, y+dcy, rx, ry, startAngle, angle)
	return p
}

func (p *PathStorage) TraceLine(tracer LineTracer, approximationScale float) {
	j := 0
	x, y := 0.0, 0.0
	firstX, firstY := x, y
	if len(p.commands) > 0 {
		if p.commands[0] == MoveTo {
			firstX, firstY = p.vertices[0], p.vertices[1]
		}
	}
	for _, cmd := range p.commands {
		switch cmd {
		case MoveTo:
			tracer.MoveTo(p.vertices[j], p.vertices[j+1])
			x, y = p.vertices[j], p.vertices[j+1]
			firstX, firstY = x, y
			j = j + 2
		case LineTo:
			tracer.LineTo(p.vertices[j], p.vertices[j+1])
			x, y = p.vertices[j], p.vertices[j+1]
			j = j + 2
		case QuadCurveTo:
			quadraticBezier(tracer, x, y, p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3], approximationScale, 0)
			x, y = p.vertices[j+2], p.vertices[j+3]
			j = j + 4
		case CubicCurveTo:
			cubicBezier(tracer, x, y, p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3], p.vertices[j+4], p.vertices[j+5], approximationScale, 0, 0)
			x, y = p.vertices[j+4], p.vertices[j+5]
			j = j + 6
		case ArcTo:
			arc(tracer, p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3], p.vertices[j+4], p.vertices[j+5], approximationScale)
			j = j + 6
		case Close:
			tracer.LineTo(firstX, firstY)
			x, y = firstX, firstY
		}
	}
}

func (p *PathStorage) String() string {
	s := ""
	j := 0
	for _, cmd := range p.commands {
		switch cmd {
		case MoveTo:
			s += fmt.Sprintf("MoveTo: %f, %f\n", p.vertices[j], p.vertices[j+1])
			j = j + 2
		case LineTo:
			s += fmt.Sprintf("LineTo: %f, %f\n", p.vertices[j], p.vertices[j+1])
			j = j + 2
		case QuadCurveTo:
			s += fmt.Sprintf("QuadCurveTo: %f, %f, %f, %f\n", p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3])
			j = j + 4
		case CubicCurveTo:
			s += fmt.Sprintf("CubicCurveTo: %f, %f, %f, %f, %f, %f\n", p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3], p.vertices[j+4], p.vertices[j+5])
			j = j + 6
		case ArcTo:
			s += fmt.Sprintf("ArcTo: %f, %f, %f, %f, %f, %f\n", p.vertices[j], p.vertices[j+1], p.vertices[j+2], p.vertices[j+3], p.vertices[j+4], p.vertices[j+5])
			j = j + 6
		case Close:
			s += "Close\n"
		}
	}
	return s
}
