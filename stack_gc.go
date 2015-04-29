// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"github.com/llgcode/draw2d/path"
	"image"
	"image/color"

	"code.google.com/p/freetype-go/freetype/truetype"
)

type StackGraphicContext struct {
	Current *ContextStack
}

type ContextStack struct {
	Tr          MatrixTransform
	Path        *path.Path
	LineWidth   float64
	Dash        []float64
	DashOffset  float64
	StrokeColor color.Color
	FillColor   color.Color
	FillRule    FillRule
	Cap         path.Cap
	Join        path.Join
	FontSize    float64
	FontData    FontData

	font *truetype.Font
	// fontSize and dpi are used to calculate scale. scale is the number of
	// 26.6 fixed point units in 1 em.
	scale int32

	previous *ContextStack
}

/**
 * Create a new Graphic context from an image
 */
func NewStackGraphicContext() *StackGraphicContext {
	gc := &StackGraphicContext{}
	gc.Current = new(ContextStack)
	gc.Current.Tr = NewIdentityMatrix()
	gc.Current.Path = new(path.Path)
	gc.Current.LineWidth = 1.0
	gc.Current.StrokeColor = image.Black
	gc.Current.FillColor = image.White
	gc.Current.Cap = path.RoundCap
	gc.Current.FillRule = FillRuleEvenOdd
	gc.Current.Join = path.RoundJoin
	gc.Current.FontSize = 10
	gc.Current.FontData = defaultFontData
	return gc
}

func (gc *StackGraphicContext) GetMatrixTransform() MatrixTransform {
	return gc.Current.Tr
}

func (gc *StackGraphicContext) SetMatrixTransform(Tr MatrixTransform) {
	gc.Current.Tr = Tr
}

func (gc *StackGraphicContext) ComposeMatrixTransform(Tr MatrixTransform) {
	gc.Current.Tr = Tr.Multiply(gc.Current.Tr)
}

func (gc *StackGraphicContext) Rotate(angle float64) {
	gc.Current.Tr = NewRotationMatrix(angle).Multiply(gc.Current.Tr)
}

func (gc *StackGraphicContext) Translate(tx, ty float64) {
	gc.Current.Tr = NewTranslationMatrix(tx, ty).Multiply(gc.Current.Tr)
}

func (gc *StackGraphicContext) Scale(sx, sy float64) {
	gc.Current.Tr = NewScaleMatrix(sx, sy).Multiply(gc.Current.Tr)
}

func (gc *StackGraphicContext) SetStrokeColor(c color.Color) {
	gc.Current.StrokeColor = c
}

func (gc *StackGraphicContext) SetFillColor(c color.Color) {
	gc.Current.FillColor = c
}

func (gc *StackGraphicContext) SetFillRule(f FillRule) {
	gc.Current.FillRule = f
}

func (gc *StackGraphicContext) SetLineWidth(LineWidth float64) {
	gc.Current.LineWidth = LineWidth
}

func (gc *StackGraphicContext) SetLineCap(cap path.Cap) {
	gc.Current.Cap = cap
}

func (gc *StackGraphicContext) SetLineJoin(join path.Join) {
	gc.Current.Join = join
}

func (gc *StackGraphicContext) SetLineDash(Dash []float64, DashOffset float64) {
	gc.Current.Dash = Dash
	gc.Current.DashOffset = DashOffset
}

func (gc *StackGraphicContext) SetFontSize(FontSize float64) {
	gc.Current.FontSize = FontSize
}

func (gc *StackGraphicContext) GetFontSize() float64 {
	return gc.Current.FontSize
}

func (gc *StackGraphicContext) SetFontData(FontData FontData) {
	gc.Current.FontData = FontData
}

func (gc *StackGraphicContext) GetFontData() FontData {
	return gc.Current.FontData
}

func (gc *StackGraphicContext) BeginPath() {
	gc.Current.Path.Clear()
}

func (gc *StackGraphicContext) IsEmpty() bool {
	return gc.Current.Path.IsEmpty()
}

func (gc *StackGraphicContext) LastPoint() (float64, float64) {
	return gc.Current.Path.LastPoint()
}

func (gc *StackGraphicContext) MoveTo(x, y float64) {
	gc.Current.Path.MoveTo(x, y)
}

func (gc *StackGraphicContext) LineTo(x, y float64) {
	gc.Current.Path.LineTo(x, y)
}

func (gc *StackGraphicContext) QuadCurveTo(cx, cy, x, y float64) {
	gc.Current.Path.QuadCurveTo(cx, cy, x, y)
}

func (gc *StackGraphicContext) CubicCurveTo(cx1, cy1, cx2, cy2, x, y float64) {
	gc.Current.Path.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
}

func (gc *StackGraphicContext) ArcTo(cx, cy, rx, ry, startAngle, angle float64) {
	gc.Current.Path.ArcTo(cx, cy, rx, ry, startAngle, angle)
}

func (gc *StackGraphicContext) Close() {
	gc.Current.Path.Close()
}

func (gc *StackGraphicContext) Save() {
	context := new(ContextStack)
	context.FontSize = gc.Current.FontSize
	context.FontData = gc.Current.FontData
	context.LineWidth = gc.Current.LineWidth
	context.StrokeColor = gc.Current.StrokeColor
	context.FillColor = gc.Current.FillColor
	context.FillRule = gc.Current.FillRule
	context.Dash = gc.Current.Dash
	context.DashOffset = gc.Current.DashOffset
	context.Cap = gc.Current.Cap
	context.Join = gc.Current.Join
	context.Path = gc.Current.Path.Copy()
	context.font = gc.Current.font
	context.scale = gc.Current.scale
	copy(context.Tr[:], gc.Current.Tr[:])
	context.previous = gc.Current
	gc.Current = context
}

func (gc *StackGraphicContext) Restore() {
	if gc.Current.previous != nil {
		oldContext := gc.Current
		gc.Current = gc.Current.previous
		oldContext.previous = nil
	}
}
