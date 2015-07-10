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
	Commands []PathCmd
	Vertices []float64
	x, y     float64
}

func NewPathStorage() (p *PathStorage) {
	p = new(PathStorage)
	p.Commands = make([]PathCmd, 0, 256)
	p.Vertices = make([]float64, 0, 256)
	return
}

func (p *PathStorage) appendToPath(cmd PathCmd, Vertices ...float64) {
	if cap(p.Vertices) <= len(p.Vertices)+6 {
		a := make([]PathCmd, len(p.Commands), cap(p.Commands)+256)
		b := make([]float64, len(p.Vertices), cap(p.Vertices)+256)
		copy(a, p.Commands)
		p.Commands = a
		copy(b, p.Vertices)
		p.Vertices = b
	}
	p.Commands = p.Commands[0 : len(p.Commands)+1]
	p.Commands[len(p.Commands)-1] = cmd
	copy(p.Vertices[len(p.Vertices):len(p.Vertices)+len(Vertices)], Vertices)
	p.Vertices = p.Vertices[0 : len(p.Vertices)+len(Vertices)]
}

func (src *PathStorage) Copy() (dest *PathStorage) {
	dest = new(PathStorage)
	dest.Commands = make([]PathCmd, len(src.Commands))
	copy(dest.Commands, src.Commands)
	dest.Vertices = make([]float64, len(src.Vertices))
	copy(dest.Vertices, src.Vertices)
	return dest
}

func (p *PathStorage) LastPoint() (x, y float64) {
	return p.x, p.y
}

func (p *PathStorage) IsEmpty() bool {
	return len(p.Commands) == 0
}

func (p *PathStorage) Close() *PathStorage {
	p.appendToPath(Close)
	return p
}

func (p *PathStorage) MoveTo(x, y float64) *PathStorage {
	p.appendToPath(MoveTo, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RMoveTo(dx, dy float64) *PathStorage {
	x, y := p.LastPoint()
	p.MoveTo(x+dx, y+dy)
	return p
}

func (p *PathStorage) LineTo(x, y float64) *PathStorage {
	p.appendToPath(LineTo, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RLineTo(dx, dy float64) *PathStorage {
	x, y := p.LastPoint()
	p.LineTo(x+dx, y+dy)
	return p
}

func (p *PathStorage) QuadCurveTo(cx, cy, x, y float64) *PathStorage {
	p.appendToPath(QuadCurveTo, cx, cy, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RQuadCurveTo(dcx, dcy, dx, dy float64) *PathStorage {
	x, y := p.LastPoint()
	p.QuadCurveTo(x+dcx, y+dcy, x+dx, y+dy)
	return p
}

func (p *PathStorage) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64) *PathStorage {
	p.appendToPath(CubicCurveTo, cx1, cy1, cx2, cy2, x, y)
	p.x = x
	p.y = y
	return p
}

func (p *PathStorage) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float64) *PathStorage {
	x, y := p.LastPoint()
	p.CubicCurveTo(x+dcx1, y+dcy1, x+dcx2, y+dcy2, x+dx, y+dy)
	return p
}

func (p *PathStorage) ArcTo(cx, cy, rx, ry, startAngle, angle float64) *PathStorage {
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
	startX := cx + math.Cos(startAngle)*rx
	startY := cy + math.Sin(startAngle)*ry
	if len(p.Commands) > 0 {
		p.LineTo(startX, startY)
	} else {
		p.MoveTo(startX, startY)
	}
	p.appendToPath(ArcTo, cx, cy, rx, ry, startAngle, angle)
	p.x = cx + math.Cos(endAngle)*rx
	p.y = cy + math.Sin(endAngle)*ry
	return p
}

func (p *PathStorage) RArcTo(dcx, dcy, rx, ry, startAngle, angle float64) *PathStorage {
	x, y := p.LastPoint()
	p.ArcTo(x+dcx, y+dcy, rx, ry, startAngle, angle)
	return p
}

func (p *PathStorage) String() string {
	s := ""
	j := 0
	for _, cmd := range p.Commands {
		switch cmd {
		case MoveTo:
			s += fmt.Sprintf("MoveTo: %f, %f\n", p.Vertices[j], p.Vertices[j+1])
			j = j + 2
		case LineTo:
			s += fmt.Sprintf("LineTo: %f, %f\n", p.Vertices[j], p.Vertices[j+1])
			j = j + 2
		case QuadCurveTo:
			s += fmt.Sprintf("QuadCurveTo: %f, %f, %f, %f\n", p.Vertices[j], p.Vertices[j+1], p.Vertices[j+2], p.Vertices[j+3])
			j = j + 4
		case CubicCurveTo:
			s += fmt.Sprintf("CubicCurveTo: %f, %f, %f, %f, %f, %f\n", p.Vertices[j], p.Vertices[j+1], p.Vertices[j+2], p.Vertices[j+3], p.Vertices[j+4], p.Vertices[j+5])
			j = j + 6
		case ArcTo:
			s += fmt.Sprintf("ArcTo: %f, %f, %f, %f, %f, %f\n", p.Vertices[j], p.Vertices[j+1], p.Vertices[j+2], p.Vertices[j+3], p.Vertices[j+4], p.Vertices[j+5])
			j = j + 6
		case Close:
			s += "Close\n"
		}
	}
	return s
}
