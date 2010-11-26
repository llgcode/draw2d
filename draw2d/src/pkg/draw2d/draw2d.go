// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"exp/draw"
	"image"
	//"math"
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
	path        *Path
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
	gc.current.path = new(Path)
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
	gc.current.path = new(Path)
}

func (gc *GraphicContext) MoveTo(x, y float) {
	gc.current.tr.Transform(&x, &y)
	gc.current.path.MoveTo(x, y)
}

func (gc *GraphicContext) RMoveTo(dx, dy float) {
	gc.current.tr.VectorTransform(&dx, &dy)
	gc.current.path.RMoveTo(dx, dy)
}

func (gc *GraphicContext) LineTo(x, y float) {
	gc.current.tr.Transform(&x, &y)
	gc.current.path.LineTo(x, y)
}

func (gc *GraphicContext) RLineTo(dx, dy float) {
	gc.current.tr.VectorTransform(&dx, &dy)
	gc.current.path.RLineTo(dx, dy)
}

func (gc *GraphicContext) Rect(x1, y1, x2, y2 float) {
	gc.current.tr.Transform(&x1, &y1, &x2, &y2)
	gc.current.path.Rect(x1, y1, x2, y2)
}

func (gc *GraphicContext) RRect(dx1, dy1, dx2, dy2 float) {
	gc.current.tr.VectorTransform(&dx1, &dy1, &dx2, &dy2)
	gc.current.path.RRect(dx1, dy1, dx2, dy2)
}

func (gc *GraphicContext) QuadCurveTo(cx, cy, x, y float) {
	gc.current.tr.Transform(&cx, &cy, &x, &y)
	gc.current.path.QuadCurveTo(cx, cy, x, y)
}

func (gc *GraphicContext) RQuadCurveTo(dcx, dcy, dx, dy float) {
	gc.current.tr.VectorTransform(&dcx, &dcy, &dx, &dy)
	gc.current.path.RQuadCurveTo(dcx, dcy, dx, dy)
}

func (gc *GraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float) {
	gc.current.tr.Transform(&cx1, &cy1, &cx2, &cy2, &x, &y)
	gc.current.path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *GraphicContext) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float) {
	gc.current.tr.VectorTransform(&dcx1, &dcy1, &dcx2, &dcy2, &dx, &dy)
	gc.current.path.RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy)
}

func (gc *GraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float) {
	gc.current.tr.Transform(&cx, &cy)
	gc.current.tr.VectorTransform(&rx, &ry)
	gc.current.path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *GraphicContext) RArcTo(dcx, dcy, rx, ry, startAngle, angle float) {
	gc.current.tr.VectorTransform(&dcx, &dcy)
	gc.current.tr.VectorTransform(&rx, &ry)
	gc.current.path.RArcTo(dcx, dcy, rx, ry, startAngle, angle)
}

func (gc *GraphicContext) ClosePath() {
	gc.current.path.Close()
}

func (gc *GraphicContext) paint(color image.Color) {
	painter := raster.NewRGBAPainter(gc.PaintedImage)
	painter.SetColor(color)
	gc.rasterizer.Rasterize(painter)
	gc.rasterizer.Clear()
	gc.current.path = new(Path)
}

func (gc *GraphicContext) Stroke(paths ...*Path) {
	paths = append(paths, gc.current.path)
	rasterPath := tracePath(gc.current.dash, gc.current.dashOffset, paths...)
	gc.rasterizer.UseNonZeroWinding = true
	gc.rasterizer.AddStroke(*rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())
	gc.paint(gc.current.strokeColor)
}

func (gc *GraphicContext) Fill(paths ...*Path) {
	paths = append(paths, gc.current.path)
	rasterPath := tracePath(nil, 0, paths...)

	gc.rasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.rasterizer.AddPath(*rasterPath)
	gc.paint(gc.current.fillColor)
}

func (gc *GraphicContext) FillStroke(paths ...*Path) {
	paths = append(paths, gc.current.path)
	rasterPath := tracePath(nil, 0, paths...)

	gc.rasterizer.UseNonZeroWinding = gc.current.fillRule.fillRule()
	gc.rasterizer.AddPath(*rasterPath)
	gc.paint(gc.current.fillColor)
	if gc.current.dash != nil {
		rasterPath = tracePath(gc.current.dash, gc.current.dashOffset, paths...)
	}
	gc.rasterizer.UseNonZeroWinding = true
	gc.rasterizer.AddStroke(*rasterPath, raster.Fix32(gc.current.lineWidth*256), gc.current.cap.capper(), gc.current.join.joiner())
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


type PathAdapter struct {
	path           *raster.Path
	x, y, distance float
	dash           []float
	currentDash    int
	dashOffset     float
}

func tracePath(dash []float, dashOffset float, paths ...*Path) *raster.Path {
	var adapter PathAdapter
	if dash != nil && len(dash) > 0 {
		adapter.dash = dash
	} else {
		adapter.dash = nil
	}
	adapter.currentDash = 0
	adapter.dashOffset = dashOffset
	adapter.path = new(raster.Path)
	for _, path := range paths {
		path.TraceLine(&adapter)
	}
	return adapter.path
}

func floatToPoint(x, y float) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}

func (p *PathAdapter) MoveTo(x, y float) {
	p.path.Start(floatToPoint(x, y))
	p.x, p.y = x, y
	p.distance = p.dashOffset
	p.currentDash = 0
}

func (p *PathAdapter) LineTo(x, y float) {
	if p.dash != nil {
		rest := p.dash[p.currentDash] - p.distance
		for rest < 0 {
			p.distance = p.distance - p.dash[p.currentDash]
			p.currentDash = (p.currentDash + 1) % len(p.dash)
			rest = p.dash[p.currentDash] - p.distance
		}
		d := distance(p.x, p.y, x, y)
		for d >= rest {
			k := rest / d
			lx := p.x + k*(x-p.x)
			ly := p.y + k*(y-p.y)
			if p.currentDash%2 == 0 {
				// line
				p.path.Add1(floatToPoint(lx, ly))
			} else {
				// gap
				p.path.Start(floatToPoint(lx, ly))
			}
			d = d - rest
			p.x, p.y = lx, ly
			p.currentDash = (p.currentDash + 1) % len(p.dash)
			rest = p.dash[p.currentDash]
		}
		p.distance = d
		if p.currentDash%2 == 0 {
			p.path.Add1(floatToPoint(x, y))
		} else {
			p.path.Start(floatToPoint(x, y))
		}
		if p.distance >= p.dash[p.currentDash] {
			p.distance = p.distance - p.dash[p.currentDash]
			p.currentDash = (p.currentDash + 1) % len(p.dash)
		}
	} else {
		p.path.Add1(floatToPoint(x, y))
	}
	p.x, p.y = x, y
}
