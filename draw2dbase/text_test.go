// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"github.com/llgcode/draw2d"
	"testing"
)

func TestNewGlyphCache(t *testing.T) {
	cache := NewGlyphCache()
	if cache == nil {
		t.Error("NewGlyphCache() returned nil")
	}

	if cache.glyphs == nil {
		t.Error("NewGlyphCache() glyphs map is nil")
	}
}

func TestGlyph_Copy(t *testing.T) {
	// Create a path for the glyph
	path := &draw2d.Path{}
	path.MoveTo(0, 0)
	path.LineTo(10, 10)
	path.LineTo(20, 0)
	path.Close()

	// Create original glyph
	original := &Glyph{
		Path:  path,
		Width: 100.0,
	}

	// Copy the glyph
	copy := original.Copy()

	// Verify the copy is not nil
	if copy == nil {
		t.Fatal("Glyph.Copy() returned nil")
	}

	// Verify independence - modifying copy should not affect original
	if copy.Path == original.Path {
		t.Error("Glyph.Copy() did not create independent Path copy")
	}

	// Verify width is preserved
	if copy.Width != original.Width {
		t.Errorf("Glyph.Copy() Width = %v, want %v", copy.Width, original.Width)
	}
}

func TestGlyph_Copy_Width(t *testing.T) {
	tests := []struct {
		name  string
		width float64
	}{
		{"Zero Width", 0.0},
		{"Small Width", 10.5},
		{"Large Width", 1000.0},
		{"Negative Width", -5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := &draw2d.Path{}
			path.MoveTo(0, 0)
			path.LineTo(10, 0)

			original := &Glyph{
				Path:  path,
				Width: tt.width,
			}

			copy := original.Copy()

			if copy.Width != tt.width {
				t.Errorf("Glyph.Copy() Width = %v, want %v", copy.Width, tt.width)
			}
		})
	}
}
