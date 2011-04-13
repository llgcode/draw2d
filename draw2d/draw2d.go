// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"exp/draw"
	"image"
	"log"
	"math"
	"freetype-go.googlecode.com/hg/freetype"
	"freetype-go.googlecode.com/hg/freetype/raster"
)

type Painter interface {
	raster.Painter
	SetColor(color image.Color)
}

type ImageGraphicContext struct {
	img              draw.Image
	painter          Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	freetype         *freetype.Context
	defaultFontData  FontData
	DPI              int
	current          *contextStack
}

type contextStack struct {
	tr          MatrixTransform
	path        *PathStorage
	lineWidth   float64
	dash        []float64
	dashOffset  float64
	strokeColor image.Color
	fillColor   image.Color
	fillRule    FillRule
	cap         Cap
	join        Join
	previous    *contextStack
	fontSize    float64
	fontData    FontData
}

/**
 * Create a new Graphic context from an image
 */
func NewImageGraphicContext(img draw.Image) *ImageGraphicContext {
	gc := new(ImageGraphicContext)
	gc.img = img
	switch selectImage := img.(type) {
	case *image.RGBA:
		gc.painter = raster.NewRGBAPainter(selectImage)
	case *image.NRGBA:
		gc.painter = NewNRGBAPainter(selectImage)
	default:
		panic("Image type not supported")
	}

	width, height := gc.img.Bounds().Dx(), gc.img.Bounds().Dy()
	gc.fillRasterizer = raster.NewRasterizer(width, height)
	gc.strokeRasterizer = raster.NewRasterizer(width, height)

	gc.DPI = 92
	gc.defaultFontData = FontData{"luxi", FontFamilySans, FontStyleNormal}
	gc.freetype = freetype.NewContext()
	gc.freetype.SetDPI(gc.DPI)
	gc.freetype.SetClip(img.Bounds())
	gc.freetype.SetDst(img)

	gc.current = new(contextStack)

	gc.current.tr = NewIdentityMatrix()
	gc.current.path = new(PathStorage)
	gc.current.lineWidth = 1.0
	gc.current.strokeColor = image.Black
	gc.current.fillColor = image.White
	gc.current.cap = RoundCap
	gc.current.fillRule = FillRuleEvenOdd
	gc.current.join = RoundJoin
	gc.current.fontSize = 10
	gc.current.fontData = gc.defaultFontData

	return gc
}

func (gc *ImageGraphicContext) GetMatrixTransform() MatrixTransform {
	return gc.current.tr
}

func (gc *ImageGraphicContext) SetMatrixTransform(tr MatrixTransform) {
	gc.current.tr = tr
}

func (gc *ImageGraphicContext) ComposeMatrixTransform(tr MatrixTransform) {
	gc.current.tr = tr.Multiply(gc.current.tr)
}

func (gc *ImageGraphicContext) Rotate(angle float64) {
	gc.current.tr = NewRotationMatrix(angle).Multiply(gc.current.tr)
}

func (gc *ImageGraphicContext) Translate(tx, ty float64) {
	gc.current.tr = NewTranslationMatrix(tx, ty).Multiply(gc.current.tr)
}

func (gc *ImageGraphicContext) Scale(sx, sy float64) {
	gc.current.tr = NewScaleMatrix(sx, sy).Multiply(gc.current.tr)
}

func (gc *ImageGraphicContext) Clear() {
	width, height := gc.img.Bounds().Dx(), gc.img.Bounds().Dy()
	gc.ClearRect(0, 0, width, height)
}

func (gc *ImageGraphicContext) ClearRect(x1, y1, x2, y2 int) {
	imageColor := image.NewColorImage(gc.current.fillColor)
	draw.Draw(gc.img, image.Rect(x1, y1, x2, y2), imageColor, image.ZP)
}

func (gc *ImageGraphicContext) SetStrokeColor(c image.Color) {
	gc.current.strokeColor = c
}

func (gc *ImageGraphicContext) SetFillColor(c image.Color) {
	gc.current.fillColor = c
}

func (gc *ImageGraphicContext) SetFillRule(f FillRule) {
	gc.current.fillRule = f
}

func (gc *ImageGraphicContext) SetLineWidth(lineWidth float64) {
	gc.current.lineWidth = lineWidth
}

func (gc *ImageGraphicContext) SetLineCap(cap Cap) {
	gc.current.cap = cap
}

func (gc *ImageGraphicContext) SetLineJoin(join Join) {
	gc.current.join = join
}

func (gc *ImageGraphicContext) SetLineDash(dash []float64, dashOffset float64) {
	gc.current.dash = dash
	gc.current.dashOffset = dashOffset
}

func (gc *ImageGraphicContext) SetFontSize(fontSize float64) {
	gc.current.fontSize = fontSize
}

func (gc *ImageGraphicContext) GetFontSize() float64 {
	return gc.current.fontSize
}

func (gc *ImageGraphicContext) SetFontData(fontData FontData) {
	gc.current.fontData = fontData
}

func (gc *ImageGraphicContext) GetFontData() FontData {
	return gc.current.fontData
}

func (gc *ImageGraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.freetype.SetDPI(dpi)
}

func (gc *ImageGraphicContext) GetDPI() int {
	return gc.DPI
}


func (gc *ImageGraphicContext) Save() {
	context := new(contextStack)
	context.fontSize = gc.current.fontSize
	context.fontData = gc.current.fontData
	context.lineWidth = gc.current.lineWidth
	context.strokeColor = gc.current.strokeColor
	context.fillColor = gc.current.fillColor
	context.fillRule = gc.current.fillRule
	context.dash = gc.current.dash
	context.dashOffset = gc.current.dashOffset
	context.cap = gc.current.cap
	context.join = gc.current.join
	context.path = gc.current.path.Copy()
	copy(context.tr[:], gc.current.tr[:])
	context.previous = gc.current
	gc.current = context
}

func (gc *ImageGraphicContext) Restore() {
	if gc.current.previous != nil {
		oldContext := gc.current
		gc.current = gc.current.previous
		oldContext.previous = nil
	}
}
//see http://pippin.gimp.org/image_processing/chap_resampling.html
func getColorLinear(img image.Image, x, y float64) image.Color {
	return img.At(int(x), int(y))
}

func getColorBilinear(img image.Image, x, y float64) image.Color {
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	dx := x - x0
	dy := y - y0

	color0 := img.At(int(x0), int(y0))
	color1 := img.At(int(x0+1), int(y0))
	color2 := img.At(int(x0+1), int(y0+1))
	color3 := img.At(int(x0), int(y0+1))

	return lerp(lerp(color0, color1, dx), lerp(color3, color2, dx), dy)
}
/**
-- LERP
-- /lerp/, vi.,n.
--
-- Quasi-acronym for Linear Interpolation, used as a verb or noun for
-- the operation. "Bresenham's algorithm lerps incrementally between the
-- two endpoints of the line." (From Jargon File (4.4.4, 14 Aug 2003)
*/
func lerp(c1, c2 image.Color, ratio float64) image.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r := int(float64(r1)*(1-ratio) + float64(r2)*ratio)
	g := int(float64(g1)*(1-ratio) + float64(g2)*ratio)
	b := int(float64(b1)*(1-ratio) + float64(b2)*ratio)
	a := int(float64(a1)*(1-ratio) + float64(a2)*ratio)
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}


func getColorCubicRow(img image.Image, x, y, offset float64) image.Color {
	c0 := img.At(int(x), int(y))
	c1 := img.At(int(x+1), int(y))
	c2 := img.At(int(x+2), int(y))
	c3 := img.At(int(x+3), int(y))
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r3, g3, b3, a3 := c3.RGBA()
	r, g, b, a := cubic(offset,float64(r0),float64(r1),float64(r2),float64(r3)), cubic(offset,float64(g0),float64(g1),float64(g2),float64(g3)), cubic(offset,float64(b0),float64(b1),float64(b2),float64(b3)), cubic(offset,float64(a0),float64(a1),float64(a2),float64(a3))
  return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func getColorBicubic(img image.Image, x, y float64) image.Color {
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	dx := x - x0
	dy := y - y0

	c0 := getColorCubicRow(img, x0-1, y0-1, dx)
	c1 := getColorCubicRow(img, x0-1, y0, dx)
	c2 := getColorCubicRow(img, x0-1, y0+1, dx)
	c3 := getColorCubicRow(img, x0-1, y0+2, dx)
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r3, g3, b3, a3 := c3.RGBA()
	r, g, b, a := cubic(dy,float64(r0),float64(r1),float64(r2),float64(r3)), cubic(dy,float64(g0),float64(g1),float64(g2),float64(g3)), cubic(dy,float64(b0),float64(b1),float64(b2),float64(b3)), cubic(dy,float64(a0),float64(a1),float64(a2),float64(a3))
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func cubic(offset,v0,v1,v2,v3 float64) uint32{
  // offset is the offset of the sampled value between v1 and v2
   return   uint32((((( -7 * v0 + 21 * v1 - 21 * v2 + 7 * v3 ) * offset +
               ( 15 * v0 - 36 * v1 + 27 * v2 - 6 * v3 ) ) * offset +
               ( -9 * v0 + 9 * v2 ) ) * offset + (v0 + 16 * v1 + v2) ) / 18.0);
}

func compose(c1, c2 image.Color) image.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	ia := M - a2
	r := ((r1 * ia) / M) + r2
	g := ((g1 * ia) / M) + g2
	b := ((b1 * ia) / M) + b2
	a := ((a1 * ia) / M) + a2
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func (gc *ImageGraphicContext) DrawImage(img image.Image) {
	width := float64(gc.img.Bounds().Dx())
	height := float64(gc.img.Bounds().Dy())
	gc.current.tr.Transform(&width, &height)
	var x, y, u, v float64
	for x = 0; x < width; x++ {
		for y = 0; y < height; y++ {
			u = x
			v = y
			gc.current.tr.InverseTransform(&u, &v)
			gc.img.Set(int(x), int(y), compose(gc.img.At(int(x), int(y)), getColorLinear(img, u, v)))
		}
	}
}

func (gc *ImageGraphicContext) BeginPath() {
	gc.current.path = new(PathStorage)
}

func (gc *ImageGraphicContext) IsEmpty() bool {
	return gc.current.path.IsEmpty()
}

func (gc *ImageGraphicContext) LastPoint() (float64, float64) {
	return gc.current.path.LastPoint()
}

func (gc *ImageGraphicContext) MoveTo(x, y float64) {
	gc.current.path.MoveTo(x, y)
}

func (gc *ImageGraphicContext) RMoveTo(dx, dy float64) {
	gc.current.path.RMoveTo(dx, dy)
}

func (gc *ImageGraphicContext) LineTo(x, y float64) {
	gc.current.path.LineTo(x, y)
}

func (gc *ImageGraphicContext) RLineTo(dx, dy float64) {
	gc.current.path.RLineTo(dx, dy)
}

func (gc *ImageGraphicContext) QuadCurveTo(cx, cy, x, y float64) {
	gc.current.path.QuadCurveTo(cx, cy, x, y)
}

func (gc *ImageGraphicContext) RQuadCurveTo(dcx, dcy, dx, dy float64) {
	gc.current.path.RQuadCurveTo(dcx, dcy, dx, dy)
}

func (gc *ImageGraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64) {
	gc.current.path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *ImageGraphicContext) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float64) {
	gc.current.path.RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy)
}

func (gc *ImageGraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float64) {
	gc.current.path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *ImageGraphicContext) RArcTo(dcx, dcy, rx, ry, startAngle, angle float64) {
	gc.current.path.RArcTo(dcx, dcy, rx, ry, startAngle, angle)
}

func (gc *ImageGraphicContext) Close() {
	gc.current.path.Close()
}

func (gc *ImageGraphicContext) FillString(text string) (cursor float64) {
	gc.freetype.SetSrc(image.NewColorImage(gc.current.strokeColor))
	// Draw the text.
	x, y := gc.current.path.LastPoint()
	gc.current.tr.Transform(&x, &y)
	x0, fontSize := 0.0, gc.current.fontSize
	gc.current.tr.VectorTransform(&x0, &fontSize)
	font := GetFont(gc.current.fontData)
	if font == nil {
		font = GetFont(gc.defaultFontData)
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
	x1, _ := gc.current.path.LastPoint()
	x2, y2 := float64(p.X)/256, float64(p.Y)/256
	gc.current.tr.InverseTransform(&x2, &y2)
	width := x2 - x1
	return width
}


func (gc *ImageGraphicContext) paint(rasterizer *raster.Rasterizer, color image.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.current.path = new(PathStorage)
}

/**** First method ****/
func (gc *ImageGraphicContext) Stroke2(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	rasterPath := new(raster.Path)

	var pathConverter *PathConverter
	if gc.current.dash != nil && len(gc.current.dash) > 0 {
		dasher := NewDashConverter(gc.current.dash, gc.current.dashOffset, NewVertexAdder(rasterPath))
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(NewVertexAdder(rasterPath))
	}

	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	mta := NewMatrixTransformAdder(gc.current.tr, gc.strokeRasterizer)
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())

	gc.paint(gc.strokeRasterizer, gc.current.strokeColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Stroke(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := NewLineStroker(gc.current.cap, gc.current.join, NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.current.lineWidth / 2
	var pathConverter *PathConverter
	if gc.current.dash != nil && len(gc.current.dash) > 0 {
		dasher := NewDashConverter(gc.current.dash, gc.current.dashOffset, stroker)
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(stroker)
	}
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.current.strokeColor)
}

/**** first method ****/
func (gc *ImageGraphicContext) Fill2(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()

	pathConverter := NewPathConverter(NewVertexAdder(NewMatrixTransformAdder(gc.current.tr, gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.fillColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()

	/**** first method ****/
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.fillColor)
}

func (gc *ImageGraphicContext) FillStroke2(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.fillRasterizer))
	rasterPath := new(raster.Path)
	stroker := NewVertexAdder(rasterPath)

	demux := NewDemuxConverter(filler, stroker)

	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	mta := NewMatrixTransformAdder(gc.current.tr, gc.strokeRasterizer)
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())

	gc.paint(gc.fillRasterizer, gc.current.fillColor)
	gc.paint(gc.strokeRasterizer, gc.current.strokeColor)
}

/* second method */
func (gc *ImageGraphicContext) FillStroke(paths ...*PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.fillRasterizer))

	stroker := NewLineStroker(gc.current.cap, gc.current.join, NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.current.lineWidth / 2

	demux := NewDemuxConverter(filler, stroker)
	paths = append(paths, gc.current.path)
	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.current.fillColor)
	gc.paint(gc.strokeRasterizer, gc.current.strokeColor)
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
