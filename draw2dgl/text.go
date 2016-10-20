package draw2dgl

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"

	"golang.org/x/image/math/fixed"
)

var glyphCache map[string]map[rune]*glyph

func init() {
	glyphCache = make(map[string]map[rune]*glyph)
}

// FillGlyph copies a glyph from the cache, copies it to the gc, and fills it
func FillGlyph(gc draw2d.GraphicContext, x, y float64, chr rune) float64 {
	g := fetchGlyph(gc, chr)
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.Path)
	gc.Restore()
	return g.Width
}

// StrokeGlyph fetches a glyph from the cache, copies it to the gc, and strokes it
func StrokeGlyph(gc draw2d.GraphicContext, x, y float64, chr rune) float64 {
	g := fetchGlyph(gc, chr)
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.Path)
	gc.Restore()
	return g.Width
}

// fetchGlyph fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func fetchGlyph(gc draw2d.GraphicContext, chr rune) *glyph {
	fontData := gc.GetFontData()
	fontName := fmt.Sprintf("%s:%d:%d:%d", fontData.Name, fontData.Family, fontData.Style, gc.GetFontSize())
	if glyphCache[fontName] == nil {
		glyphCache[fontName] = make(map[rune]*glyph, 60)
	}
	if glyphCache[fontName][chr] == nil {
		glyphCache[fontName][chr] = renderGlyph(gc, fontName, chr)
	}
	return glyphCache[fontName][chr].Copy()
}

// renderGlyph renders a Glyph then caches and returns it
func renderGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *glyph {
	gc.Save()
	defer gc.Restore()
	gc.BeginPath()
	width := gc.CreateStringPath(string(chr), 0, 0)
	path := gc.GetPath()
	return &glyph{
		Path:  &path,
		Width: width,
	}
}

// glyph represents a rune which has been converted to a Path and width
type glyph struct {
	// Path represents a glyph, it is always at (0, 0)
	Path *draw2d.Path
	// Width of the glyph
	Width float64
}

func (g *glyph) Copy() *glyph {
	return &glyph{
		Path:  g.Path.Copy(),
		Width: g.Width,
	}
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
