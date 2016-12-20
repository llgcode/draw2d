package draw2dbase

import "github.com/llgcode/draw2d"

var (
	DefaultGlyphCache = &defaultGlyphCache{make(map[string]map[rune]*draw2d.Glyph)}
)

type defaultGlyphCache struct {
	glyphMap map[string]map[rune]*draw2d.Glyph
}

// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
func (cache *defaultGlyphCache) Fetch(gc draw2d.GraphicContext, fontName string, chr rune) *draw2d.Glyph {
	if cache.glyphMap[fontName] == nil {
		cache.glyphMap[fontName] = make(map[rune]*draw2d.Glyph, 60)
	}
	if cache.glyphMap[fontName][chr] == nil {
		cache.glyphMap[fontName][chr] = cache.Render(gc, fontName, chr)
	}
	return cache.glyphMap[fontName][chr].Copy()
}

// Render renders a glyph then returns it
func (cache *defaultGlyphCache) Render(gc draw2d.GraphicContext, fontName string, chr rune) *draw2d.Glyph {
	gc.Save()
	defer gc.Restore()
	gc.BeginPath()
	width := gc.CreateStringPath(string(chr), 0, 0)
	path := gc.GetPath()
	return &draw2d.Glyph{
		Path:  &path,
		Width: width,
	}
}

// FetchGlyph fetches a glyph from the cache, calling renderGlyph first if it doesn't already exist
func FetchGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *draw2d.Glyph {
	return gc.GetGlyphCache().Fetch(gc, fontName, chr)
}

// renderGlyph renders a glyph then caches and returns it
func renderGlyph(gc draw2d.GraphicContext, fontName string, chr rune) *draw2d.Glyph {
	return gc.GetGlyphCache().Render(gc, fontName, chr)
}
