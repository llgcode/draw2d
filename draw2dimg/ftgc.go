// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2dimg

import (
	"errors"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type Painter interface {
	raster.Painter
	SetColor(color color.Color)
}

type GraphicContext struct {
	*draw2dbase.StackGraphicContext
	img              draw.Image
	painter          Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	glyphBuf         *truetype.GlyphBuf
	DPI              int
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(img draw.Image) *GraphicContext {
	var painter Painter
	switch selectImage := img.(type) {
	case *image.RGBA:
		painter = raster.NewRGBAPainter(selectImage)
	default:
		panic("Image type not supported")
	}
	return NewGraphicContextWithPainter(img, painter)
}

// Create a new Graphic context from an image and a Painter (see Freetype-go)
func NewGraphicContextWithPainter(img draw.Image, painter Painter) *GraphicContext {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	dpi := 92
	gc := &GraphicContext{
		draw2dbase.NewStackGraphicContext(),
		img,
		painter,
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		truetype.NewGlyphBuf(),
		dpi,
	}
	return gc
}

func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

func (gc *GraphicContext) Clear() {
	width, height := gc.img.Bounds().Dx(), gc.img.Bounds().Dy()
	gc.ClearRect(0, 0, width, height)
}

func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	imageColor := image.NewUniform(gc.Current.FillColor)
	draw.Draw(gc.img, image.Rect(x1, y1, x2, y2), imageColor, image.ZP, draw.Over)
}

func (gc *GraphicContext) DrawImage(img image.Image) {
	DrawImage(img, gc.img, gc.Current.Tr, draw.Over, BilinearFilter)
}

func (gc *GraphicContext) FillString(text string) (cursor float64) {
	return gc.FillStringAt(text, 0, 0)
}

func (gc *GraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	width := gc.CreateStringPath(text, x, y)
	gc.Fill()
	return width
}

func (gc *GraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	width := gc.CreateStringPath(text, x, y)
	gc.Stroke()
	return width
}

func (gc *GraphicContext) loadCurrentFont() (*truetype.Font, error) {
	font := draw2d.GetFont(gc.Current.FontData)
	if font == nil {
		font = draw2d.GetFont(draw2dbase.DefaultFontData)
	}
	if font == nil {
		return nil, errors.New("No font set, and no default font available.")
	}
	gc.SetFont(font)
	gc.SetFontSize(gc.Current.FontSize)
	return font, nil
}

// p is a truetype.Point measured in FUnits and positive Y going upwards.
// The returned value is the same thing measured in floating point and positive Y
// going downwards.

func (gc *GraphicContext) drawGlyph(glyph truetype.Index, dx, dy float64) error {
	if err := gc.glyphBuf.Load(gc.Current.Font, gc.Current.Scale, glyph, truetype.NoHinting); err != nil {
		return err
	}
	e0 := 0
	for _, e1 := range gc.glyphBuf.End {
		DrawContour(gc, gc.glyphBuf.Point[e0:e1], dx, dy)
		e0 = e1
	}
	return nil
}

// CreateStringPath creates a path from the string s at x, y, and returns the string width.
// The text is placed so that the left edge of the em square of the first character of s
// and the baseline intersect at x, y. The majority of the affected pixels will be
// above and to the right of the point, but some may be below or to the left.
// For example, drawing a string that starts with a 'J' in an italic font may
// affect pixels below and left of the point.
func (gc *GraphicContext) CreateStringPath(s string, x, y float64) float64 {
	font, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0.0
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := font.Index(rune)
		if hasPrev {
			x += fUnitsToFloat64(font.Kerning(gc.Current.Scale, prev, index))
		}
		err := gc.drawGlyph(index, x, y)
		if err != nil {
			log.Println(err)
			return startx - x
		}
		x += fUnitsToFloat64(font.HMetric(gc.Current.Scale, index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x - startx
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
// The the left edge of the em square of the first character of s
// and the baseline intersect at 0, 0 in the returned coordinates.
// Therefore the top and left coordinates may well be negative.
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	font, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0
	}
	top, left, bottom, right = 10e6, 10e6, -10e6, -10e6
	cursor := 0.0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := font.Index(rune)
		if hasPrev {
			cursor += fUnitsToFloat64(font.Kerning(gc.Current.Scale, prev, index))
		}
		if err := gc.glyphBuf.Load(gc.Current.Font, gc.Current.Scale, index, truetype.NoHinting); err != nil {
			log.Println(err)
			return 0, 0, 0, 0
		}
		e0 := 0
		for _, e1 := range gc.glyphBuf.End {
			ps := gc.glyphBuf.Point[e0:e1]
			for _, p := range ps {
				x, y := pointToF64Point(p)
				top = math.Min(top, y)
				bottom = math.Max(bottom, y)
				left = math.Min(left, x+cursor)
				right = math.Max(right, x+cursor)
			}
		}
		cursor += fUnitsToFloat64(font.HMetric(gc.Current.Scale, index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return left, top, right, bottom
}

// recalc recalculates scale and bounds values from the font size, screen
// resolution and font metrics, and invalidates the glyph cache.
func (gc *GraphicContext) recalc() {
	gc.Current.Scale = int32(gc.Current.FontSize * float64(gc.DPI) * (64.0 / 72.0))
}

// SetDPI sets the screen resolution in dots per inch.
func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.recalc()
}

// SetFont sets the font used to draw text.
func (gc *GraphicContext) SetFont(font *truetype.Font) {
	gc.Current.Font = font
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
func (gc *GraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
	gc.recalc()
}

func (gc *GraphicContext) paint(rasterizer *raster.Rasterizer, color color.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.Current.Path.Clear()
}

func (gc *GraphicContext) Stroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2dbase.Transformer{gc.Current.Tr, draw2dbase.FtLineBuilder{gc.strokeRasterizer}})
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	var liner draw2dbase.Flattener
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
	} else {
		liner = stroker
	}
	for _, p := range paths {
		draw2dbase.Flatten(p, liner, gc.Current.Tr.GetScale())
	}

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = useNonZeroWinding(gc.Current.FillRule)

	/**** first method ****/
	flattener := draw2dbase.Transformer{gc.Current.Tr, draw2dbase.FtLineBuilder{gc.fillRasterizer}}
	for _, p := range paths {
		draw2dbase.Flatten(p, flattener, gc.Current.Tr.GetScale())
	}

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
}

func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = useNonZeroWinding(gc.Current.FillRule)
	gc.strokeRasterizer.UseNonZeroWinding = true

	flattener := draw2dbase.Transformer{gc.Current.Tr, draw2dbase.FtLineBuilder{gc.fillRasterizer}}

	stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, draw2dbase.Transformer{gc.Current.Tr, draw2dbase.FtLineBuilder{gc.strokeRasterizer}})
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	var liner draw2dbase.Flattener
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
	} else {
		liner = stroker
	}

	demux := draw2dbase.DemuxFlattener{[]draw2dbase.Flattener{flattener, liner}}
	for _, p := range paths {
		draw2dbase.Flatten(p, demux, gc.Current.Tr.GetScale())
	}

	// Fill
	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	// Stroke
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func useNonZeroWinding(f draw2d.FillRule) bool {
	switch f {
	case draw2d.FillRuleEvenOdd:
		return false
	case draw2d.FillRuleWinding:
		return true
	}
	return false
}
