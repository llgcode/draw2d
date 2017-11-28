package draw2dbase

import "github.com/llgcode/draw2d"

var (
	DefaultGlyphCache = &defaultGlyphCache{make(map[string]map[rune]*Glyph)}
)

// Types implementing this interface can be passed to gc.SetGlyphCache to change the
// way glyphs are being stored and retrieved.
type GlyphCache interface {
	// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
	Fetch(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph

	// Render renders a glyph then returns it
	Render(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph
}

// NewGlyphCache creates and returns a new GlyphCache
func NewGlyphCache() GlyphCache {
	return &defaultGlyphCache{make(map[string]map[rune]*Glyph)}
}

type defaultGlyphCache struct {
	glyphMap map[string]map[rune]*Glyph
}

// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
func (cache *defaultGlyphCache) Fetch(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	if cache.glyphMap[fontName] == nil {
		cache.glyphMap[fontName] = make(map[rune]*Glyph, 60)
	}
	if cache.glyphMap[fontName][chr] == nil {
		cache.glyphMap[fontName][chr] = cache.Render(gc, fontName, chr)
	}
	return cache.glyphMap[fontName][chr].Copy()
}

// Render renders a glyph then returns it
func (cache *defaultGlyphCache) Render(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	gc.Save()
	defer gc.Restore()
	gc.BeginPath()
	width := gc.CreateStringPath(string(chr), 0, 0)
	path := gc.GetPath()
	return &Glyph{
		Path:  &path,
		Width: width,
	}
}

// Glyph represents a rune which has been converted to a Path and width
type Glyph struct {
	// Path represents a glyph, it is always at (0, 0)
	Path *draw2d.Path
	// Width of the glyph
	Width float64
}

// Copy copys the Glyph, and returns the copy
func (g *Glyph) Copy() *Glyph {
	return &Glyph{
		Path:  g.Path.Copy(),
		Width: g.Width,
	}
}

// Fill fills a Glyph in the specified coordinates
func (g *Glyph) Fill(gc draw2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.Path)
	gc.Restore()
	return g.Width
}

// Stroke strokes a Glyph in the specified coordinates
func (g *Glyph) Stroke(gc draw2d.GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.Path)
	gc.Restore()
	return g.Width
}
