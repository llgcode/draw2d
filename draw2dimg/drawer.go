package draw2dimg

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
	"image"
	"image/draw"
)

type Drawer struct {
	matrix           draw2d.Matrix
	img              draw.Image
	painter          Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	glyphBuf         *truetype.GlyphBuf
}

func NewDrawer(img *image.RGBA) *Drawer {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	return &Drawer{
		draw2d.NewIdentityMatrix(),
		img,
		raster.NewRGBAPainter(img),
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		truetype.NewGlyphBuf(),
	}
}

func (d *Drawer) Matrix() *draw2d.Matrix {
	return d.Matrix()
}

func (d *Drawer) Fill(path *draw2d.Path, style draw2d.FillStyle) {
	switch fillStyle := style.(type) {
	case draw2d.SolidFillStyle:
		d.fillRasterizer.UseNonZeroWinding = useNonZeroWinding(fillStyle.FillRule)
		d.painter.SetColor(fillStyle.Color)
	default:
		panic("FillStyle not supported")
	}

	flattener := draw2dbase.Transformer{d.matrix, draw2dbase.FtLineBuilder{d.fillRasterizer}}

	draw2dbase.Flatten(path, flattener, d.matrix.GetScale())

	d.fillRasterizer.Rasterize(d.painter)
	d.fillRasterizer.Clear()
}

func (d *Drawer) Stroke(path *draw2d.Path, style draw2d.StrokeStyle) {
	d.strokeRasterizer.UseNonZeroWinding = true

	stroker := draw2dbase.NewLineStroker(style.LineCap, style.LineJoin, draw2dbase.Transformer{d.matrix, draw2dbase.FtLineBuilder{d.strokeRasterizer}})
	stroker.HalfLineWidth = style.Width / 2

	var liner draw2dbase.Flattener
	if style.Dash != nil && len(style.Dash) > 0 {
		liner = draw2dbase.NewDashConverter(style.Dash, style.DashOffset, stroker)
	} else {
		liner = stroker
	}

	draw2dbase.Flatten(path, liner, d.matrix.GetScale())

	d.painter.SetColor(style.Color)
	d.strokeRasterizer.Rasterize(d.painter)
	d.strokeRasterizer.Clear()
}

func (d *Drawer) Text(text string, x, y float64, style draw2d.TextStyle) {

}

func (d *Drawer) Image(image image.Image, x, y float64, scaling draw2d.ImageScaling) {
}
