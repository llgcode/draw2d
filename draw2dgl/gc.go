package draw2dgl

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/freetype/raster"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
	"github.com/llgcode/draw2d/draw2dimg"
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
		a := uint8((s.Alpha * p.ca / M16) >> 8)

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
	*draw2dbase.StackGraphicContext
	painter          *Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	DPI              int
}

// NewGraphicContext creates a new Graphic context from an image.
func NewGraphicContext(width, height int) *GraphicContext {
	gc := &GraphicContext{
		draw2dbase.NewStackGraphicContext(),
		NewPainter(),
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		92,
	}
	return gc
}

// CreateStringPath creates a path from the string s at x, y, and returns the string width.
// The text is placed so that the left edge of the em square of the first character of s
// and the baseline intersect at x, y. The majority of the affected pixels will be
// above and to the right of the point, but some may be below or to the left.
// For example, drawing a string that starts with a 'J' in an italic font may
// affect pixels below and left of the point.
func (gc *GraphicContext) CreateStringPath(s string, x, y float64) float64 {
	return draw2dbase.CreateStringPath(gc, s, x, y)
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
// The the left edge of the em square of the first character of s
// and the baseline intersect at 0, 0 in the returned coordinates.
// Therefore the top and left coordinates may well be negative.
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	return draw2dbase.GetStringBounds(gc, s)
}

func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return gc.FillStringAt(text, 0, 0)
}

func (gc *GraphicContext) FillStringAt(text string, x, y float64) (width float64) {
	width = gc.CreateStringPath(text, x, y)
	gc.Fill()
	return width
}

func (gc *GraphicContext) StrokeString(text string) (width float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (width float64) {
	width = gc.CreateStringPath(text, x, y)
	gc.Stroke()
	return width
}

func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
}

func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
func (gc *GraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
}

//TODO
func (gc *GraphicContext) Clear() {
	panic("not implemented")
}

//TODO
func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	panic("not implemented")
}

//TODO
func (gc *GraphicContext) DrawImage(img image.Image) {
	panic("not implemented")
}

func (gc *GraphicContext) paint(rasterizer *raster.Rasterizer, color color.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.painter.Flush()
	gc.Current.Path.Clear()
}

func (gc *GraphicContext) Stroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2dbase.Transformer{Tr: gc.Current.Tr, Flattener: draw2dimg.FtLineBuilder{Adder: gc.strokeRasterizer}})
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	var liner draw2dbase.Flattener
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
	} else {
		liner = stroker
	}
	for _, p := range paths {
		draw2dbase.Flatten(p, liner, gc.Current.Tr.GetScale())
	}

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = useNonZeroWinding(gc.Current.FillRule)

	/**** first method ****/
	flattener := draw2dbase.Transformer{Tr: gc.Current.Tr, Flattener: draw2dimg.FtLineBuilder{Adder: gc.fillRasterizer}}
	for _, p := range paths {
		draw2dbase.Flatten(p, flattener, gc.Current.Tr.GetScale())
	}

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
}

func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = useNonZeroWinding(gc.Current.FillRule)
	gc.strokeRasterizer.UseNonZeroWinding = true

	flattener := draw2dbase.Transformer{Tr: gc.Current.Tr, Flattener: draw2dimg.FtLineBuilder{Adder: gc.fillRasterizer}}

	stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2dbase.Transformer{Tr: gc.Current.Tr, Flattener: draw2dimg.FtLineBuilder{Adder: gc.strokeRasterizer}})
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	var liner draw2dbase.Flattener
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
	} else {
		liner = stroker
	}

	demux := draw2dbase.DemuxFlattener{Flatteners: []draw2dbase.Flattener{flattener, liner}}
	for _, p := range paths {
		draw2dbase.Flatten(p, demux, gc.Current.Tr.GetScale())
	}

	// Fill
	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	// Stroke
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func useNonZeroWinding(f draw2d.FillRule) bool {
	switch f {
	case draw2d.FillRuleEvenOdd:
		return false
	case draw2d.FillRuleWinding:
		return true
	}
	return false
}
