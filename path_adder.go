// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2d

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"github.com/llgcode/draw2d/path"
)

type VertexAdder struct {
	adder raster.Adder
}

func NewVertexAdder(adder raster.Adder) *VertexAdder {
	return &VertexAdder{adder}
}

func (vertexAdder *VertexAdder) MoveTo(x, y float64) {
	vertexAdder.adder.Start(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (vertexAdder *VertexAdder) LineTo(x, y float64) {
	vertexAdder.adder.Add1(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (vertexAdder *VertexAdder) LineJoin() {
}

func (vertexAdder *VertexAdder) Close() {
}

func (vertexAdder *VertexAdder) End() {
}

type PathAdder struct {
	adder              raster.Adder
	firstPoint         raster.Point
	ApproximationScale float64
}

func NewPathAdder(adder raster.Adder) *PathAdder {
	return &PathAdder{adder, raster.Point{0, 0}, 1}
}

func (pathAdder *PathAdder) Convert(paths ...*path.Path) {
	for _, apath := range paths {
		j := 0
		for _, cmd := range apath.Components {
			switch cmd {
			case path.MoveToCmp:
				pathAdder.firstPoint = raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}
				pathAdder.adder.Start(pathAdder.firstPoint)
				j += 2
			case path.LineToCmp:
				pathAdder.adder.Add1(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)})
				j += 2
			case path.QuadCurveToCmp:
				pathAdder.adder.Add2(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}, raster.Point{raster.Fix32(apath.Points[j+2] * 256), raster.Fix32(apath.Points[j+3] * 256)})
				j += 4
			case path.CubicCurveToCmp:
				pathAdder.adder.Add3(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}, raster.Point{raster.Fix32(apath.Points[j+2] * 256), raster.Fix32(apath.Points[j+3] * 256)}, raster.Point{raster.Fix32(apath.Points[j+4] * 256), raster.Fix32(apath.Points[j+5] * 256)})
				j += 6
			case path.ArcToCmp:
				lastPoint := arcAdder(pathAdder.adder, apath.Points[j], apath.Points[j+1], apath.Points[j+2], apath.Points[j+3], apath.Points[j+4], apath.Points[j+5], pathAdder.ApproximationScale)
				pathAdder.adder.Add1(lastPoint)
				j += 6
			case path.CloseCmp:
				pathAdder.adder.Add1(pathAdder.firstPoint)
			}
		}
	}
}
