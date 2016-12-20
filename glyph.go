package draw2d

// Types implementing this interface can be passed to gc.SetGlyphCache to change the
// way glyphs are being stored and retrieved.
type GlyphCache interface {
	// Fetch fetches a glyph from the cache, storing with Render first if it doesn't already exist
	Fetch(gc GraphicContext, fontName string, chr rune) *Glyph

	// Render renders a glyph then returns it
	Render(gc GraphicContext, fontName string, chr rune) *Glyph
}

// Glyph represents a rune which has been converted to a Path and width
type Glyph struct {
	// Path represents a glyph, it is always at (0, 0)
	Path *Path
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
func (g *Glyph) Fill(gc GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Fill(g.Path)
	gc.Restore()
	return g.Width
}

// Stroke strokes a Glyph in the specified coordinates
func (g *Glyph) Stroke(gc GraphicContext, x, y float64) float64 {
	gc.Save()
	gc.BeginPath()
	gc.Translate(x, y)
	gc.Stroke(g.Path)
	gc.Restore()
	return g.Width
}
