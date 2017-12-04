package draw2dbase

import "github.com/llgcode/draw2d"

// GlyphCache manage a map of glyphs
type GlyphCache struct {
	glyphs map[string]map[rune]*Glyph
}


// NewGlyphCache initializes a GlyphCache
func NewGlyphCache() *GlyphCache {
	glyphs := make(map[string]map[rune]*Glyph)
	return &GlyphCache {
		glyphs: glyphs,
	}
}

// FetchGlyph fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func (glyphCache *GlyphCache) FetchGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	if glyphCache.glyphs[fontName] == nil {
		glyphCache.glyphs[fontName] = make(map[rune]*Glyph, 60)
	}
	if glyphCache.glyphs[fontName][chr] == nil {
		glyphCache.glyphs[fontName][chr] = renderGlyph(gc, fontName, chr)
	}
	return glyphCache.glyphs[fontName][chr].Copy()
}

// renderGlyph renders a glyph then caches and returns it
func renderGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	gc.Save()
	defer gc.Restore()
	gc.BeginPath()
	width := gc.CreateStringPath(string(chr), 0, 0)
	path := gc.GetPath()
	return &Glyph{
		path:  &path,
		Width: width,
	}
}

// Glyph represents a rune which has been converted to a Path and width
type Glyph struct {
	// path represents a glyph, it is always at (0, 0)
	path *draw2d.Path
	// Width of the glyph
	Width float64
}

// Returns a copy of a Glyph
func (g *Glyph) Copy() *Glyph {
	return &Glyph{
		path:  g.path.Copy(),
		Width: g.Width,
	}
}

// Fill copies a glyph from the cache, and fills it
func (g *Glyph) Fill(gc draw2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.path)
	gc.Restore()
	return g.Width
}

// Stroke fetches a glyph from the cache, and strokes it
func (g *Glyph) Stroke(gc draw2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.path)
	gc.Restore()
	return g.Width
}
