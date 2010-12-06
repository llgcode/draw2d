// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

type Path interface {
	MoveTo(x, y float)
	RMoveTo(dx, dy float)
	LineTo(x, y float)
	RLineTo(dx, dy float)
	QuadCurveTo(cx, cy, x, y float)
	RQuadCurveTo(dcx, dcy, dx, dy float)
	CubicCurveTo(cx1, cy1, cx2, cy2, x, y float)
	RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float)
	ArcTo(cx, cy, rx, ry, startAngle, angle float)
	RArcTo(dcx, dcy, rx, ry, startAngle, angle float)
	Close()
}


type VertexCommand byte

const (
	VertexNoCommand VertexCommand = iota
	VertexStartCommand
	VertexJoinCommand
	VertexCloseCommand
	VertexStopCommand
)

type VertexConverter interface {
	NextCommand(cmd VertexCommand)
	Vertex(x, y float)
}
