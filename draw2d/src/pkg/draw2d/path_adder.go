package draw2d


import (
	"freetype-go.googlecode.com/hg/freetype/raster"
)


type VertexAdder struct {
	command VertexCommand
	adder   raster.Adder
}


func floatToPoint(x, y float) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}


func NewVertexAdder(adder raster.Adder) *VertexAdder {
	return &VertexAdder{VertexNoCommand, adder}
}

func (vertexAdder *VertexAdder) NextCommand(cmd VertexCommand) {
	vertexAdder.command = cmd
}

func (vertexAdder *VertexAdder) Vertex(x, y float) {
	switch vertexAdder.command {
	case VertexStartCommand:
		vertexAdder.adder.Start(floatToPoint(x, y))
	default:
		vertexAdder.adder.Add1(floatToPoint(x, y))
	}
	vertexAdder.command = VertexNoCommand
}
