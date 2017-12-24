// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

package draw2dsvg

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"bytes"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
)

const (
)

var (
)

type drawType int
const (
	filled drawType = 1 << iota
	stroked
)

type SVG bytes.Buffer

func NewSvg() *Svg {
	return &Svg{Xmlns: "http://www.w3.org/2000/svg"}
}

// GraphicContext implements the draw2d.GraphicContext interface
// It provides draw2d with a svg backend
type GraphicContext struct {
	*draw2dbase.StackGraphicContext
	svg *Svg
}

func NewGraphicContext(svg *Svg) *GraphicContext {
	gc := &GraphicContext{draw2dbase.NewStackGraphicContext(), svg}
	return gc
}

// Clear fills the current canvas with a default transparent color
func (gc *GraphicContext) Clear() {
	gc.svg.Groups = nil
	gc.svg.Groups = append(gc.svg.Groups, Group{
		// TODO add background color?
	})
}

// Stroke strokes the paths with the color specified by SetStrokeColor
func (gc *GraphicContext) Stroke(paths ...*draw2d.Path) {
	gc.drawPaths(stroked, paths...)
	gc.Current.Path.Clear()
}
// Fill fills the paths with the color specified by SetFillColor
func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {
	gc.drawPaths(filled, paths...)
	gc.Current.Path.Clear()
}
// FillStroke first fills the paths and than strokes them
func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {
	gc.drawPaths(filled | stroked, paths...)
	gc.Current.Path.Clear()
}

func (gc *GraphicContext) drawPaths (drawType drawType, paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)

	svgPaths := make([]Path, len(paths))

	for i, path := range paths {
		svgPaths[i].Desc = toSvgPathDesc(path)
		if drawType & stroked == stroked {
			svgPaths[i].Stroke = toSvgRGBA(gc.Current.StrokeColor)
			svgPaths[i].StrokeWidth = toSvgLength(gc.Current.LineWidth)
			svgPaths[i].StrokeLinecap = gc.Current.Cap.String()
			svgPaths[i].StrokeLinejoin = gc.Current.Join.String()
		} else {
			svgPaths[i].Stroke = "none"
		}
		if drawType & filled == filled {
			svgPaths[i].Fill = toSvgRGBA(gc.Current.FillColor)
		} else {
			svgPaths[i].Fill = "none"
		}
	}

	gc.svg.Groups = append(gc.svg.Groups, Group{
		Paths: svgPaths,
	})
}

func toSvgRGBA (c color.Color) string { // TODO move elsewhere
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%v, %v, %v, %.3f)", r>>8, g>>8, b>>8, float64(a>>8)/255)
}

func toSvgLength (l float64) string {
	return fmt.Sprintf("%.4f", l)
}

func toSvgPathDesc (p *draw2d.Path) string { // TODO move elsewhere
	parts := make([]string, len(p.Components))
	i := 0
	for j, cmp := range p.Components {
		switch cmp {
		case draw2d.MoveToCmp:
			parts[j] = fmt.Sprintf("M %v %v", p.Points[i], p.Points[i+1])
			i += 2
		case draw2d.LineToCmp:
			parts[j] = fmt.Sprintf("L %v %v", p.Points[i], p.Points[i+1])
			i += 2
		case draw2d.QuadCurveToCmp:
			parts[j] = fmt.Sprintf("Q %v %v %v %v", p.Points[i], p.Points[i+1], p.Points[i+2], p.Points[i+3])
			i += 4
		case draw2d.CubicCurveToCmp:
			parts[j] = fmt.Sprintf("C %v %v %v %v %v %v", p.Points[i], p.Points[i+1], p.Points[i+2], p.Points[i+3], p.Points[i+4], p.Points[i+5])
			i += 6
		case draw2d.ArcToCmp:
			large := 0
			sweep := 0
			if p.Points[i+4] - p.Points[i+5] > 0 { // TODO this is probably not correct
				large = 1
				sweep = 1
			}
			// rx ry x-axis-rotation large-arc-flag sweep-flag x y
			parts[j] = fmt.Sprintf("A %v %v %v %v %v %v %v",
				p.Points[i+2], // rx
				p.Points[i+3], // ry
				0, // x-axis-rotation
				large, // large-arc-flag
				sweep, // // sweep-flag
				p.Points[i], // x // TODO this is center of arc not an endpont
				p.Points[i+1], // y // TODO -//-
			)
			i += 6
		case draw2d.CloseCmp:
			parts[j] = "Z"
		}
	}
	println("parts", parts)
	return strings.Join(parts, " ")
}
///////////////////////////////////////
// TODO implement following methods (or remove if not neccesary)

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
