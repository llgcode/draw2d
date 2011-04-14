// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"exp/draw"
	"image"
	"log"
	"freetype-go.googlecode.com/hg/freetype"
	"freetype-go.googlecode.com/hg/freetype/raster"
)

type Painter interface {
	raster.Painter
	SetColor(color image.Color)
}

var (
	defaultFontData = FontData{"luxi", FontFamilySans, FontStyleNormal}
)

type ImageGraphicContext struct {
	*StackGraphicContext
	img     draw.Image
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
	case *image.NRGBA:
		painter = NewNRGBAPainter(selectImage)
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
	imageColor := image.NewColorImage(gc.current.FillColor)
	draw.Draw(gc.img, image.Rect(x1, y1, x2, y2), imageColor, image.ZP)
}

func (gc *ImageGraphicContext) DrawImage(img image.Image) {
	DrawImage(img, gc.img, gc.current.Tr, draw.Over, BilinearFilter)
}

func (gc *ImageGraphicContext) FillString(text string) (cursor float64) {
	gc.freetype.SetSrc(image.NewColorImage(gc.current.StrokeColor))
	// Draw the text.
	x, y := gc.current.Path.LastPoint()
	gc.current.Tr.Transform(&x, &y)
	x0, fontSize := 0.0, gc.current.FontSize
	gc.current.Tr.VectorTransform(&x0, &fontSize)
	font := GetFont(gc.current.FontData)
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
	x1, _ := gc.current.Path.LastPoint()
	x2, y2 := float64(p.X)/256, float64(p.Y)/256
	gc.current.Tr.InverseTransform(&x2, &y2)
	width := x2 - x1
	return width
}


func (gc *ImageGraphicContext) paint(rasterizer *raster.Rasterizer, color image.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.current.Path = new(PathStorage)
}

/**** First method ****/
func (gc *ImageGraphicContext) Stroke2(paths ...*PathStorage) {
	paths = append(paths, gc.current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	rasterPath := new(raster.Path)

	var pathConverter *PathConverter
	if gc.current.Dash != nil && len(gc.current.Dash) > 0 {
		dasher := NewDashConverter(gc.current.Dash, gc.current.DashOffset, NewVertexAdder(rasterPath))
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(NewVertexAdder(rasterPath))
	}

	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	mta := NewMatrixTransformAdder(gc.current.Tr, gc.strokeRasterizer)
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.LineWidth*256), gc.current.Cap.capper(), gc.current.Join.joiner())

	gc.paint(gc.strokeRasterizer, gc.current.StrokeColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Stroke(paths ...*PathStorage) {
	paths = append(paths, gc.current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := NewLineStroker(gc.current.Cap, gc.current.Join, NewVertexMatrixTransform(gc.current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.current.LineWidth / 2
	var pathConverter *PathConverter
	if gc.current.Dash != nil && len(gc.current.Dash) > 0 {
		dasher := NewDashConverter(gc.current.Dash, gc.current.DashOffset, stroker)
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(stroker)
	}
	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.current.StrokeColor)
}

/**** first method ****/
func (gc *ImageGraphicContext) Fill2(paths ...*PathStorage) {
	paths = append(paths, gc.current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.FillRule.fillRule()

	pathConverter := NewPathConverter(NewVertexAdder(NewMatrixTransformAdder(gc.current.Tr, gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.FillColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.FillRule.fillRule()

	/**** first method ****/
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.current.Tr, NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.FillColor)
}

func (gc *ImageGraphicContext) FillStroke2(paths ...*PathStorage) {
	paths = append(paths, gc.current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.FillRule.fillRule()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.current.Tr, NewVertexAdder(gc.fillRasterizer))
	rasterPath := new(raster.Path)
	stroker := NewVertexAdder(rasterPath)

	demux := NewDemuxConverter(filler, stroker)

	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	mta := NewMatrixTransformAdder(gc.current.Tr, gc.strokeRasterizer)
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.LineWidth*256), gc.current.Cap.capper(), gc.current.Join.joiner())

	gc.paint(gc.fillRasterizer, gc.current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.current.StrokeColor)
}

/* second method */
func (gc *ImageGraphicContext) FillStroke(paths ...*PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.current.FillRule.fillRule()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.current.Tr, NewVertexAdder(gc.fillRasterizer))

	stroker := NewLineStroker(gc.current.Cap, gc.current.Join, NewVertexMatrixTransform(gc.current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.current.LineWidth / 2

	demux := NewDemuxConverter(filler, stroker)
	paths = append(paths, gc.current.Path)
	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.current.Tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.current.StrokeColor)
}


func (f FillRule) fillRule() bool {
	switch f {
	case FillRuleEvenOdd:
		return false
	case FillRuleWinding:
		return true
	}
	return false
}

func (c Cap) capper() raster.Capper {
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

func (j Join) joiner() raster.Joiner {
	switch j {
	case RoundJoin:
		return raster.RoundJoiner
	case BevelJoin:
		return raster.BevelJoiner
	}
	return raster.RoundJoiner
}
