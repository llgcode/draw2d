package draw2dbase

import "github.com/llgcode/draw2d"

var glyphCache map[string]map[rune]*glyph

func init() {
	glyphCache = make(map[string]map[rune]*glyph)
}

// FillGlyph copies a glyph from the cache, copies it to the gc, and fills it
func FillGlyph(gc draw2d.GraphicContext, x, y float64, fontName string, chr rune) float64 {
	g := fetchGlyph(gc, fontName, chr)
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.Path)
	gc.Restore()
	return g.Width
}

// StrokeGlyph fetches a glyph from the cache, copies it to the gc, and strokes it
func StrokeGlyph(gc draw2d.GraphicContext, x, y float64, fontName string, chr rune) float64 {
	g := fetchGlyph(gc, fontName, chr)
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.Path)
	gc.Restore()
	return g.Width
}

// fetchGlyph fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func fetchGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *glyph {
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
