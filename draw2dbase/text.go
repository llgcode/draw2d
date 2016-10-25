package draw2dbase

import (
	"errors"
	"log"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"


	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	glyphBuf        = &truetype.GlyphBuf{}
	DefaultFontData = draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}
)


func loadCurrentFont(gc draw2d.GraphicContext) (*truetype.Font, error) {
	font := draw2d.GetFont(gc.GetFontData())
	if font == nil {
		font = draw2d.GetFont(DefaultFontData)
	}
	if font == nil {
		return nil, errors.New("No font set, and no default font available.")
	}
	gc.SetFont(font)
	gc.SetFontSize(gc.GetFontSize())
	recalc(gc)
	return font, nil
}

func drawGlyph(gc draw2d.GraphicContext, glyph truetype.Index, dx, dy float64) error {
	if err := glyphBuf.Load(gc.GetFont(), fixed.Int26_6(gc.GetScale()), glyph, font.HintingNone); err != nil {
		return err
	}
	e0 := 0
	for _, e1 := range glyphBuf.Ends {
		DrawContour(gc, glyphBuf.Points[e0:e1], dx, dy)
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
func CreateStringPath(gc draw2d.GraphicContext, s string, x, y float64) float64 {
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return 0.0
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	scale := gc.GetScale()
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(scale), prev, index))
		}
		err := drawGlyph(gc, index, x, y)
		if err != nil {
			log.Println(err)
			return startx - x
		}
		x += fUnitsToFloat64(f.HMetric(fixed.Int26_6(scale), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x - startx
}

// GetStringBounds returns the approximate pixel bounds of the string s at x, y.
// The the left edge of the em square of the first character of s
// and the baseline intersect at 0, 0 in the returned coordinates.
// Therefore the top and left coordinates may well be negative.
func GetStringBounds(gc draw2d.GraphicContext, s string) (left, top, right, bottom float64) {
	f, err := loadCurrentFont(gc)
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0
	}
	top, left, bottom, right = 10e6, 10e6, -10e6, -10e6
	cursor := 0.0
	prev, hasPrev := truetype.Index(0), false
	scale := gc.GetScale()
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			cursor += fUnitsToFloat64(f.Kern(fixed.Int26_6(scale), prev, index))
		}
		if err := glyphBuf.Load(gc.GetFont(), fixed.Int26_6(scale), index, font.HintingNone); err != nil {
			log.Println(err)
			return 0, 0, 0, 0
		}
		e0 := 0
		for _, e1 := range glyphBuf.Ends {
			ps := glyphBuf.Points[e0:e1]
			for _, p := range ps {
				x, y := pointToF64Point(p)
				top = math.Min(top, y)
				bottom = math.Max(bottom, y)
				left = math.Min(left, x+cursor)
				right = math.Max(right, x+cursor)
			}
		}
		cursor += fUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.GetScale()), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return left, top, right, bottom
}

// recalc recalculates scale and bounds values from the font size, screen
// resolution and font metrics, and invalidates the glyph cache.
func recalc(gc draw2d.GraphicContext) {
	gc.SetScale(gc.GetFontSize() * float64(gc.GetDPI()) * (64.0 / 72.0))
}

// DrawContour draws the given closed contour at the given sub-pixel offset.
func DrawContour(path draw2d.PathBuilder, ps []truetype.Point, dx, dy float64) {
	if len(ps) == 0 {
		return
	}
	startX, startY := pointToF64Point(ps[0])
	path.MoveTo(startX+dx, startY+dy)
	q0X, q0Y, on0 := startX, startY, true
	for _, p := range ps[1:] {
		qX, qY := pointToF64Point(p)
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				path.LineTo(qX+dx, qY+dy)
			} else {
				path.QuadCurveTo(q0X+dx, q0Y+dy, qX+dx, qY+dy)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				midX := (q0X + qX) / 2
				midY := (q0Y + qY) / 2
				path.QuadCurveTo(q0X+dx, q0Y+dy, midX+dx, midY+dy)
			}
		}
		q0X, q0Y, on0 = qX, qY, on
	}
	// Close the curve.
	if on0 {
		path.LineTo(startX+dx, startY+dy)
	} else {
		path.QuadCurveTo(q0X+dx, q0Y+dy, startX+dx, startY+dy)
	}
}

func pointToF64Point(p truetype.Point) (x, y float64) {
	return fUnitsToFloat64(p.X), -fUnitsToFloat64(p.Y)
}

func fUnitsToFloat64(x fixed.Int26_6) float64 {
	scaled := x << 2
	return float64(scaled/256) + float64(scaled%256)/256.0
}

// FontExtents contains font metric information.
type FontExtents struct {
	// Ascent is the distance that the text
	// extends above the baseline.
	Ascent float64

	// Descent is the distance that the text
	// extends below the baseline.  The descent
	// is given as a negative value.
	Descent float64

	// Height is the distance from the lowest
	// descending point to the highest ascending
	// point.
	Height float64
}

// Extents returns the FontExtents for a font.
// TODO needs to read this https://developer.apple.com/fonts/TrueType-Reference-Manual/RM02/Chap2.html#intro
func Extents(font *truetype.Font, size float64) FontExtents {
	bounds := font.Bounds(fixed.Int26_6(font.FUnitsPerEm()))
	scale := size / float64(font.FUnitsPerEm())
	return FontExtents{
		Ascent:  float64(bounds.Max.Y) * scale,
		Descent: float64(bounds.Min.Y) * scale,
		Height:  float64(bounds.Max.Y-bounds.Min.Y) * scale,
	}
}
