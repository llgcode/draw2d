// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels
// TODO: fonts, dpi

package pdf2d

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"strconv"

	"code.google.com/p/freetype-go/freetype/truetype"

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

var (
	imageCount uint32
	white      color.Color = color.RGBA{255, 255, 255, 255}
)

func rgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(float64(r) * c255), int(float64(g) * c255), int(float64(b) * c255)
}

func clearRect(gc *GraphicContext, x1, y1, x2, y2 float64) {
	// save state
	f := gc.Current.FillColor
	x, y := gc.pdf.GetXY()
	// cover page with white rectangle
	gc.SetFillColor(white)
	draw2d.Rect(gc, x1, y1, x2, y2)
	gc.Fill()
	// restore state
	gc.SetFillColor(f)
	gc.pdf.MoveTo(x, y)
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

// DrawImage draws an image as JPG at 96dpi
func (gc *GraphicContext) DrawImage(image image.Image) {
	name := strconv.Itoa(int(imageCount))
	tp := "JPG" // "JPG", "JPEG", "PNG" and "GIF"
	b := &bytes.Buffer{}
	jpeg.Encode(b, image, nil)
	gc.pdf.RegisterImageReader(name, tp, b)
	gc.pdf.Image(name, 0, 0, 0, 0, false, tp, 0, "")
	// bounds := image.Bounds()
	// x, y, w, h := float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Dx()), float64(bounds.Dy())
	//gc.pdf.Image(name, x, y, w, h, false, tp, 0, "")
}

// Clear draws a white rectangle over the whole page
func (gc *GraphicContext) Clear() {
	width, height := gc.pdf.GetPageSize()
	clearRect(gc, 0, 0, width, height)
}

// ClearRect draws a white rectangle over the specified area
func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	clearRect(gc, float64(x1), float64(y1), float64(x2), float64(y2))
}

// SetDPI is a dummy method to implement the GraphicContext interface
func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	// gc.recalc()
}

// GetDPI is a dummy method to implement the GraphicContext interface
func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	_, h := gc.pdf.GetFontSize()
	return 0, 0, gc.pdf.GetStringWidth(s), h
}

// CreateStringPath creates a path from the string s at x, y, and returns the string width.
func (gc *GraphicContext) CreateStringPath(text string, x, y float64) (cursor float64) {
	gc.pdf.MoveTo(x, y)
	_, _, w, h := gc.GetStringBounds(text)
	gc.pdf.Cell(w, h, text)
	return w
}

// FillString draws a string at 0, 0
func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return gc.FillStringAt(text, 0, 0)
}

// FillStringAt draws a string at x, y
func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	return gc.CreateStringPath(text, x, y)
}

// StrokeString draws a string at 0, 0
func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

// StrokeStringAt draws a string at x, y
func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	return gc.CreateStringPath(text, x, y)
}

// Stroke strokes the paths
func (gc *GraphicContext) Stroke(paths ...*draw2d.PathStorage) {
	gc.draw("D", paths...)
}

// Fill strokes the paths
func (gc *GraphicContext) Fill(paths ...*draw2d.PathStorage) {
	gc.draw("F", paths...)
}

// FillStroke first fills the paths and than strokes them
func (gc *GraphicContext) FillStroke(paths ...*draw2d.PathStorage) {
	gc.draw("FD", paths...)
}

var logger = log.New(os.Stdout, "", log.Lshortfile)

// draw fills and/or strokes paths
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

// SetStrokeColor sets the stroke color
func (gc *GraphicContext) SetStrokeColor(c color.Color) {
	gc.StackGraphicContext.SetStrokeColor(c)
	gc.pdf.SetDrawColor(rgb(c))
}

// SetFillColor sets the fill color
func (gc *GraphicContext) SetFillColor(c color.Color) {
	gc.StackGraphicContext.SetFillColor(c)
	gc.pdf.SetFillColor(rgb(c))
}

// SetFont sets the font used to draw text.
// It is mandatory to call this method at least once before printing
// text or the resulting document will not be valid.
func (gc *GraphicContext) SetFont(font *truetype.Font) {
	// TODO: this api conflict needs to be fixed
	gc.pdf.SetFont("Helvetica", "", 12)
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
func (gc *GraphicContext) SetFontSize(fontSize float64) {
	gc.StackGraphicContext.SetFontSize(fontSize)
	gc.pdf.SetFontSize(fontSize)
	//gc.recalc()
}

// SetLineWidth sets the line width
func (gc *GraphicContext) SetLineWidth(LineWidth float64) {
	gc.StackGraphicContext.SetLineWidth(LineWidth)
	gc.pdf.SetLineWidth(LineWidth)
}

// SetLineCap sets the line cap (round, but or square)
func (gc *GraphicContext) SetLineCap(Cap draw2d.Cap) {
	gc.StackGraphicContext.SetLineCap(Cap)
	gc.pdf.SetLineCapStyle(caps[Cap])
}
