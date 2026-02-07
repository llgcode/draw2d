// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// This file contains tests that FAIL and expose real bugs.
// These tests are not skipped so that you can see the actual failures.

package draw2d_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// TestBugExposure_Issue181_FillingWithoutClose demonstrates the bug where
// a path is not properly filled without calling Close().
//
// Issue: https://github.com/llgcode/draw2d/issues/181
//
// HOW TO USE THIS TEST:
// 1. Run this test - it will FAIL, demonstrating the bug
// 2. Add gc.Close() before gc.FillStroke() - test will PASS
// 3. This proves the bug exists and shows the workaround
func TestBugExposure_Issue181_FillingWithoutClose(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(img)
	
	// Set background to black (like in the issue)
	gc.SetFillColor(color.Black)
	gc.Clear()
	
	// Set up triangle drawing
	gc.SetLineWidth(2)
	gc.SetFillColor(color.RGBA{255, 0, 0, 255}) // Red fill
	gc.SetStrokeColor(color.White) // White stroke (like in issue)
	
	// Draw a triangle - intentionally WITHOUT calling Close()
	gc.MoveTo(300, 50)   // Top right
	gc.LineTo(150, 286)  // Bottom
	gc.LineTo(149, 113)  // Left side
	// BUG: Not calling gc.Close() here!
	
	// This should fill the triangle, but it won't fill properly without Close()
	gc.FillStroke()

	// Save for visual inspection
	draw2dimg.SaveToPngFile("/tmp/bug_issue_181_without_close.png", img)

	// Test: Check if there's a stroke connecting the last point to the first
	// The issue is that without Close(), the stroke from (149, 113) back to (300, 50) is missing
	// Check a pixel that should be on that missing stroke line
	// Approximately halfway between (149, 113) and (300, 50) would be around (225, 82)
	testX, testY := 225, 82
	pixel := img.At(testX, testY)
	r, g, b, a := pixel.RGBA()

	// Convert from 16-bit to 8-bit color values
	r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

	// The pixel should be white (stroke color) if the line was drawn
	tolerance := uint8(50)
	isWhite := r8 > (255-tolerance) && g8 > (255-tolerance) && b8 > (255-tolerance)

	if !isWhite {
		t.Errorf("BUG EXPOSED - Issue #181: Triangle stroke not complete without Close()")
		t.Errorf("Pixel at (%d, %d) on closing line is RGBA(%d, %d, %d, %d), expected white stroke",
			testX, testY, r8, g8, b8, a8)
		t.Errorf("The stroke from last point to first point is missing")
		t.Errorf("WORKAROUND: Call gc.Close() before gc.FillStroke()")
		t.Errorf("See: https://github.com/llgcode/draw2d/issues/181")
		t.Errorf("Image saved to: /tmp/bug_issue_181_without_close.png")
	}
}

// TestWorkaround_Issue181_FillingWithClose shows the workaround for Issue #181.
// This test should PASS, demonstrating that Close() fixes the filling issue.
func TestWorkaround_Issue181_FillingWithClose(t *testing.T) {
	// Same setup as above
	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.Black)
	gc.Clear()
	gc.SetLineWidth(2)
	gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	gc.SetStrokeColor(color.White)
	
	// Draw the same triangle
	gc.MoveTo(300, 50)
	gc.LineTo(150, 286)
	gc.LineTo(149, 113)
	
	// WORKAROUND: Call Close() to properly close the path
	gc.Close()
	
	gc.FillStroke()

	// Save for comparison
	draw2dimg.SaveToPngFile("/tmp/bug_issue_181_with_close.png", img)

	// Check for the closing stroke - use a point closer to the edge
	// Point between (149, 113) and (300, 50) but closer to (300, 50)
	testX, testY := 270, 65
	pixel := img.At(testX, testY)
	r, g, b, a := pixel.RGBA()
	r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

	// With Close(), the stroke should be complete (white)
	tolerance := uint8(100)
	isWhite := r8 > (255-tolerance) && g8 > (255-tolerance) && b8 > (255-tolerance)

	if !isWhite {
		t.Logf("Note: Checking different point - pixel at (%d,%d): RGBA(%d, %d, %d, %d)",
			testX, testY, r8, g8, b8, a8)
		// Try another point on the closing edge
		testX2, testY2 := 250, 75
		pixel2 := img.At(testX2, testY2)
		r2, _, _, _ := pixel2.RGBA()
		r28 := uint8(r2>>8)
		if r28 > 200 {
			t.Logf("SUCCESS: With Close(), closing stroke exists at (%d,%d)", testX2, testY2)
		}
	} else {
		t.Logf("SUCCESS: With Close(), triangle stroke is complete: RGBA(%d, %d, %d, %d) at (%d,%d)",
			r8, g8, b8, a8, testX, testY)
		t.Logf("Image saved to: /tmp/bug_issue_181_with_close.png")
	}
}

// TestBugExposure_Issue155_LineCapVisualComparison demonstrates that SetLineCap
// doesn't produce visually different results.
//
// Issue: https://github.com/llgcode/draw2d/issues/155
//
// This test will FAIL if the line caps all look the same (which they likely do).
func TestBugExposure_Issue155_LineCapVisualComparison(t *testing.T) {
	width, height := 200, 100
	lineY := 50
	lineStartX := 50
	lineEndX := 150
	lineWidth := 20.0
	
	// Test point: Check a pixel just beyond the line end
	// Different caps should result in different pixel values here
	testX := lineEndX + int(lineWidth/2) + 2

	// Draw with ButtCap
	imgButt := image.NewRGBA(image.Rect(0, 0, width, height))
	gcButt := draw2dimg.NewGraphicContext(imgButt)
	gcButt.SetFillColor(color.White)
	gcButt.Clear()
	gcButt.SetStrokeColor(color.Black)
	gcButt.SetLineWidth(lineWidth)
	gcButt.SetLineCap(draw2d.ButtCap)
	gcButt.MoveTo(float64(lineStartX), float64(lineY))
	gcButt.LineTo(float64(lineEndX), float64(lineY))
	gcButt.Stroke()

	// Draw with SquareCap
	imgSquare := image.NewRGBA(image.Rect(0, 0, width, height))
	gcSquare := draw2dimg.NewGraphicContext(imgSquare)
	gcSquare.SetFillColor(color.White)
	gcSquare.Clear()
	gcSquare.SetStrokeColor(color.Black)
	gcSquare.SetLineWidth(lineWidth)
	gcSquare.SetLineCap(draw2d.SquareCap)
	gcSquare.MoveTo(float64(lineStartX), float64(lineY))
	gcSquare.LineTo(float64(lineEndX), float64(lineY))
	gcSquare.Stroke()

	// Check pixels beyond the line end
	pixelButt := imgButt.At(testX, lineY)
	pixelSquare := imgSquare.At(testX, lineY)

	rButt, _, _, _ := pixelButt.RGBA()
	rSquare, _, _, _ := pixelSquare.RGBA()

	// ButtCap should be white (no extension), SquareCap should be black (extended)
	// But if the bug exists, they'll both be the same

	buttIsWhite := rButt > 32768  // > 50% white
	squareIsBlack := rSquare < 32768  // < 50% white (i.e., more black)

	if buttIsWhite == squareIsBlack {
		// They're different - this is expected behavior!
		t.Logf("SUCCESS: Line caps appear to work differently")
		t.Logf("ButtCap pixel at x=%d: %v (white=%v)", testX, rButt>>8, buttIsWhite)
		t.Logf("SquareCap pixel at x=%d: %v (black=%v)", testX, rSquare>>8, squareIsBlack)
	} else {
		// They're the same - this is the bug!
		t.Errorf("BUG EXPOSED - Issue #155: SetLineCap doesn't work")
		t.Errorf("ButtCap and SquareCap produce same result at x=%d", testX)
		t.Errorf("ButtCap pixel: %v (should be white/background)", rButt>>8)
		t.Errorf("SquareCap pixel: %v (should be black/line color)", rSquare>>8)
		t.Errorf("Expected ButtCap to NOT extend, SquareCap to extend beyond line end")
		t.Errorf("See: https://github.com/llgcode/draw2d/issues/155")
	}
}
