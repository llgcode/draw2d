// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dimg

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/llgcode/draw2d/draw2dkit"
)

func TestNewGraphicContext_RGBA(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	if gc == nil {
		t.Error("NewGraphicContext should not return nil for RGBA image")
	}
}

func TestNewGraphicContext_UnsupportedImageType(t *testing.T) {
	// Test related to issue #143: Unsupported image types should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewGraphicContext should panic for unsupported image type")
		}
	}()
	img := image.NewPaletted(image.Rect(0, 0, 100, 100), nil)
	NewGraphicContext(img)
}

func TestGraphicContext_Clear(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.NRGBA{255, 0, 0, 255})
	gc.Clear()
	// Check that a pixel has the fill color
	r, g, b, a := img.At(50, 50).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Errorf("Clear should fill with fill color, got RGBA(%d, %d, %d, %d)", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestGraphicContext_ClearRect(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	// Clear entire image with white
	gc.SetFillColor(color.White)
	gc.Clear()
	// Clear a rect with red
	gc.SetFillColor(color.NRGBA{255, 0, 0, 255})
	gc.ClearRect(10, 10, 20, 20)
	// Check inside the rect
	r, g, b, _ := img.At(15, 15).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Error("ClearRect should fill rect area with fill color")
	}
	// Check outside the rect
	r2, g2, b2, _ := img.At(50, 50).RGBA()
	if r2>>8 != 255 || g2>>8 != 255 || b2>>8 != 255 {
		t.Error("ClearRect should not affect area outside rect")
	}
}

func TestGraphicContext_GetDPI(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	if gc.GetDPI() != 92 {
		t.Errorf("Default DPI = %d, want 92", gc.GetDPI())
	}
}

func TestGraphicContext_SetDPI(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetDPI(150)
	if gc.GetDPI() != 150 {
		t.Errorf("SetDPI(150): DPI = %d, want 150", gc.GetDPI())
	}
}

func TestGraphicContext_StrokeRectangle(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()
	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(2)
	draw2dkit.Rectangle(gc, 10, 10, 90, 90)
	gc.Stroke()
	// Check that an edge pixel is colored
	r, g, b, _ := img.At(10, 10).RGBA()
	if r>>8 == 255 && g>>8 == 255 && b>>8 == 255 {
		t.Error("StrokeRectangle should draw on edges")
	}
}

func TestGraphicContext_FillRectangle(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()
	gc.SetFillColor(color.NRGBA{0, 255, 0, 255})
	draw2dkit.Rectangle(gc, 20, 20, 80, 80)
	gc.Fill()
	// Check that an inside pixel is colored
	r, g, b, _ := img.At(50, 50).RGBA()
	if r>>8 != 0 || g>>8 != 255 || b>>8 != 0 {
		t.Errorf("FillRectangle should fill interior, got RGB(%d, %d, %d)", r>>8, g>>8, b>>8)
	}
}

func TestGraphicContext_FillStroke(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()
	gc.SetFillColor(color.NRGBA{0, 255, 0, 255})
	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(2)
	draw2dkit.Rectangle(gc, 20, 20, 80, 80)
	gc.FillStroke()
	// Check that interior is filled
	r, g, b, _ := img.At(50, 50).RGBA()
	if r>>8 != 0 || g>>8 != 255 || b>>8 != 0 {
		t.Error("FillStroke should fill interior")
	}
}

func TestGraphicContext_SaveRestoreColors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetFillColor(color.NRGBA{0, 255, 0, 255})
	gc.Save()
	gc.SetStrokeColor(color.NRGBA{0, 0, 255, 255})
	gc.SetFillColor(color.NRGBA{255, 255, 0, 255})
	gc.Restore()
	// Check that colors are restored
	r1, g1, b1, _ := gc.Current.StrokeColor.RGBA()
	if r1>>8 != 255 || g1>>8 != 0 || b1>>8 != 0 {
		t.Error("Restore should restore StrokeColor")
	}
	r2, g2, b2, _ := gc.Current.FillColor.RGBA()
	if r2>>8 != 0 || g2>>8 != 255 || b2>>8 != 0 {
		t.Error("Restore should restore FillColor")
	}
}

func TestGraphicContext_SetLineCap(t *testing.T) {
	// Test related to issue #155: LineCap should be stored correctly
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetLineCap(0) // ButtCap
	if gc.Current.Cap != 0 {
		t.Errorf("SetLineCap(ButtCap): got %v, want ButtCap", gc.Current.Cap)
	}
	gc.SetLineCap(1) // RoundCap
	if gc.Current.Cap != 1 {
		t.Errorf("SetLineCap(RoundCap): got %v, want RoundCap", gc.Current.Cap)
	}
	gc.SetLineCap(2) // SquareCap
	if gc.Current.Cap != 2 {
		t.Errorf("SetLineCap(SquareCap): got %v, want SquareCap", gc.Current.Cap)
	}
}

func TestGraphicContext_DrawImage(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DrawImage panicked: %v", r)
		}
	}()
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	srcImg := image.NewRGBA(image.Rect(0, 0, 50, 50))
	gc.DrawImage(srcImg)
}

func TestSaveToPngFile(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	tmpFile := t.TempDir() + "/test.png"
	err := SaveToPngFile(tmpFile, img)
	if err != nil {
		t.Errorf("SaveToPngFile failed: %v", err)
	}
	// Verify file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("SaveToPngFile did not create file")
	}
}

func TestGraphicContext_TransformAffectsDrawing(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()
	gc.Translate(10, 10)
	gc.SetFillColor(color.NRGBA{255, 0, 0, 255})
	draw2dkit.Rectangle(gc, 0, 0, 10, 10)
	gc.Fill()
	// Rectangle should be at (10,10) due to translation
	r, g, b, _ := img.At(15, 15).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Error("Transform should affect drawing position")
	}
}

func TestGraphicContext_LineWidth(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetLineWidth(5.0)
	if gc.Current.LineWidth != 5.0 {
		t.Errorf("SetLineWidth(5.0): got %f, want 5.0", gc.Current.LineWidth)
	}
}

func TestGraphicContext_Circle_Fill(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	gc := NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.Clear()
	gc.SetFillColor(color.NRGBA{0, 0, 255, 255})
	draw2dkit.Circle(gc, 50, 50, 20)
	gc.Fill()
	// Check that center pixel is colored
	r, g, b, _ := img.At(50, 50).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 {
		t.Errorf("Circle Fill should color center pixel, got RGB(%d, %d, %d)", r>>8, g>>8, b>>8)
	}
}
