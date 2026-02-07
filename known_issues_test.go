// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// This file contains tests for known bugs and limitations tracked in GitHub issues.
// These tests are expected to FAIL and demonstrate real problems with the current implementation.
// Each test is documented with the issue number and describes the expected vs actual behavior.

package draw2d_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dpdf"
)

// TestIssue181_TriangleFillingWithoutClose tests the bug where a triangle
// doesn't fill properly when Close() is not called.
// Issue: https://github.com/llgcode/draw2d/issues/181
// Expected: Triangle should be filled even without explicit Close()
// Actual: Triangle is not filled from starting to ending points
func TestIssue181_TriangleFillingWithoutClose(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(img)
	
	// Setup
	gc.SetFillColor(color.Black)
	gc.Clear()
	gc.SetLineWidth(2)
	gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	gc.SetStrokeColor(color.White)
	
	// Draw triangle WITHOUT calling Close()
	gc.MoveTo(300, 50)
	gc.LineTo(150, 286)
	gc.LineTo(149, 113)
	// Intentionally NOT calling gc.Close() - this is the bug
	
	gc.FillStroke()
	
	// Check that the triangle interior is filled
	// The center of the triangle should be red
	centerX, centerY := 200, 150
	r, g, b, _ := img.At(centerX, centerY).RGBA()
	
	// Expected: center should be red (255, 0, 0)
	// Actual: center is NOT red because path is not closed
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Logf("KNOWN BUG: Triangle without Close() doesn't fill properly")
		t.Logf("Center pixel (%d, %d) should be red (255,0,0) but got RGB(%d,%d,%d)", 
			centerX, centerY, r>>8, g>>8, b>>8)
		t.Errorf("Issue #181: Triangle should be filled even without explicit Close()")
	}
}

// TestIssue155_SetLineCapButtCap tests whether different line caps are actually rendered.
// Issue: https://github.com/llgcode/draw2d/issues/155
// Expected: ButtCap should render differently than RoundCap
// Actual: Line caps appear to render the same way
func TestIssue155_SetLineCapButtCap(t *testing.T) {
	// Create two images with different line caps
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc1 := draw2dimg.NewGraphicContext(img1)
	gc1.SetFillColor(color.White)
	gc1.Clear()
	gc1.SetStrokeColor(color.Black)
	gc1.SetLineWidth(20)
	gc1.SetLineCap(draw2d.ButtCap)
	gc1.MoveTo(50, 20)
	gc1.LineTo(50, 80)
	gc1.Stroke()
	
	img2 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc2 := draw2dimg.NewGraphicContext(img2)
	gc2.SetFillColor(color.White)
	gc2.Clear()
	gc2.SetStrokeColor(color.Black)
	gc2.SetLineWidth(20)
	gc2.SetLineCap(draw2d.RoundCap)
	gc2.MoveTo(50, 20)
	gc2.LineTo(50, 80)
	gc2.Stroke()
	
	// Check the end points - they should be different
	// For ButtCap, the line should end exactly at y=80
	// For RoundCap, the line should extend beyond y=80
	
	// Check if the pixel just beyond the end is different
	buttEndPixel := img1.At(50, 90)
	roundEndPixel := img2.At(50, 90)
	
	br, bg, bb, _ := buttEndPixel.RGBA()
	rr, rg, rb, _ := roundEndPixel.RGBA()
	
	// Expected: RoundCap extends beyond line end, ButtCap does not
	// So roundEndPixel should be darker than buttEndPixel
	// Actual: They are the same (LineCap doesn't work)
	if br == rr && bg == rg && bb == rb {
		t.Logf("KNOWN BUG: SetLineCap doesn't produce different rendering")
		t.Logf("ButtCap pixel at end+10: RGB(%d,%d,%d)", br>>8, bg>>8, bb>>8)
		t.Logf("RoundCap pixel at end+10: RGB(%d,%d,%d)", rr>>8, rg>>8, rb>>8)
		t.Errorf("Issue #155: ButtCap and RoundCap render identically")
	}
}

// TestIssue155_SetLineCapSquareCap tests whether SquareCap renders differently.
// Issue: https://github.com/llgcode/draw2d/issues/155
func TestIssue155_SetLineCapSquareCap(t *testing.T) {
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc1 := draw2dimg.NewGraphicContext(img1)
	gc1.SetFillColor(color.White)
	gc1.Clear()
	gc1.SetStrokeColor(color.Black)
	gc1.SetLineWidth(20)
	gc1.SetLineCap(draw2d.ButtCap)
	gc1.MoveTo(50, 30)
	gc1.LineTo(50, 70)
	gc1.Stroke()
	
	img2 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc2 := draw2dimg.NewGraphicContext(img2)
	gc2.SetFillColor(color.White)
	gc2.Clear()
	gc2.SetStrokeColor(color.Black)
	gc2.SetLineWidth(20)
	gc2.SetLineCap(draw2d.SquareCap)
	gc2.MoveTo(50, 30)
	gc2.LineTo(50, 70)
	gc2.Stroke()
	
	// SquareCap should extend the line by half the line width (10 pixels)
	// Check if pixel beyond the end is different
	buttPixel := img1.At(50, 80)
	squarePixel := img2.At(50, 80)
	
	br, bg, bb, _ := buttPixel.RGBA()
	sr, sg, sb, _ := squarePixel.RGBA()
	
	if br == sr && bg == sg && bb == sb {
		t.Logf("KNOWN BUG: SetLineCap(SquareCap) doesn't produce different rendering")
		t.Errorf("Issue #155: ButtCap and SquareCap render identically")
	}
}

// TestIssue139_PDFVerticalFlip tests Y-axis flipping with PDF backend.
// Issue: https://github.com/llgcode/draw2d/issues/139
// Expected: Y-axis flip should work with PDF backend like it does with image backend
// Actual: Scale(1, -1) silently fails with draw2dpdf.GraphicContext
func TestIssue139_PDFVerticalFlip(t *testing.T) {
	// Create a PDF context
	pdf := draw2dpdf.NewPdf("L", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(pdf)
	
	// Try to flip Y axis
	gc.Save()
	gc.Translate(0, 100)
	gc.Scale(1, -1)
	
	// Draw a simple rectangle
	gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	gc.MoveTo(10, 10)
	gc.LineTo(50, 10)
	gc.LineTo(50, 30)
	gc.LineTo(10, 30)
	gc.Close()
	gc.Fill()
	
	gc.Restore()
	
	// We can't easily verify the PDF output in a unit test, but we can check
	// that the transformation matrix was set
	m := gc.GetMatrixTransform()
	
	// Expected: m[3] should be -1 (Y scale factor)
	// Actual: Transformation may not be applied properly to PDF backend
	if m[3] != -1.0 {
		t.Logf("KNOWN BUG: Y-axis flip may not work properly with PDF backend")
		t.Logf("Expected matrix Y scale = -1, got: %f", m[3])
		t.Errorf("Issue #139: Scale(1, -1) doesn't work properly with draw2dpdf.GraphicContext")
	}
}

// TestIssue171_TextStrokeDisconnected tests text stroke rendering quality.
// Issue: https://github.com/llgcode/draw2d/issues/171
// Expected: Text stroke should be continuous and connected
// Actual: Text stroke has gaps and disconnections, especially for letters like 'i' and 't'
func TestIssue171_TextStrokeDisconnected(t *testing.T) {
	t.Skip("This test requires font loading and visual inspection - see issue #171")
	
	// This is a visual bug that's hard to test programmatically
	// The issue is that SetLineCap doesn't work (issue #155), which affects text stroke
	// The test would need to:
	// 1. Load a font (Roboto-Medium)
	// 2. Render text with stroke
	// 3. Check for gaps in the stroke
	
	// For now, we acknowledge this is a known issue related to #155
	t.Logf("Issue #171: Text stroke rendering has disconnections")
	t.Logf("This is related to issue #155 (SetLineCap not working)")
}

// TestIssue181_TriangleFillingWithClose verifies the workaround works.
// This test should PASS to show that calling Close() is the current workaround.
func TestIssue181_TriangleFillingWithClose(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(img)
	
	gc.SetFillColor(color.Black)
	gc.Clear()
	gc.SetLineWidth(2)
	gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	gc.SetStrokeColor(color.White)
	
	// Draw triangle WITH Close() - this should work
	gc.MoveTo(300, 50)
	gc.LineTo(150, 286)
	gc.LineTo(149, 113)
	gc.Close() // This makes it work
	
	gc.FillStroke()
	
	// Check that the triangle interior is filled
	centerX, centerY := 200, 150
	r, g, b, _ := img.At(centerX, centerY).RGBA()
	
	// This should PASS - Close() makes it work
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("With Close(), triangle should be filled. Center RGB(%d,%d,%d)", r>>8, g>>8, b>>8)
	} else {
		t.Logf("WORKAROUND VERIFIED: Calling Close() makes triangle fill work")
	}
}

// TestPerformanceNote documents the performance issue.
// Issue: https://github.com/llgcode/draw2d/issues/147
// This is not a failing test per se, but documents that draw2d is 10-30x slower than Cairo
func TestPerformanceNote(t *testing.T) {
	t.Logf("KNOWN ISSUE #147: draw2d performance is ~10-30x slower than Cairo")
	t.Logf("This is a known limitation of the current implementation")
	t.Logf("See: https://github.com/llgcode/draw2d/issues/147")
	
	// We don't fail this test, but document the limitation
	// To actually measure this, run: go test -bench=. -benchmem
}
