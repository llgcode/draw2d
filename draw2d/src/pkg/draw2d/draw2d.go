// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"exp/draw"
	"image"
	"freetype-go.googlecode.com/hg/freetype/raster"
)

type FillRule int

const (
	FillRuleEvenOdd FillRule = iota
	FillRuleWinding
)

type Cap int

const (
	RoundCap Cap = iota
	ButtCap
	SquareCap
)

type Join int

const (
	BevelJoin Join = iota
	RoundJoin
	MiterJoin
)

type GraphicContext struct {
	PaintedImage *image.RGBA
	rasterizer   *raster.Rasterizer
	current      *contextStack
}

type contextStack struct {
	tr			MatrixTransform
	path        *PathStorage
	lineWidth   float
	dash        []float
	dashOffset  float
	strokeColor image.Color
	fillColor   image.Color
	fillRule    FillRule
	cap         Cap
	join        Join
	previous *contextStack
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(pi *image.RGBA) *GraphicContext {
	gc := new(GraphicContext)
	gc.PaintedImage = pi
	width, height := gc.PaintedImage.Bounds().Dx(), gc.PaintedImage.Bounds().Dy()
	gc.rasterizer = raster.NewRasterizer(width, height)

	gc.current = new(contextStack)
	
	gc.current.tr = NewIdentityMatrix()
	gc.current.path = new(PathStorage)
	gc.current.lineWidth = 1.0
	gc.current.strokeColor = image.Black
	gc.current.fillColor = image.White
	gc.current.cap = RoundCap
	gc.current.fillRule = FillRuleEvenOdd
	gc.current.join = RoundJoin
	return gc
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

func (gc *GraphicContext) Save() {
	context := new(contextStack)
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

func (gc *GraphicContext) BeginPath() {
	gc.current.path = new(PathStorage)
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

func (gc *GraphicContext) paint(color image.Color) {
	painter := raster.NewRGBAPainter(gc.PaintedImage)
	painter.SetColor(color)
	gc.rasterizer.Rasterize(painter)
	gc.rasterizer.Clear()
	gc.current.path = new(PathStorage)
}

func (gc *GraphicContext) Stroke(paths ...*PathStorage) {	
	paths = append(paths, gc.current.path)
	gc.rasterizer.UseNonZeroWinding = true
	rasterPath := new(raster.Path)
	if(gc.current.dash == nil) {
		tracePath(gc.current.tr.GetMaxAbsScaling(), rasterPath, paths...)
	} else {
		traceDashPath(gc.current.dash, gc.current.dashOffset, gc.current.tr.GetMaxAbsScaling(), rasterPath, paths...)
	}
	mta := NewMatrixTransformAdder(gc.current.tr, gc.rasterizer)
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())
	gc.paint(gc.current.strokeColor)
}

func (gc *GraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	gc.rasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	mta := NewMatrixTransformAdder(gc.current.tr, gc.rasterizer)
	tracePath(gc.current.tr.GetMaxAbsScaling(), mta,  paths...)
	gc.paint(gc.current.fillColor)
}

func (gc *GraphicContext) FillStroke(paths ...*PathStorage) {
	paths = append(paths, gc.current.path)
	mta := NewMatrixTransformAdder(gc.current.tr, gc.rasterizer)
	tracePath(gc.current.tr.GetMaxAbsScaling(), mta, paths...)

	gc.rasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.paint(gc.current.fillColor)
	
	gc.rasterizer.UseNonZeroWinding = true
	rasterPath := new(raster.Path)
	if(gc.current.dash == nil) {
		tracePath(gc.current.tr.GetMaxAbsScaling(), rasterPath, paths...)
	} else {
		traceDashPath(gc.current.dash, gc.current.dashOffset, gc.current.tr.GetMaxAbsScaling(), rasterPath, paths...)
	}
	raster.Stroke(mta, *rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())
	gc.paint(gc.current.strokeColor)
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

