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

type FillRule int

const (
	FillRuleEvenOdd FillRule = iota
	FillRuleWinding
)

type GraphicContext struct {
	PaintedImage     *image.RGBA
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
	lineWidth   float
	dash        []float
	dashOffset  float
	strokeColor image.Color
	fillColor   image.Color
	fillRule    FillRule
	cap         Cap
	join        Join
	previous    *contextStack
	fontSize    float
	fontData    FontData
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(pi *image.RGBA) *GraphicContext {
	gc := new(GraphicContext)
	gc.PaintedImage = pi
	width, height := gc.PaintedImage.Bounds().Dx(), gc.PaintedImage.Bounds().Dy()
	gc.fillRasterizer = raster.NewRasterizer(width, height)
	gc.strokeRasterizer = raster.NewRasterizer(width, height)

	gc.DPI = 92
	gc.defaultFontData = FontData{"luxi", FontFamilySans, FontStyleNormal}
	gc.freetype = freetype.NewContext()
	gc.freetype.SetDPI(gc.DPI)
	gc.freetype.SetClip(pi.Bounds())
	gc.freetype.SetDst(pi)

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

func (gc *GraphicContext) GetMatrixTransform() MatrixTransform {
	return gc.current.tr
}

func (gc *GraphicContext) SetMatrixTransform(tr MatrixTransform) {
	gc.current.tr = tr
}

func (gc *GraphicContext) ComposeMatrixTransform(tr MatrixTransform) {
	gc.current.tr = tr.Multiply(gc.current.tr)
}

func (gc *GraphicContext) Rotate(angle float) {
	gc.current.tr = NewRotationMatrix(angle).Multiply(gc.current.tr)
}

func (gc *GraphicContext) Translate(tx, ty float) {
	gc.current.tr = NewTranslationMatrix(tx, ty).Multiply(gc.current.tr)
}

func (gc *GraphicContext) Scale(sx, sy float) {
	gc.current.tr = NewScaleMatrix(sx, sy).Multiply(gc.current.tr)
}

func (gc *GraphicContext) Clear() {
	width, height := gc.PaintedImage.Bounds().Dx(), gc.PaintedImage.Bounds().Dy()
	gc.ClearRect(0, 0, width, height)
}

func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	imageColor := image.NewColorImage(gc.current.fillColor)
	draw.Draw(gc.PaintedImage, image.Rect(x1, y1, x2, y2), imageColor, image.ZP)
}

func (gc *GraphicContext) SetStrokeColor(c image.Color) {
	gc.current.strokeColor = c
}

func (gc *GraphicContext) SetFillColor(c image.Color) {
	gc.current.fillColor = c
}

func (gc *GraphicContext) SetFillRule(f FillRule) {
	gc.current.fillRule = f
}

func (gc *GraphicContext) SetLineWidth(lineWidth float) {
	gc.current.lineWidth = lineWidth
}

func (gc *GraphicContext) SetLineCap(cap Cap) {
	gc.current.cap = cap
}

func (gc *GraphicContext) SetLineJoin(join Join) {
	gc.current.join = join
}

func (gc *GraphicContext) SetLineDash(dash []float, dashOffset float) {
	gc.current.dash = dash
	gc.current.dashOffset = dashOffset
}

func (gc *GraphicContext) SetFontSize(fontSize float) {
	gc.current.fontSize = fontSize
}

func (gc *GraphicContext) GetFontSize() float {
	return gc.current.fontSize
}

func (gc *GraphicContext) SetFontData(fontData FontData) {
	gc.current.fontData = fontData
}

func (gc *GraphicContext) GetFontData() FontData {
	return gc.current.fontData
}

func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.freetype.SetDPI(dpi)
}

func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}


func (gc *GraphicContext) Save() {
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

func (gc *GraphicContext) Restore() {
	if gc.current.previous != nil {
		oldContext := gc.current
		gc.current = gc.current.previous
		oldContext.previous = nil
	}
}

func (gc *GraphicContext) DrawImage(image image.Image) {
	width := raster.Fix32(gc.PaintedImage.Bounds().Dx()* 256)
	height := raster.Fix32(gc.PaintedImage.Bounds().Dy()* 256)

	painter := raster.NewRGBAPainter(gc.PaintedImage)

	p0 := raster.Point{0,0}
	p1 := raster.Point{0,0}
	p2 := raster.Point{0,0}
	p3 := raster.Point{0,0}
	var i raster.Fix32 = 0
	for ; i < width; i+=256 {
		var j raster.Fix32 = 0
		for ; j < height; j+=256 {
			p0.X, p0.Y = i, j
			p1.X, p1.Y = p0.X + 256, p0.Y
			p2.X, p2.Y = p1.X, p0.Y + 256
			p3.X, p3.Y = p0.X, p2.Y
			
			gc.current.tr.TransformRasterPoint(&p0, &p1, &p2, &p3)
			gc.fillRasterizer.Start(p0)
			gc.fillRasterizer.Add1(p1)
			gc.fillRasterizer.Add1(p2)
			gc.fillRasterizer.Add1(p3)
			gc.fillRasterizer.Add1(p0)
			painter.SetColor(image.At(int(i>>8), int(j>>8)))
			gc.fillRasterizer.Rasterize(painter)
			gc.fillRasterizer.Clear()
		}
	}
}

func (gc *GraphicContext) BeginPath() {
	gc.current.path = new(PathStorage)
}

func (gc *GraphicContext) IsEmpty() bool{
	return gc.current.path.IsEmpty()
}

func (gc *GraphicContext) LastPoint() (float, float){
	return gc.current.path.LastPoint()
}

func (gc *GraphicContext) MoveTo(x, y float) {
	gc.current.path.MoveTo(x, y)
}

func (gc *GraphicContext) RMoveTo(dx, dy float) {
	gc.current.path.RMoveTo(dx, dy)
}

func (gc *GraphicContext) LineTo(x, y float) {
	gc.current.path.LineTo(x, y)
}

func (gc *GraphicContext) RLineTo(dx, dy float) {
	gc.current.path.RLineTo(dx, dy)
}

func (gc *GraphicContext) QuadCurveTo(cx, cy, x, y float) {
	gc.current.path.QuadCurveTo(cx, cy, x, y)
}

func (gc *GraphicContext) RQuadCurveTo(dcx, dcy, dx, dy float) {
	gc.current.path.RQuadCurveTo(dcx, dcy, dx, dy)
}

func (gc *GraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float) {
	gc.current.path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *GraphicContext) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float) {
	gc.current.path.RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy)
}

func (gc *GraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float) {
	gc.current.path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *GraphicContext) RArcTo(dcx, dcy, rx, ry, startAngle, angle float) {
	gc.current.path.RArcTo(dcx, dcy, rx, ry, startAngle, angle)
}

func (gc *GraphicContext) Close() {
	gc.current.path.Close()
}

func (gc *GraphicContext) FillString(text string) (cursor float) {
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
	x2, y2 := float(p.X)/256, float(p.Y)/256
	gc.current.tr.InverseTransform(&x2, &y2)
	width := x2 - x1
	return width
}


func (gc *GraphicContext) paint(rasterizer *raster.Rasterizer, color image.Color) {
	painter := raster.NewRGBAPainter(gc.PaintedImage)
	painter.SetColor(color)
	rasterizer.Rasterize(painter)
	rasterizer.Clear()
	gc.current.path = new(PathStorage)
}

/**** First method ****/
func (gc *GraphicContext) Stroke2(paths ...*PathStorage) {
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
func (gc *GraphicContext) Stroke(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := NewLineStroker(NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.strokeRasterizer)))
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
func (gc *GraphicContext) Fill2(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()

	pathConverter := NewPathConverter(NewVertexAdder(NewMatrixTransformAdder(gc.current.tr, gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.fillColor)
}

/**** second method ****/
func (gc *GraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()

	/**** first method ****/
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.current.tr.GetMaxAbsScaling()
	pathConverter.Convert(paths...)
	gc.paint(gc.fillRasterizer, gc.current.fillColor)
}

func (gc *GraphicContext) FillStroke2(paths ...*PathStorage) {
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
func (gc *GraphicContext) FillStroke(paths ...*PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.fillRasterizer))

	stroker := NewLineStroker(NewVertexMatrixTransform(gc.current.tr, NewVertexAdder(gc.strokeRasterizer)))
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
