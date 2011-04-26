// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff
package draw2d


import (
	"freetype-go.googlecode.com/hg/freetype/raster"
)


type VertexAdder struct {
	command VertexCommand
	adder   raster.Adder
}


func floatToPoint(x, y float64) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}


func NewVertexAdder(adder raster.Adder) *VertexAdder {
	return &VertexAdder{VertexNoCommand, adder}
}

func (vertexAdder *VertexAdder) NextCommand(cmd VertexCommand) {
	vertexAdder.command = cmd
}

func (vertexAdder *VertexAdder) Vertex(x, y float64) {
	switch vertexAdder.command {
	case VertexStartCommand:
		vertexAdder.adder.Start(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
	default:
		vertexAdder.adder.Add1(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
	}
	vertexAdder.command = VertexNoCommand
}


type PathAdder struct {
	adder              raster.Adder
	lastPoint          raster.Point
	ApproximationScale float64
}

func NewPathAdder(adder raster.Adder) *PathAdder {
	return &PathAdder{adder, raster.Point{0, 0}, 1}
}


func (pathAdder *PathAdder) Convert(paths ...*PathStorage) {
	for _, path := range paths {
		j := 0
		for _, cmd := range path.commands {
			j = j + pathAdder.ConvertCommand(cmd, path.vertices[j:]...)
		}
	}
}


func (pathAdder *PathAdder) ConvertCommand(cmd PathCmd, vertices ...float64) int {
	switch cmd {
	case MoveTo:
		pathAdder.lastPoint = floatToPoint(vertices[0], vertices[1])
		pathAdder.adder.Start(pathAdder.lastPoint)
		return 2
	case LineTo:
		pathAdder.lastPoint = floatToPoint(vertices[0], vertices[1])
		pathAdder.adder.Add1(pathAdder.lastPoint)
		return 2
	case QuadCurveTo:
		pathAdder.lastPoint = floatToPoint(vertices[2], vertices[3])
		pathAdder.adder.Add2(floatToPoint(vertices[0], vertices[1]), pathAdder.lastPoint)
		return 4
	case CubicCurveTo:
		pathAdder.lastPoint = floatToPoint(vertices[4], vertices[5])
		pathAdder.adder.Add3(floatToPoint(vertices[0], vertices[1]), floatToPoint(vertices[2], vertices[3]), pathAdder.lastPoint)
		return 6
	case ArcTo:
		pathAdder.lastPoint = arcAdder(pathAdder.adder, vertices[0], vertices[1], vertices[2], vertices[3], vertices[4], vertices[5], pathAdder.ApproximationScale)
		pathAdder.adder.Add1(pathAdder.lastPoint)
		return 6
	case Close:
		pathAdder.adder.Add1(pathAdder.lastPoint)
		return 0
	}
	return 0
}
