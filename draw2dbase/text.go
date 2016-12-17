package draw2dbase

import "github.com/llgcode/draw2d"

var (
	glyphCacheDefault            = &defaultGlyphCache{make(map[string]map[rune]*Glyph)}
	glyphCache        GlyphCache = glyphCacheDefault
)

// Types implementing this interface can be passed to SetGlyphCache to change the
// way glyphs are being stored and retrieved.
type GlyphCache interface {
	// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
	Fetch(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph

	// Render renders a glyph then returns it
	Render(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph
}

// Changes the glyph cache backend used by the package.
// To restore the default glyph cache, call this function passing nil as argument.
func SetGlyphCache(cache GlyphCache) {
	if cache == nil {
		glyphCache = glyphCacheDefault
	} else {
		glyphCache = cache
	}
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
		path:  &path,
		Width: width,
	}
}

// FetchGlyph fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func FetchGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	return glyphCache.Fetch(gc, fontName, chr)
}

// renderGlyph renders a glyph then caches and returns it
func renderGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *Glyph {
	return glyphCache.Render(gc, fontName, chr)
}

// Glyph represents a rune which has been converted to a Path and width
type Glyph struct {
	// path represents a glyph, it is always at (0, 0)
	path *draw2d.Path
	// Width of the glyph
	Width float64
}

// Copy copys the Glyph, and returns the copy
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
