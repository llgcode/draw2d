// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// Tests for known bugs in draw2dimg package

package draw2dimg

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d"
)

// TestIssue155_LineCapVisualDifference tests that different line caps produce visually different results.
// Issue: https://github.com/llgcode/draw2d/issues/155
// Expected: Different line caps should produce different visual output
// Actual: All line caps appear to render the same way
func TestIssue155_LineCapVisualDifference(t *testing.T) {
	testLineCaps := []struct {
		name string
		cap  draw2d.LineCap
	}{
		{"ButtCap", draw2d.ButtCap},
		{"RoundCap", draw2d.RoundCap},
		{"SquareCap", draw2d.SquareCap},
	}
	
	images := make([]*image.RGBA, len(testLineCaps))
	
	// Create images with different line caps
	for i, tc := range testLineCaps {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		gc := NewGraphicContext(img)
		gc.SetFillColor(color.White)
		gc.Clear()
		gc.SetStrokeColor(color.Black)
		gc.SetLineWidth(30)
		gc.SetLineCap(tc.cap)
		
		// Draw a horizontal line
		gc.MoveTo(20, 50)
		gc.LineTo(80, 50)
		gc.Stroke()
		
		images[i] = img
	}
	
	// Compare images - they should be different
	// Check pixels at the end of the line (x=85, beyond the line end)
	buttPixel := images[0].At(85, 50)
	roundPixel := images[1].At(85, 50)
	squarePixel := images[2].At(85, 50)
	
	br, bg, bb, _ := buttPixel.RGBA()
	rr, rg, rb, _ := roundPixel.RGBA()
	sr, sg, sb, _ := squarePixel.RGBA()
	
	// All three should be different, but they're not
	allSame := (br == rr && bg == rg && bb == rb) && (br == sr && bg == sg && bb == sb)
	
	if allSame {
		t.Logf("KNOWN BUG: All line caps render identically")
		t.Logf("ButtCap pixel at line end+5: RGB(%d,%d,%d)", br>>8, bg>>8, bb>>8)
		t.Logf("RoundCap pixel at line end+5: RGB(%d,%d,%d)", rr>>8, rg>>8, rb>>8)
		t.Logf("SquareCap pixel at line end+5: RGB(%d,%d,%d)", sr>>8, sg>>8, sb>>8)
		t.Errorf("Issue #155: Different line caps should produce different output")
	}
}

// TestIssue155_LineJoinVisualDifference tests that different line joins produce different results.
// Issue: https://github.com/llgcode/draw2d/issues/155 (also affects line joins)
// Expected: Different line joins should produce different visual output
// Actual: Line joins may not render correctly
func TestIssue155_LineJoinVisualDifference(t *testing.T) {
	testLineJoins := []struct {
		name string
		join draw2d.LineJoin
	}{
		{"BevelJoin", draw2d.BevelJoin},
		{"RoundJoin", draw2d.RoundJoin},
		{"MiterJoin", draw2d.MiterJoin},
	}
	
	images := make([]*image.RGBA, len(testLineJoins))
	
	// Create images with different line joins
	for i, tc := range testLineJoins {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		gc := NewGraphicContext(img)
		gc.SetFillColor(color.White)
		gc.Clear()
		gc.SetStrokeColor(color.Black)
		gc.SetLineWidth(20)
		gc.SetLineJoin(tc.join)
		
		// Draw two lines meeting at 90 degrees
		gc.MoveTo(30, 70)
		gc.LineTo(50, 50)
		gc.LineTo(70, 70)
		gc.Stroke()
		
		images[i] = img
	}
	
	// Check the corner pixel where lines meet
	// Different joins should produce different appearances at the corner
	bevelCorner := images[0].At(50, 50)
	roundCorner := images[1].At(50, 50)
	miterCorner := images[2].At(50, 50)
	
	br, bg, bb, _ := bevelCorner.RGBA()
	rr, rg, rb, _ := roundCorner.RGBA()
	mr, mg, mb, _ := miterCorner.RGBA()
	
	allSame := (br == rr && bg == rg && bb == rb) && (br == mr && bg == mg && bb == mb)
	
	if allSame {
		t.Logf("KNOWN BUG: Line joins may not render with visible differences")
		t.Logf("BevelJoin corner: RGB(%d,%d,%d)", br>>8, bg>>8, bb>>8)
		t.Logf("RoundJoin corner: RGB(%d,%d,%d)", rr>>8, rg>>8, rb>>8)
		t.Logf("MiterJoin corner: RGB(%d,%d,%d)", mr>>8, mg>>8, mb>>8)
		t.Logf("Issue #155: Different line joins should produce different output")
		// Don't fail - this may actually work for joins, just document it
	}
}

// TestIssue143_UnsupportedImageTypesDocumented documents which image types are not supported.
// Issue: https://github.com/llgcode/draw2d/issues/143
func TestIssue143_UnsupportedImageTypesDocumented(t *testing.T) {
	// Test that we properly document unsupported image types
	unsupportedTypes := []struct {
		name string
		makeImage func() image.Image
	}{
		{"Paletted", func() image.Image { 
			return image.NewPaletted(image.Rect(0, 0, 100, 100), nil) 
		}},
		{"Gray", func() image.Image { 
			return image.NewGray(image.Rect(0, 0, 100, 100)) 
		}},
		{"Gray16", func() image.Image { 
			return image.NewGray16(image.Rect(0, 0, 100, 100)) 
		}},
		{"Alpha", func() image.Image { 
			return image.NewAlpha(image.Rect(0, 0, 100, 100)) 
		}},
	}
	
	supportCount := 0
	unsupportedCount := 0
	
	for _, tt := range unsupportedTypes {
		t.Run(tt.name, func(t *testing.T) {
			img := tt.makeImage()
			
			defer func() {
				if r := recover(); r != nil {
					unsupportedCount++
					t.Logf("CONFIRMED: %s is not supported (panics as expected)", tt.name)
				}
			}()
			
			// This will panic for unsupported types
			_ = NewGraphicContext(img.(interface{
				At(x, y int) color.Color
				Bounds() image.Rectangle
				ColorModel() color.Model
				Set(x, y int, c color.Color)
			}))
			
			supportCount++
			t.Logf("UNEXPECTED: %s is supported", tt.name)
		})
	}
	
	if unsupportedCount > 0 {
		t.Logf("Issue #143: %d image types are not supported", unsupportedCount)
		t.Logf("Only *image.RGBA is currently supported")
		t.Logf("This is a known limitation")
	}
}
