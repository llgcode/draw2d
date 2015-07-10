package draw2dgl

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"

	"code.google.com/p/freetype-go/freetype/raster"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/llgcode/draw2d"
)

func init() {
	runtime.LockOSThread()
}

type Painter struct {
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
func (p *Painter) Paint(ss []raster.Span, done bool) {
	//gl.Begin(gl.LINES)
	sslen := len(ss)
	clenrequired := sslen * 8
	vlenrequired := sslen * 4
	if clenrequired >= (cap(p.colors) - len(p.colors)) {
		p.Flush()

		if clenrequired >= cap(p.colors) {
			p.vertices = make([]int32, 0, vlenrequired+(vlenrequired/2))
			p.colors = make([]uint8, 0, clenrequired+(clenrequired/2))
		}
	}
	vi := len(p.vertices)
	ci := len(p.colors)
	p.vertices = p.vertices[0 : vi+vlenrequired]
	p.colors = p.colors[0 : ci+clenrequired]
	var (
		colors   []uint8
		vertices []int32
	)
	for _, s := range ss {
		ma := s.A >> 16
		a := uint8((ma * p.ca / M16) >> 8)
		colors = p.colors[ci:]
		colors[0] = p.cr
		colors[1] = p.cg
		colors[2] = p.cb
		colors[3] = a
		colors[4] = p.cr
		colors[5] = p.cg
		colors[6] = p.cb
		colors[7] = a
		ci += 8
		vertices = p.vertices[vi:]
		vertices[0] = int32(s.X0)
		vertices[1] = int32(s.Y)
		vertices[2] = int32(s.X1)
		vertices[3] = int32(s.Y)
		vi += 4
	}
}

func (p *Painter) Flush() {
	if len(p.vertices) != 0 {
		gl.EnableClientState(gl.COLOR_ARRAY)
		gl.EnableClientState(gl.VERTEX_ARRAY)
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, 0, gl.Ptr(p.colors))
		gl.VertexPointer(2, gl.INT, 0, gl.Ptr(p.vertices))

		// draw lines
		gl.DrawArrays(gl.LINES, 0, int32(len(p.vertices)/2))
		gl.DisableClientState(gl.VERTEX_ARRAY)
		gl.DisableClientState(gl.COLOR_ARRAY)
		p.vertices = p.vertices[0:0]
		p.colors = p.colors[0:0]
	}
}

// SetColor sets the color to paint the spans.
func (p *Painter) SetColor(c color.Color) {
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
func NewPainter() *Painter {
	p := new(Painter)
	p.vertices = make([]int32, 0, 1024)
	p.colors = make([]uint8, 0, 1024)
	return p
}

type GraphicContext struct {
	*draw2d.StackGraphicContext
	painter          *Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
}

// NewGraphicContext creates a new Graphic context from an image.
func NewGraphicContext(width, height int) *GraphicContext {
	gc := &GraphicContext{
		draw2d.NewStackGraphicContext(),
		NewPainter(),
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
	}
	return gc
}

func (gc *GraphicContext) CreateStringPath(s string, x, y float64) float64 {
	panic("not implemented")
}

func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	panic("not implemented")
}

func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	panic("not implemented")
}

func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	width := gc.CreateStringPath(text, x, y)
	gc.Stroke()
	return width
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

func (gc *GraphicContext) paint(rasterizer *raster.Rasterizer, color color.Color) {
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
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale() // From agg code
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
	gc.Current.Path = draw2d.NewPathStorage()
}

func (gc *GraphicContext) Fill(paths ...*draw2d.PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()

	pathConverter := draw2d.NewPathConverter(draw2d.NewVertexMatrixTransform(gc.Current.Tr, draw2d.NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale() // From agg code
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.Current.Path = draw2d.NewPathStorage()
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
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale() // From agg code
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
	gc.Current.Path = draw2d.NewPathStorage()
}
