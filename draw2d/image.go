// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"image"
	"image/color"
	"image/draw"
	"log"
)

type Painter interface {
	raster.Painter
	SetColor(color color.Color)
}

var (
	defaultFontData = FontData{"luxi", FontFamilySans, FontStyleNormal}
)

type ImageGraphicContext struct {
	*StackGraphicContext
	img              draw.Image
	painter          Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	freetype         *freetype.Context
	DPI              int
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(img draw.Image) *ImageGraphicContext {
	var painter Painter
	switch selectImage := img.(type) {
	case *image.RGBA:
		painter = raster.NewRGBAPainter(selectImage)
	default:
		panic("Image type not supported")
	}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	dpi := 92
	ftContext := freetype.NewContext()
	ftContext.SetDPI(dpi)
	ftContext.SetClip(img.Bounds())
	ftContext.SetDst(img)
	gc := &ImageGraphicContext{
		NewStackGraphicContext(),
		img,
		painter,
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		ftContext,
		dpi,
	}
	return gc
}

// Create a new Graphic context from an image and a Painter (see Freetype-go)
func NewGraphicContextWithPainter(img draw.Image, painter Painter) *ImageGraphicContext {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	dpi := 92
	ftContext := freetype.NewContext()
	ftContext.SetDPI(dpi)
	ftContext.SetClip(img.Bounds())
	ftContext.SetDst(img)
	gc := &ImageGraphicContext{
		NewStackGraphicContext(),
		img,
		painter,
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		ftContext,
		dpi,
	}
	return gc
}

func (gc *ImageGraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.freetype.SetDPI(dpi)
}

func (gc *ImageGraphicContext) GetDPI() int {
	return gc.DPI
}

func (gc *ImageGraphicContext) Clear() {
	width, height := gc.img.Bounds().Dx(), gc.img.Bounds().Dy()
	gc.ClearRect(0, 0, width, height)
}

func (gc *ImageGraphicContext) ClearRect(x1, y1, x2, y2 int) {
	imageColor := image.NewUniform(gc.Current.FillColor)
	draw.Draw(gc.img, image.Rect(x1, y1, x2, y2), imageColor, image.ZP, draw.Over)
}

func (gc *ImageGraphicContext) DrawImage(img image.Image) {
	DrawImage(img, gc.img, gc.Current.Tr, draw.Over, BilinearFilter)
}

func (gc *ImageGraphicContext) FillString(text string) (cursor float64) {
	gc.freetype.SetSrc(image.NewUniform(gc.Current.StrokeColor))
	// Draw the text.
	x, y := gc.Current.Path.LastPoint()
	gc.Current.Tr.Transform(&x, &y)
	x0, fontSize := 0.0, gc.Current.FontSize
	gc.Current.Tr.VectorTransform(&x0, &fontSize)
	font := GetFont(gc.Current.FontData)
	if font == nil {
		font = GetFont(defaultFontData)
	}
	if font == nil {
		return 0
	}
	gc.freetype.SetFont(font)
	gc.freetype.SetFontSize(fontSize)
	pt := freetype.Pt(int(x), int(y))
	p, err := gc.freetype.DrawString(text, pt)
	if err != nil {
		log.Println(err)
	}
	x1, _ := gc.Current.Path.LastPoint()
	x2, y2 := float64(p.X)/256, float64(p.Y)/256
	gc.Current.Tr.InverseTransform(&x2, &y2)
	width := x2 - x1
	return width
}

func (gc *ImageGraphicContext) paint(rasterizer *raster.Rasterizer, color color.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.Current.Path.Clear()
}

/**** second method ****/
func (gc *ImageGraphicContext) Stroke(paths ...*PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := NewLineStroker(gc.Current.Cap, gc.Current.Join, NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2
	var pathConverter *PathConverter
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		dasher := NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(stroker)
	}
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()

	/**** first method ****/
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
}

/* second method */
func (gc *ImageGraphicContext) FillStroke(paths ...*PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.fillRasterizer))

	stroker := NewLineStroker(gc.Current.Cap, gc.Current.Join, NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	demux := NewDemuxConverter(filler, stroker)
	paths = append(paths, gc.Current.Path)
	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func (f FillRule) UseNonZeroWinding() bool {
	switch f {
	case FillRuleEvenOdd:
		return false
	case FillRuleWinding:
		return true
	}
	return false
}

func (c Cap) Convert() raster.Capper {
	switch c {
	case RoundCap:
		return raster.RoundCapper
	case ButtCap:
		return raster.ButtCapper
	case SquareCap:
		return raster.SquareCapper
	}
	return raster.RoundCapper
}

func (j Join) Convert() raster.Joiner {
	switch j {
	case RoundJoin:
		return raster.RoundJoiner
	case BevelJoin:
		return raster.BevelJoiner
	}
	return raster.RoundJoiner
}
