package draw2d


import(
	"freetype-go.googlecode.com/hg/freetype/raster"	
)


type PathAdder struct {
	adder raster.Adder
}


func floatToPoint(x, y float) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}

func tracePath(approximationScale float, adder raster.Adder, paths ...*PathStorage) {
	pathAdder := &PathAdder{adder}
	for _, path := range paths {
		path.TraceLine(pathAdder, approximationScale)
	}
}

func (pathAdder *PathAdder) MoveTo(x, y float) {
	pathAdder.adder.Start(floatToPoint(x, y))
}

func (pathAdder *PathAdder) LineTo(x, y float) {
	pathAdder.adder.Add1(floatToPoint(x, y))
}
