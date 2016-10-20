// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2dbase

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/llgcode/draw2d"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var DefaultFontData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}

type StackGraphicContext struct {
	Current *ContextStack
}

type ContextStack struct {
	Tr          draw2d.Matrix
	Path        *draw2d.Path
	LineWidth   float64
	Dash        []float64
	DashOffset  float64
	StrokeColor color.Color
	FillColor   color.Color
	FillRule    draw2d.FillRule
	Cap         draw2d.LineCap
	Join        draw2d.LineJoin
	FontSize    float64
	FontData    draw2d.FontData

	Font *truetype.Font
	// fontSize and dpi are used to calculate scale. scale is the number of
	// 26.6 fixed point units in 1 em.
	Scale float64
	glyphBuf         *truetype.GlyphBuf
	DPI		int

	Previous *ContextStack
}

/**
 * Create a new Graphic context from an image
 */
func NewStackGraphicContext() *StackGraphicContext {
	gc := &StackGraphicContext{}
	gc.Current = new(ContextStack)
	gc.Current.Tr = draw2d.NewIdentityMatrix()
	gc.Current.Path = new(draw2d.Path)
	gc.Current.LineWidth = 1.0
	gc.Current.StrokeColor = image.Black
	gc.Current.FillColor = image.White
	gc.Current.Cap = draw2d.RoundCap
	gc.Current.FillRule = draw2d.FillRuleEvenOdd
	gc.Current.Join = draw2d.RoundJoin
	gc.Current.FontSize = 10
	gc.Current.FontData = DefaultFontData
	gc.Current.glyphBuf = &truetype.GlyphBuf{}
	gc.Current.DPI = 92
	return gc
}

func (gc *StackGraphicContext) GetMatrixTransform() draw2d.Matrix {
	return gc.Current.Tr
}

func (gc *StackGraphicContext) SetMatrixTransform(Tr draw2d.Matrix) {
	gc.Current.Tr = Tr
}

func (gc *StackGraphicContext) ComposeMatrixTransform(Tr draw2d.Matrix) {
	gc.Current.Tr.Compose(Tr)
}

func (gc *StackGraphicContext) Rotate(angle float64) {
	gc.Current.Tr.Rotate(angle)
}

func (gc *StackGraphicContext) Translate(tx, ty float64) {
	gc.Current.Tr.Translate(tx, ty)
}

func (gc *StackGraphicContext) Scale(sx, sy float64) {
	gc.Current.Tr.Scale(sx, sy)
}

func (gc *StackGraphicContext) SetStrokeColor(c color.Color) {
	gc.Current.StrokeColor = c
}

func (gc *StackGraphicContext) SetFillColor(c color.Color) {
	gc.Current.FillColor = c
}

func (gc *StackGraphicContext) SetFillRule(f draw2d.FillRule) {
	gc.Current.FillRule = f
}

func (gc *StackGraphicContext) SetLineWidth(lineWidth float64) {
	gc.Current.LineWidth = lineWidth
}

func (gc *StackGraphicContext) SetLineCap(cap draw2d.LineCap) {
	gc.Current.Cap = cap
}

func (gc *StackGraphicContext) SetLineJoin(join draw2d.LineJoin) {
	gc.Current.Join = join
}

func (gc *StackGraphicContext) SetLineDash(dash []float64, dashOffset float64) {
	gc.Current.Dash = dash
	gc.Current.DashOffset = dashOffset
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
func (gc *StackGraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
	gc.recalc()
}

func (gc *StackGraphicContext) GetFontSize() float64 {
	return gc.Current.FontSize
}

func (gc *StackGraphicContext) SetFontData(fontData draw2d.FontData) {
	gc.Current.FontData = fontData
}

func (gc *StackGraphicContext) GetFontData() draw2d.FontData {
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
	context.Font = gc.Current.Font
	context.Scale = gc.Current.Scale
	context.glyphBuf = gc.Current.glyphBuf
	context.DPI = gc.Current.DPI
	copy(context.Tr[:], gc.Current.Tr[:])
	context.Previous = gc.Current
	gc.Current = context
}

func (gc *StackGraphicContext) Restore() {
	if gc.Current.Previous != nil {
		oldContext := gc.Current
		gc.Current = gc.Current.Previous
		oldContext.Previous = nil
	}
}

func (gc *StackGraphicContext) loadCurrentFont() (*truetype.Font, error) {
	font := draw2d.GetFont(gc.Current.FontData)
	if font == nil {
		font = draw2d.GetFont(DefaultFontData)
	}
	if font == nil {
		return nil, errors.New("No font set, and no default font available.")
	}
	gc.SetFont(font)
	gc.SetFontSize(gc.Current.FontSize)
	return font, nil
}

func (gc *StackGraphicContext) drawGlyph(glyph truetype.Index, dx, dy float64) error {
	if err := gc.Current.glyphBuf.Load(gc.Current.Font, fixed.Int26_6(gc.Current.Scale), glyph, font.HintingNone); err != nil {
		return err
	}
	e0 := 0
	for _, e1 := range gc.Current.glyphBuf.Ends {
		DrawContour(gc, gc.Current.glyphBuf.Points[e0:e1], dx, dy)
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
func (gc *StackGraphicContext) CreateStringPath(s string, x, y float64) float64 {
	f, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0.0
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		err := gc.drawGlyph(index, x, y)
		if err != nil {
			log.Println(err)
			return startx - x
		}
		x += fUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.Current.Scale), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x - startx
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
// The the left edge of the em square of the first character of s
// and the baseline intersect at 0, 0 in the returned coordinates.
// Therefore the top and left coordinates may well be negative.
func (gc *StackGraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	f, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0
	}
	top, left, bottom, right = 10e6, 10e6, -10e6, -10e6
	cursor := 0.0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			cursor += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		if err := gc.Current.glyphBuf.Load(gc.Current.Font, fixed.Int26_6(gc.Current.Scale), index, font.HintingNone); err != nil {
			log.Println(err)
			return 0, 0, 0, 0
		}
		e0 := 0
		for _, e1 := range gc.Current.glyphBuf.Ends {
			ps := gc.Current.glyphBuf.Points[e0:e1]
			for _, p := range ps {
				x, y := pointToF64Point(p)
				top = math.Min(top, y)
				bottom = math.Max(bottom, y)
				left = math.Min(left, x+cursor)
				right = math.Max(right, x+cursor)
			}
		}
		cursor += fUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.Current.Scale), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return left, top, right, bottom
}

// recalc recalculates scale and bounds values from the font size, screen
// resolution and font metrics, and invalidates the glyph cache.
func (gc *StackGraphicContext) recalc() {
	gc.Current.Scale = gc.Current.FontSize * float64(gc.Current.DPI) * (64.0 / 72.0)
}

func (gc *StackGraphicContext) SetDPI(dpi int) {
	gc.Current.DPI = dpi
	gc.recalc()
}

// SetFont sets the font used to draw text.
func (gc *StackGraphicContext) SetFont(font *truetype.Font) {
	gc.Current.Font = font
}

func (gc *StackGraphicContext) GetDPI() int {
	return gc.Current.DPI
}
