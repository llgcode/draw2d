package draw2d

import (
	"freetype-go.googlecode.com/hg/freetype/raster"
)

type VertexRasterizer struct {
	rasterizer   *raster.Rasterizer
	command VertexCommand
}


func NewVertexRasterizer(rasterizer *raster.Rasterizer) (*VertexRasterizer) {
	vr := new(VertexRasterizer)
	vr.rasterizer = rasterizer
	return vr
}


func (vr *VertexRasterizer) NextCommand(command VertexCommand) {
	vr.command = command
}

func (vr *VertexRasterizer) Vertex(x, y float) {
	switch vr.command {
	case VertexStartCommand:	
		vr.rasterizer.Start(floatToPoint(x,y))
	default:
		vr.rasterizer.Add1(floatToPoint(x,y))
	}
	vr.command = VertexNoCommand
}

func floatToPoint(x, y float) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}
