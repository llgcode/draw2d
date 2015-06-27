// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels
// TODO: fonts, dpi

package pdf2d

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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

const c255 = 255.0 / 65535.0

var (
	imageCount uint32
	white      color.Color = color.RGBA{255, 255, 255, 255}
)

// NewPdf creates a new pdf document with the draw2d fontfolder, adds
// a page and set fill color to white.
func NewPdf(orientationStr, unitStr, sizeStr string) *gofpdf.Fpdf {
	pdf := gofpdf.New(orientationStr, unitStr, sizeStr, draw2d.GetFontFolder())
	pdf.AddPage()
	pdf.SetFillColor(255, 255, 255) // to be compatible with draw2d
	return pdf
}

// rgb converts a color (used by draw2d) into 3 int (used by gofpdf)
func rgb(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(float64(r) * c255), int(float64(g) * c255), int(float64(b) * c255)
}

// clearRect draws a white rectangle
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

// DrawImage draws an image as PNG
// TODO: add type (tp) as parameter to argument list?
func (gc *GraphicContext) DrawImage(image image.Image) {
	// TODO: fix rotation gc.Current.Tr
	name := strconv.Itoa(int(imageCount))
	tp := "PNG" // "JPG", "JPEG", "PNG" and "GIF"
	b := &bytes.Buffer{}
	png.Encode(b, image)
	gc.pdf.RegisterImageReader(name, tp, b)
	bounds := image.Bounds()
	//x0, y0, x1, y1 := float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Dx()), float64(bounds.Dy())
	x0, y0, x1, y1 := float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Max.X), float64(bounds.Max.Y)
	tr := gc.Current.Tr
	tr.TransformRectangle(&x0, &y0, &x1, &y1)
	gc.pdf.Image(name, x0, y0, x1-x0, y1-y0, false, tp, 0, "")
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
	// TODO: fix rotation of gc.Current.Tr
	_, _, w, h := gc.GetStringBounds(text)
	x1, y1 := x+w, y+h
	tr := gc.Current.Tr
	tr.TransformRectangle(&x, &y, &x1, &y1)
	gc.pdf.MoveTo(x, y)
	gc.pdf.Cell(x1-x, y1-y, text)
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

// StrokeString draws a string at 0, 0 (stroking is unsupported,
// string will be filled)
func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

// StrokeStringAt draws a string at x, y (stroking is unsupported,
// string will be filled)
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
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.Current.Tr, gc.pdf))
	// pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.Current.Tr,NewPathLogger(logger, gc.pdf)))
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

// SetFont is unsupported by the pdf graphic context, use SetFontData
// instead.
func (gc *GraphicContext) SetFont(font *truetype.Font) {
	// TODO: what to do with this api conflict between draw2d and gofpdf?!
}

// SetFontData sets the current font used to draw text. Always use
// this method, as SetFont is unsupported by the pdf graphic context.
// It is mandatory to call this method at least once before printing
// text or the resulting document will not be valid.
// It is necessary to generate a font definition file first with the
// makefont utility. It is not necessary to call this function for the
// core PDF fonts (courier, helvetica, times, zapfdingbats).
// go get github.com/jung-kurt/gofpdf/makefont
// http://godoc.org/github.com/jung-kurt/gofpdf#Fpdf.AddFont
func (gc *GraphicContext) SetFontData(fontData draw2d.FontData) {
	gc.StackGraphicContext.SetFontData(fontData)
	var style string
	if fontData.Style&draw2d.FontStyleBold != 0 {
		style += "B"
	}
	if fontData.Style&draw2d.FontStyleItalic != 0 {
		style += "I"
	}
	fn := draw2d.FontFileName(fontData)
	fn = fn[:len(fn)-4]
	gc.pdf.AddFont(fn, style, fn+".json")
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
