// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels

package pdf2d

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/stanim/draw2d"
	"github.com/stanim/gofpdf"
)

var (
	caps = map[draw2d.Cap]string{
		draw2d.RoundCap:  "round",
		draw2d.ButtCap:   "butt",
		draw2d.SquareCap: "square"}
)

func notImplemented(method string) {
	fmt.Printf("%s: not implemented\n", method)
}

const c255 = 255.0 / 65535.0

func rgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(float64(r) * c255), int(float64(g) * c255), int(float64(b) * c255)
}

// GraphicContext implements the draw2d.GraphicContext interface
// It provides draw2d with a pdf backend (based on gofpdf)
type GraphicContext struct {
	*draw2d.StackGraphicContext
	pdf *gofpdf.Fpdf
	DPI int
}

// NewGraphicContext creates a new pdf GraphicContext
func NewGraphicContext(pdf *gofpdf.Fpdf) *GraphicContext {
	dpi := 92
	return &GraphicContext{draw2d.NewStackGraphicContext(), pdf, dpi}
}

func (gc *GraphicContext) DrawImage(image image.Image) {
	notImplemented("DrawImage")
}
func (gc *GraphicContext) Clear() {
	notImplemented("Clear")
}

func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	notImplemented("ClearRect")
}

func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
}

func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	notImplemented("GetStringBounds")
	return 0, 0, 0, 0
}

func (gc *GraphicContext) CreateStringPath(text string, x, y float64) (cursor float64) {
	notImplemented("CreateStringPath")
	return 0
}

func (gc *GraphicContext) FillString(text string) (cursor float64) {
	notImplemented("FillString")
	return 0
}

func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	notImplemented("FillStringAt")
	return 0
}

func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	notImplemented("StrokeString")
	return 0
}

func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	notImplemented("StrokeStringAt")
	return 0
}

func (gc *GraphicContext) Stroke(paths ...*draw2d.PathStorage) {
	gc.draw("D", paths...)
}

func (gc *GraphicContext) Fill(paths ...*draw2d.PathStorage) {
	gc.draw("F", paths...)
}

func (gc *GraphicContext) FillStroke(paths ...*draw2d.PathStorage) {
	gc.draw("FD", paths...)
}

var logger *log.Logger = log.New(os.Stdout, "", log.Lshortfile)

func (gc *GraphicContext) draw(style string, paths ...*draw2d.PathStorage) {
	paths = append(paths, gc.Current.Path)
	pathConverter := NewPathConverter(
		NewVertexMatrixTransform(gc.Current.Tr,
			NewPathLogger(logger, gc.pdf)))
	pathConverter.Convert(paths...)
	if gc.Current.FillRule.UseNonZeroWinding() {
		style += "*"
	}
	gc.pdf.DrawPath(style)
}

// overwrite StackGraphicContext methods

func (gc *GraphicContext) SetStrokeColor(c color.Color) {
	gc.StackGraphicContext.SetStrokeColor(c)
	gc.pdf.SetDrawColor(rgb(c))
}

func (gc *GraphicContext) SetFillColor(c color.Color) {
	gc.StackGraphicContext.SetFillColor(c)
	gc.pdf.SetFillColor(rgb(c))
}

func (gc *GraphicContext) SetLineWidth(LineWidth float64) {
	gc.StackGraphicContext.SetLineWidth(LineWidth)
	gc.pdf.SetLineWidth(LineWidth)
}

func (gc *GraphicContext) SetLineCap(Cap draw2d.Cap) {
	gc.StackGraphicContext.SetLineCap(Cap)
	gc.pdf.SetLineCapStyle(caps[Cap])
}
