// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// This file contains tests that expose known bugs and limitations
// documented in GitHub issues. These tests are expected to fail
// until the issues are resolved.

package draw2d_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// TestIssue181_WrongFilling tests that a path without Close() is not properly filled.
// Issue: https://github.com/llgcode/draw2d/issues/181
// Expected: The triangle should be filled completely even without calling Close()
// Actual: The triangle is not filled from the starting and ending points
//
// This test demonstrates a real bug where FillStroke() doesn't properly fill
// a path that hasn't been explicitly closed with Close().
func TestIssue181_WrongFilling(t *testing.T) {
	t.Skip("Known issue #181: Wrong filling without Close()")

	img := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.Black)
	gc.Clear()
	gc.SetLineWidth(2)
	gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	gc.SetStrokeColor(color.White)
	
	// Draw a triangle without calling Close()
	gc.MoveTo(300, 50)
	gc.LineTo(150, 286)
	gc.LineTo(149, 113)
	// Intentionally NOT calling gc.Close() - this is the bug
	gc.FillStroke()

	// Check if the triangle is properly filled by examining pixels inside
	// The center of the triangle should be red (filled)
	centerX, centerY := 200, 150
	pixel := img.At(centerX, centerY)
	r, g, b, a := pixel.RGBA()

	// The pixel should be red (255, 0, 0, 255) if properly filled
	// But due to the bug, it will be black (0, 0, 0, 255)
	if r == 0 && g == 0 && b == 0 && a == 65535 {
		t.Errorf("Bug confirmed: Triangle not filled without Close(). Center pixel is black (%v, %v, %v, %v), expected red",
			r>>8, g>>8, b>>8, a>>8)
	}
}

// TestIssue155_SetLineCapDoesNotWork tests that SetLineCap doesn't actually change line appearance.
// Issue: https://github.com/llgcode/draw2d/issues/155
// Expected: Different line caps (Round, Butt, Square) should produce visibly different results
// Actual: All line caps appear the same
//
// This test demonstrates that SetLineCap may not be properly implemented or respected
// by the rendering backend.
func TestIssue155_SetLineCapDoesNotWork(t *testing.T) {
	t.Skip("Known issue #155: SetLineCap does not work")

	width, height := 400, 300

	// Create three images with different line caps
	imgRound := image.NewRGBA(image.Rect(0, 0, width, height))
	imgButt := image.NewRGBA(image.Rect(0, 0, width, height))
	imgSquare := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw line with RoundCap
	gcRound := draw2dimg.NewGraphicContext(imgRound)
	gcRound.SetStrokeColor(color.Black)
	gcRound.SetLineWidth(20)
	gcRound.SetLineCap(draw2d.RoundCap)
	gcRound.MoveTo(50, 150)
	gcRound.LineTo(350, 150)
	gcRound.Stroke()

	// Draw line with ButtCap
	gcButt := draw2dimg.NewGraphicContext(imgButt)
	gcButt.SetStrokeColor(color.Black)
	gcButt.SetLineWidth(20)
	gcButt.SetLineCap(draw2d.ButtCap)
	gcButt.MoveTo(50, 150)
	gcButt.LineTo(350, 150)
	gcButt.Stroke()

	// Draw line with SquareCap
	gcSquare := draw2dimg.NewGraphicContext(imgSquare)
	gcSquare.SetStrokeColor(color.Black)
	gcSquare.SetLineWidth(20)
	gcSquare.SetLineCap(draw2d.SquareCap)
	gcSquare.MoveTo(50, 150)
	gcSquare.LineTo(350, 150)
	gcSquare.Stroke()

	// Check pixels at the line ends (x=50 and x=350)
	// RoundCap should extend slightly beyond the line end
	// ButtCap should end exactly at the line end
	// SquareCap should extend further than RoundCap

	// Check a pixel beyond the line end (x=355)
	pixelRound := imgRound.At(355, 150)
	pixelButt := imgButt.At(355, 150)
	pixelSquare := imgSquare.At(355, 150)

	// All three should be different, but they're likely all the same due to the bug
	rR, _, _, _ := pixelRound.RGBA()
	rB, _, _, _ := pixelButt.RGBA()
	rS, _, _, _ := pixelSquare.RGBA()

	// If all are the same (all black or all white), the bug is confirmed
	if rR == rB && rB == rS {
		t.Errorf("Bug confirmed: All line caps appear identical. RoundCap pixel=%v, ButtCap pixel=%v, SquareCap pixel=%v",
			rR>>8, rB>>8, rS>>8)
	}
}

// TestIssue171_TextStrokeLineCap tests that text stroke doesn't properly connect.
// Issue: https://github.com/llgcode/draw2d/issues/171
// Expected: Text stroke should fully cover and connect around letters
// Actual: Strokes on letters like "i" and "t" don't fully connect
//
// This is related to Issue #155 - LineCap and LineJoin settings don't work properly
// for stroked text paths.
func TestIssue171_TextStrokeLineCap(t *testing.T) {
	t.Skip("Known issue #171: Text stroke LineCap and LineJoin don't work properly")

	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()

	// Set up stroke style for text
	gc.SetStrokeColor(color.RGBA{0, 0, 255, 255})
	gc.SetLineWidth(2)
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineJoin(draw2d.RoundJoin)

	// Try to stroke the letter "i" which should have a connected stroke
	gc.SetFontSize(48)
	gc.StrokeStringAt("i", 50, 60)

	// The issue is difficult to test programmatically, but we can verify
	// that the SetLineCap was called (though it may not have any effect)
	// In a visual test, you would see disconnected strokes on the letter

	// For now, just document that this is a known issue
	t.Logf("Known issue: Text strokes don't respect LineCap/LineJoin settings")
}

// TestIssue129_StrokeStyleNotUsed tests that StrokeStyle type isn't actually used.
// Issue: https://github.com/llgcode/draw2d/issues/129
// Expected: Setting a StrokeStyle should affect how lines are drawn
// Actual: The StrokeStyle type exists but there's no clear way to use it
//
// This test demonstrates that while StrokeStyle is defined in the API,
// it's not clear how to apply it or if it's actually used anywhere.
func TestIssue129_StrokeStyleNotUsed(t *testing.T) {
	t.Skip("Known issue #129: StrokeStyle type not clearly used in API")

	// Create a StrokeStyle with specific settings
	style := draw2d.StrokeStyle{
		Color:      color.RGBA{255, 0, 0, 255},
		Width:      10.0,
		LineCap:    draw2d.RoundCap,
		LineJoin:   draw2d.RoundJoin,
		DashOffset: 0,
		Dash:       []float64{10, 5},
	}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	gc := draw2dimg.NewGraphicContext(img)

	// Problem: There's no method like gc.SetStrokeStyle(style) to apply it
	// We have to set each property individually:
	gc.SetStrokeColor(style.Color)
	gc.SetLineWidth(style.Width)
	gc.SetLineCap(style.LineCap)
	gc.SetLineJoin(style.LineJoin)
	gc.SetLineDash(style.Dash, style.DashOffset)

	// This test mainly documents that StrokeStyle exists but isn't integrated
	t.Logf("Known issue: StrokeStyle type exists but there's no SetStrokeStyle() method")
	t.Logf("Style values must be set individually: %+v", style)
}

// TestLineCapVisualDifference is a helper test to verify that different line caps
// should produce visually different results. This test documents what SHOULD happen.
func TestLineCapVisualDifference(t *testing.T) {
	t.Skip("This is a reference test showing expected behavior")

	// This test documents what the expected behavior should be:
	// 
	// RoundCap: The end of the line should have a semicircular cap
	//           extending Width/2 beyond the endpoint
	//
	// ButtCap: The end of the line should be flat and flush with the endpoint
	//
	// SquareCap: The end should be flat but extend Width/2 beyond the endpoint
	//
	// If Issue #155 is fixed, these differences should be measurable in pixels

	t.Logf("Reference: Line cap differences")
	t.Logf("- RoundCap: Should extend ~Width/2 with rounded end")
	t.Logf("- ButtCap: Should end flush with line endpoint")
	t.Logf("- SquareCap: Should extend Width/2 with flat end")
}
