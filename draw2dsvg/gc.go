// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

package draw2dsvg

import (
	"image"
	"bytes"
	"image/color"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
)

const (
)

var (
)

type SVG bytes.Buffer

func NewSvg() *SVG {
	return &SVG{}
}

// GraphicContext implements the draw2d.GraphicContext interface
// It provides draw2d with a svg backend
type GraphicContext struct {
	*draw2dbase.StackGraphicContext
	svg *SVG
}

func NewGraphicContext(svg *SVG) *GraphicContext {
	gc := &GraphicContext{draw2dbase.NewStackGraphicContext(), svg}
	return gc
}

// TODO implement all following methods

// BeginPath creates a new path
func (gc *GraphicContext) BeginPath() {

}
// GetPath copies the current path, then returns it
func (gc *GraphicContext) GetPath() draw2d.Path {
	return draw2d.Path{}
}
// GetMatrixTransform returns the current transformation matrix
func (gc *GraphicContext) GetMatrixTransform() draw2d.Matrix {
	return draw2d.Matrix{}
}
// SetMatrixTransform sets the current transformation matrix
func (gc *GraphicContext) SetMatrixTransform(tr draw2d.Matrix) {

}
// ComposeMatrixTransform composes the current transformation matrix with tr
func (gc *GraphicContext) ComposeMatrixTransform(tr draw2d.Matrix) {

}
// Rotate applies a rotation to the current transformation matrix. angle is in radian.
func (gc *GraphicContext) Rotate(angle float64) {

}
// Translate applies a translation to the current transformation matrix.
func (gc *GraphicContext) Translate(tx, ty float64) {

}
// Scale applies a scale to the current transformation matrix.
func (gc *GraphicContext) Scale(sx, sy float64) {

}
// SetStrokeColor sets the current stroke color
func (gc *GraphicContext) SetStrokeColor(c color.Color) {

}
// SetFillColor sets the current fill color
func (gc *GraphicContext) SetFillColor(c color.Color) {

}
// SetFillRule sets the current fill rule
func (gc *GraphicContext) SetFillRule(f draw2d.FillRule) {

}
// SetLineWidth sets the current line width
func (gc *GraphicContext) SetLineWidth(lineWidth float64) {

}
// SetLineCap sets the current line cap
func (gc *GraphicContext) SetLineCap(cap draw2d.LineCap) {

}
// SetLineJoin sets the current line join
func (gc *GraphicContext) SetLineJoin(join draw2d.LineJoin) {

}
// SetLineDash sets the current dash
func (gc *GraphicContext) SetLineDash(dash []float64, dashOffset float64) {

}
// SetFontSize sets the current font size
func (gc *GraphicContext) SetFontSize(fontSize float64) {

}
// GetFontSize gets the current font size
func (gc *GraphicContext) GetFontSize() float64 {
	return 0
}
// SetFontData sets the current FontData
func (gc *GraphicContext) SetFontData(fontData draw2d.FontData) {

}
// GetFontData gets the current FontData
func (gc *GraphicContext) GetFontData() draw2d.FontData {
	return draw2d.FontData{}
}
// GetFontName gets the current FontData as a string
func (gc *GraphicContext) GetFontName() string {
	return ""
}
// DrawImage draws the raster image in the current canvas
func (gc *GraphicContext) DrawImage(image image.Image) {

}
// Save the context and push it to the context stack
func (gc *GraphicContext) Save() {

}
// Restore remove the current context and restore the last one
func (gc *GraphicContext) Restore() {

}
// Clear fills the current canvas with a default transparent color
func (gc *GraphicContext) Clear() {

}
// ClearRect fills the specified rectangle with a default transparent color
func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {

}
// SetDPI sets the current DPI
func (gc *GraphicContext) SetDPI(dpi int) {

}
// GetDPI gets the current DPI
func (gc *GraphicContext) GetDPI() int {
	return 0
}
// GetStringBounds gets pixel bounds(dimensions) of given string
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	return 0, 0, 0, 0
}
// CreateStringPath creates a path from the string s at x, y
func (gc *GraphicContext) CreateStringPath(text string, x, y float64) (cursor float64) {
	return 0
}
// FillString draws the text at point (0, 0)
func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return 0
}
// FillStringAt draws the text at the specified point (x, y)
func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	return 0
}
// StrokeString draws the contour of the text at point (0, 0)
func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return 0
}
// StrokeStringAt draws the contour of the text at point (x, y)
func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	return 0
}
// Stroke strokes the paths with the color specified by SetStrokeColor
func (gc *GraphicContext) Stroke(paths ...*draw2d.Path) {

}
// Fill fills the paths with the color specified by SetFillColor
func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {

}
// FillStroke first fills the paths and than strokes them
func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {

}