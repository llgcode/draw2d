// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"image"
)

type StackGraphicContext struct {
	current *ContextStack
}

type ContextStack struct {
	Tr          MatrixTransform
	Path        *PathStorage
	LineWidth   float64
	Dash        []float64
	DashOffset  float64
	StrokeColor image.Color
	FillColor   image.Color
	FillRule    FillRule
	Cap         Cap
	Join        Join
	FontSize    float64
	FontData    FontData
	previous    *ContextStack
}


/**
 * Create a new Graphic context from an image
 */
func NewStackGraphicContext() *StackGraphicContext {
	gc := &StackGraphicContext{}
	gc.current = new(ContextStack)
	gc.current.Tr = NewIdentityMatrix()
	gc.current.Path = new(PathStorage)
	gc.current.LineWidth = 1.0
	gc.current.StrokeColor = image.Black
	gc.current.FillColor = image.White
	gc.current.Cap = RoundCap
	gc.current.FillRule = FillRuleEvenOdd
	gc.current.Join = RoundJoin
	gc.current.FontSize = 10
	gc.current.FontData = defaultFontData
	return gc
}


func (gc *StackGraphicContext) GetMatrixTransform() MatrixTransform {
	return gc.current.Tr
}

func (gc *StackGraphicContext) SetMatrixTransform(Tr MatrixTransform) {
	gc.current.Tr = Tr
}

func (gc *StackGraphicContext) ComposeMatrixTransform(Tr MatrixTransform) {
	gc.current.Tr = Tr.Multiply(gc.current.Tr)
}

func (gc *StackGraphicContext) Rotate(angle float64) {
	gc.current.Tr = NewRotationMatrix(angle).Multiply(gc.current.Tr)
}

func (gc *StackGraphicContext) Translate(tx, ty float64) {
	gc.current.Tr = NewTranslationMatrix(tx, ty).Multiply(gc.current.Tr)
}

func (gc *StackGraphicContext) Scale(sx, sy float64) {
	gc.current.Tr = NewScaleMatrix(sx, sy).Multiply(gc.current.Tr)
}

func (gc *StackGraphicContext) SetStrokeColor(c image.Color) {
	gc.current.StrokeColor = c
}

func (gc *StackGraphicContext) SetFillColor(c image.Color) {
	gc.current.FillColor = c
}

func (gc *StackGraphicContext) SetFillRule(f FillRule) {
	gc.current.FillRule = f
}

func (gc *StackGraphicContext) SetLineWidth(LineWidth float64) {
	gc.current.LineWidth = LineWidth
}

func (gc *StackGraphicContext) SetLineCap(Cap Cap) {
	gc.current.Cap = Cap
}

func (gc *StackGraphicContext) SetLineJoin(Join Join) {
	gc.current.Join = Join
}

func (gc *StackGraphicContext) SetLineDash(Dash []float64, DashOffset float64) {
	gc.current.Dash = Dash
	gc.current.DashOffset = DashOffset
}

func (gc *StackGraphicContext) SetFontSize(FontSize float64) {
	gc.current.FontSize = FontSize
}

func (gc *StackGraphicContext) GetFontSize() float64 {
	return gc.current.FontSize
}

func (gc *StackGraphicContext) SetFontData(FontData FontData) {
	gc.current.FontData = FontData
}

func (gc *StackGraphicContext) GetFontData() FontData {
	return gc.current.FontData
}

func (gc *StackGraphicContext) BeginPath() {
	gc.current.Path = new(PathStorage)
}

func (gc *StackGraphicContext) IsEmpty() bool {
	return gc.current.Path.IsEmpty()
}

func (gc *StackGraphicContext) LastPoint() (float64, float64) {
	return gc.current.Path.LastPoint()
}

func (gc *StackGraphicContext) MoveTo(x, y float64) {
	gc.current.Path.MoveTo(x, y)
}

func (gc *StackGraphicContext) RMoveTo(dx, dy float64) {
	gc.current.Path.RMoveTo(dx, dy)
}

func (gc *StackGraphicContext) LineTo(x, y float64) {
	gc.current.Path.LineTo(x, y)
}

func (gc *StackGraphicContext) RLineTo(dx, dy float64) {
	gc.current.Path.RLineTo(dx, dy)
}

func (gc *StackGraphicContext) QuadCurveTo(cx, cy, x, y float64) {
	gc.current.Path.QuadCurveTo(cx, cy, x, y)
}

func (gc *StackGraphicContext) RQuadCurveTo(dcx, dcy, dx, dy float64) {
	gc.current.Path.RQuadCurveTo(dcx, dcy, dx, dy)
}

func (gc *StackGraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64) {
	gc.current.Path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *StackGraphicContext) RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy float64) {
	gc.current.Path.RCubicCurveTo(dcx1, dcy1, dcx2, dcy2, dx, dy)
}

func (gc *StackGraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float64) {
	gc.current.Path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *StackGraphicContext) RArcTo(dcx, dcy, rx, ry, startAngle, angle float64) {
	gc.current.Path.RArcTo(dcx, dcy, rx, ry, startAngle, angle)
}

func (gc *StackGraphicContext) Close() {
	gc.current.Path.Close()
}

func (gc *StackGraphicContext) Save() {
	context := new(ContextStack)
	context.FontSize = gc.current.FontSize
	context.FontData = gc.current.FontData
	context.LineWidth = gc.current.LineWidth
	context.StrokeColor = gc.current.StrokeColor
	context.FillColor = gc.current.FillColor
	context.FillRule = gc.current.FillRule
	context.Dash = gc.current.Dash
	context.DashOffset = gc.current.DashOffset
	context.Cap = gc.current.Cap
	context.Join = gc.current.Join
	context.Path = gc.current.Path.Copy()
	copy(context.Tr[:], gc.current.Tr[:])
	context.previous = gc.current
	gc.current = context
}

func (gc *StackGraphicContext) Restore() {
	if gc.current.previous != nil {
		oldContext := gc.current
		gc.current = gc.current.previous
		oldContext.previous = nil
	}
}
