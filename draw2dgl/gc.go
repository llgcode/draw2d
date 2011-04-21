package draw2dgl

import (
	"image"
	"exp/draw"
	"gl"
	"unsafe"
	"freetype-go.googlecode.com/hg/freetype/raster"
	"draw2d.googlecode.com/hg/draw2d"
	//"log"
)

type GLPainter struct {
	// The Porter-Duff composition operator.
	Op draw.Op
	// The 16-bit color to paint the spans.
	cr, cg, cb uint8
	ca         uint32
	colors     []uint8
	vertices   []int32
}

const M16 uint32 = 1<<16 - 1

// Paint satisfies the Painter interface by painting ss onto an image.RGBA.
func (p *GLPainter) Paint(ss []raster.Span, done bool) {
	//gl.Begin(gl.LINES)
	for _, s := range ss {
		ma := s.A >> 16
		a := ma * p.ca / M16
		/*gl.Color4ub(p.cr, p.cg, p.cb, uint8(a>>8))
		gl.Vertex2i(s.X0, s.Y)
		gl.Vertex2i(s.X1, s.Y)*/
		p.colors = append(p.colors, p.cr, p.cg, p.cb, uint8(a>>8), p.cr, p.cg, p.cb, uint8(a>>8))
		p.vertices = append(p.vertices, int32(s.X0), int32(s.Y), int32(s.X1), int32(s.Y))
	}
	//gl.End()
}

func (p *GLPainter) Flush() {
	if len(p.vertices) != 0 {
		gl.EnableClientState(gl.COLOR_ARRAY)
		gl.EnableClientState(gl.VERTEX_ARRAY)
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, 0, unsafe.Pointer(&(p.colors[0])))
		gl.VertexPointer(2, gl.INT, 0, unsafe.Pointer(&(p.vertices[0])))

		// draw lines
		gl.DrawArrays(gl.LINES, 0, len(p.vertices)/2)
		gl.DisableClientState(gl.VERTEX_ARRAY)
		gl.DisableClientState(gl.COLOR_ARRAY)
		p.vertices = make([]int32, 0, 1024)
		p.colors = make([]uint8, 0, 1024)
	}
}


// SetColor sets the color to paint the spans.
func (p *GLPainter) SetColor(c image.Color) {
	r, g, b, a := c.RGBA()
	if a == 0 {
		p.cr = 0
		p.cg = 0
		p.cb = 0
		p.ca = a
	} else {
		p.cr = uint8((r * M16 / a) >> 8)
		p.cg = uint8((g * M16 / a) >> 8)
		p.cb = uint8((b * M16 / a) >> 8)
		p.ca = a
	}
}

// NewRGBAPainter creates a new RGBAPainter for the given image.
func NewGLPainter() *GLPainter {
	p := new(GLPainter)
	p.vertices = make([]int32, 0, 1024)
	p.colors = make([]uint8, 0, 1024)
	return p
}

type GraphicContext struct {
	*draw2d.StackGraphicContext
	painter          *GLPainter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(width, height int) *GraphicContext {
	gc := &GraphicContext{
		draw2d.NewStackGraphicContext(),
		NewGLPainter(),
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
	}
	return gc
}

func (gc *GraphicContext) SetDPI(dpi int) {

}

func (gc *GraphicContext) GetDPI() int {
	return -1
}

func (gc *GraphicContext) Clear() {

}

func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {

}

func (gc *GraphicContext) DrawImage(img image.Image) {

}


func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return 0
}


func (gc *GraphicContext) paint(rasterizer *raster.Rasterizer, color image.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.painter.Flush()
}

func (gc *GraphicContext) Stroke(paths ...*draw2d.PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := draw2d.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2d.NewVertexMatrixTransform(gc.Current.Tr, draw2d.NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2
	var pathConverter *draw2d.PathConverter
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		dasher := draw2d.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		pathConverter = draw2d.NewPathConverter(dasher)
	} else {
		pathConverter = draw2d.NewPathConverter(stroker)
	}
	pathConverter.ApproximationScale = gc.Current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
	gc.Current.Path = new(draw2d.PathStorage)
}

func (gc *GraphicContext) Fill(paths ...*draw2d.PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()

	pathConverter := draw2d.NewPathConverter(draw2d.NewVertexMatrixTransform(gc.Current.Tr, draw2d.NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.Current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.Current.Path = new(draw2d.PathStorage)
}

func (gc *GraphicContext) FillStroke(paths ...*draw2d.PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := draw2d.NewVertexMatrixTransform(gc.Current.Tr, draw2d.NewVertexAdder(gc.fillRasterizer))

	stroker := draw2d.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2d.NewVertexMatrixTransform(gc.Current.Tr, draw2d.NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	demux := draw2d.NewDemuxConverter(filler, stroker)
	paths = append(paths, gc.Current.Path)
	pathConverter := draw2d.NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.Current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
	gc.Current.Path = new(draw2d.PathStorage)
}
