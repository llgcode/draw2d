// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2dbase

import (
	"code.google.com/p/freetype-go/freetype/raster"
)

type FtLineBuilder struct {
	Adder raster.Adder
}

func (liner FtLineBuilder) MoveTo(x, y float64) {
	liner.Adder.Start(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (liner FtLineBuilder) LineTo(x, y float64) {
	liner.Adder.Add1(raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)})
}

func (liner FtLineBuilder) LineJoin() {
}

func (liner FtLineBuilder) Close() {
}

func (liner FtLineBuilder) End() {
}
