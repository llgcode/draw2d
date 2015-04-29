// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package path

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"github.com/llgcode/draw2d/curve"
)

type FtLineBuilder struct {
	adder raster.Adder
}

func NewFtLineBuilder(adder raster.Adder) *FtLineBuilder {
	return &FtLineBuilder{adder}
}

func (FtLineBuilder *FtLineBuilder) MoveTo(x, y float64) {
	FtLineBuilder.adder.Start(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (FtLineBuilder *FtLineBuilder) LineTo(x, y float64) {
	FtLineBuilder.adder.Add1(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (FtLineBuilder *FtLineBuilder) LineJoin() {
}

func (FtLineBuilder *FtLineBuilder) Close() {
}

func (FtLineBuilder *FtLineBuilder) End() {
}

type PathAdder struct {
	adder              raster.Adder
	firstPoint         raster.Point
	ApproximationScale float64
}

func NewPathAdder(adder raster.Adder) *PathAdder {
	return &PathAdder{adder, raster.Point{0, 0}, 1}
}

func (pathAdder *PathAdder) Convert(paths ...*Path) {
	for _, apath := range paths {
		j := 0
		for _, cmd := range apath.Components {
			switch cmd {
			case MoveToCmp:
				pathAdder.firstPoint = raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}
				pathAdder.adder.Start(pathAdder.firstPoint)
				j += 2
			case LineToCmp:
				pathAdder.adder.Add1(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)})
				j += 2
			case QuadCurveToCmp:
				pathAdder.adder.Add2(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}, raster.Point{raster.Fix32(apath.Points[j+2] * 256), raster.Fix32(apath.Points[j+3] * 256)})
				j += 4
			case CubicCurveToCmp:
				pathAdder.adder.Add3(raster.Point{raster.Fix32(apath.Points[j] * 256), raster.Fix32(apath.Points[j+1] * 256)}, raster.Point{raster.Fix32(apath.Points[j+2] * 256), raster.Fix32(apath.Points[j+3] * 256)}, raster.Point{raster.Fix32(apath.Points[j+4] * 256), raster.Fix32(apath.Points[j+5] * 256)})
				j += 6
			case ArcToCmp:
				lastPoint := curve.TraceArcFt(pathAdder.adder, apath.Points[j], apath.Points[j+1], apath.Points[j+2], apath.Points[j+3], apath.Points[j+4], apath.Points[j+5], pathAdder.ApproximationScale)
				pathAdder.adder.Add1(lastPoint)
				j += 6
			case CloseCmp:
				pathAdder.adder.Add1(pathAdder.firstPoint)
			}
		}
	}
}
