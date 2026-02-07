// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"image"
	"testing"

	"github.com/llgcode/draw2d"
)

func TestNewStackGraphicContext_Defaults(t *testing.T) {
	gc := NewStackGraphicContext()
	if gc.Current.LineWidth != 1.0 {
		t.Errorf("Default LineWidth = %f, want 1.0", gc.Current.LineWidth)
	}
	if gc.Current.Cap != draw2d.RoundCap {
		t.Errorf("Default Cap = %v, want RoundCap", gc.Current.Cap)
	}
	if gc.Current.Join != draw2d.RoundJoin {
		t.Errorf("Default Join = %v, want RoundJoin", gc.Current.Join)
	}
	if gc.Current.FillRule != draw2d.FillRuleEvenOdd {
		t.Errorf("Default FillRule = %v, want EvenOdd", gc.Current.FillRule)
	}
	if gc.Current.FontSize != 10 {
		t.Errorf("Default FontSize = %f, want 10", gc.Current.FontSize)
	}
	if gc.Current.FontData.Name != "luxi" {
		t.Errorf("Default FontData.Name = %s, want 'luxi'", gc.Current.FontData.Name)
	}
	if !gc.Current.Tr.IsIdentity() {
		t.Error("Default matrix should be identity")
	}
	if gc.Current.StrokeColor != image.Black {
		t.Error("Default StrokeColor should be Black")
	}
	if gc.Current.FillColor != image.White {
		t.Error("Default FillColor should be White")
	}
	if gc.Current.Path == nil {
		t.Error("Default Path should not be nil")
	}
}

func TestStackGraphicContext_SetStrokeColor(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetStrokeColor(image.White)
	if gc.Current.StrokeColor != image.White {
		t.Error("SetStrokeColor failed")
	}
}

func TestStackGraphicContext_SetFillColor(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetFillColor(image.Black)
	if gc.Current.FillColor != image.Black {
		t.Error("SetFillColor failed")
	}
}

func TestStackGraphicContext_SetLineWidth(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetLineWidth(5.0)
	if gc.Current.LineWidth != 5.0 {
		t.Errorf("SetLineWidth = %f, want 5.0", gc.Current.LineWidth)
	}
}

func TestStackGraphicContext_SetLineCap(t *testing.T) {
	gc := NewStackGraphicContext()
	caps := []draw2d.LineCap{draw2d.RoundCap, draw2d.ButtCap, draw2d.SquareCap}
	for _, cap := range caps {
		gc.SetLineCap(cap)
		if gc.Current.Cap != cap {
			t.Errorf("SetLineCap(%v) failed", cap)
		}
	}
}

func TestStackGraphicContext_SetLineJoin(t *testing.T) {
	gc := NewStackGraphicContext()
	joins := []draw2d.LineJoin{draw2d.RoundJoin, draw2d.BevelJoin, draw2d.MiterJoin}
	for _, join := range joins {
		gc.SetLineJoin(join)
		if gc.Current.Join != join {
			t.Errorf("SetLineJoin(%v) failed", join)
		}
	}
}

func TestStackGraphicContext_SetLineDash(t *testing.T) {
	gc := NewStackGraphicContext()
	dash := []float64{5, 5}
	offset := 2.5
	gc.SetLineDash(dash, offset)
	if len(gc.Current.Dash) != 2 || gc.Current.Dash[0] != 5 || gc.Current.Dash[1] != 5 {
		t.Error("SetLineDash: dash array not set correctly")
	}
	if gc.Current.DashOffset != offset {
		t.Errorf("SetLineDash: offset = %f, want %f", gc.Current.DashOffset, offset)
	}
}

func TestStackGraphicContext_SetFillRule(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetFillRule(draw2d.FillRuleWinding)
	if gc.Current.FillRule != draw2d.FillRuleWinding {
		t.Error("SetFillRule failed")
	}
	gc.SetFillRule(draw2d.FillRuleEvenOdd)
	if gc.Current.FillRule != draw2d.FillRuleEvenOdd {
		t.Error("SetFillRule failed")
	}
}

func TestStackGraphicContext_FontSize(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetFontSize(12.0)
	if gc.GetFontSize() != 12.0 {
		t.Errorf("FontSize = %f, want 12.0", gc.GetFontSize())
	}
}

func TestStackGraphicContext_FontData(t *testing.T) {
	gc := NewStackGraphicContext()
	fontData := draw2d.FontData{Name: "test", Family: draw2d.FontFamilySerif, Style: draw2d.FontStyleBold}
	gc.SetFontData(fontData)
	result := gc.GetFontData()
	if result.Name != "test" || result.Family != draw2d.FontFamilySerif || result.Style != draw2d.FontStyleBold {
		t.Error("SetFontData/GetFontData failed")
	}
}

func TestStackGraphicContext_GetFontName(t *testing.T) {
	gc := NewStackGraphicContext()
	name := gc.GetFontName()
	if name == "" {
		t.Error("GetFontName should return non-empty string")
	}
}

func TestStackGraphicContext_Translate(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.Translate(5, 10)
	m := gc.GetMatrixTransform()
	x, y := m.GetTranslation()
	if x != 5 || y != 10 {
		t.Errorf("Translate: translation = (%f, %f), want (5, 10)", x, y)
	}
}

func TestStackGraphicContext_Scale(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.Scale(2, 3)
	m := gc.GetMatrixTransform()
	sx, sy := m.GetScaling()
	if sx != 2 || sy != 3 {
		t.Errorf("Scale: scaling = (%f, %f), want (2, 3)", sx, sy)
	}
}

func TestStackGraphicContext_Rotate(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.Rotate(1.57) // ~π/2
	m := gc.GetMatrixTransform()
	if m.IsIdentity() {
		t.Error("After Rotate, matrix should not be identity")
	}
}

func TestStackGraphicContext_SetGetMatrixTransform(t *testing.T) {
	gc := NewStackGraphicContext()
	tr := draw2d.NewTranslationMatrix(10, 20)
	gc.SetMatrixTransform(tr)
	result := gc.GetMatrixTransform()
	if !result.Equals(tr) {
		t.Error("SetMatrixTransform/GetMatrixTransform round trip failed")
	}
}

func TestStackGraphicContext_ComposeMatrixTransform(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.Translate(5, 10)
	scale := draw2d.NewScaleMatrix(2, 3)
	gc.ComposeMatrixTransform(scale)
	m := gc.GetMatrixTransform()
	// Should be composed transformation
	if m.IsIdentity() {
		t.Error("After ComposeMatrixTransform, matrix should not be identity")
	}
}

func TestStackGraphicContext_SaveRestore(t *testing.T) {
	gc := NewStackGraphicContext()
	// Set some values
	gc.SetLineWidth(5.0)
	gc.SetStrokeColor(image.White)
	gc.SetFillColor(image.Black)
	gc.SetFontSize(20.0)
	
	// Save
	gc.Save()
	
	// Change values
	gc.SetLineWidth(10.0)
	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	gc.SetFontSize(30.0)
	
	// Restore
	gc.Restore()
	
	// Check restored values
	if gc.Current.LineWidth != 5.0 {
		t.Errorf("After Restore, LineWidth = %f, want 5.0", gc.Current.LineWidth)
	}
	if gc.Current.StrokeColor != image.White {
		t.Error("After Restore, StrokeColor should be White")
	}
	if gc.Current.FillColor != image.Black {
		t.Error("After Restore, FillColor should be Black")
	}
	if gc.Current.FontSize != 20.0 {
		t.Errorf("After Restore, FontSize = %f, want 20.0", gc.Current.FontSize)
	}
}

func TestStackGraphicContext_SaveRestore_MatrixIndependence(t *testing.T) {
	gc := NewStackGraphicContext()
	// Save
	gc.Save()
	// Translate
	gc.Translate(10, 20)
	// Restore
	gc.Restore()
	// Should be back to identity
	m := gc.GetMatrixTransform()
	if !m.IsIdentity() {
		t.Error("After Save/Translate/Restore, matrix should be identity")
	}
}

func TestStackGraphicContext_RestoreWithoutSave(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetLineWidth(5.0)
	// Restore without Save should not crash
	gc.Restore()
	// Values should be unchanged
	if gc.Current.LineWidth != 5.0 {
		t.Error("Restore without Save should not change values")
	}
}

func TestStackGraphicContext_MultipleSaveRestore(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.SetLineWidth(1.0)
	
	gc.Save()
	gc.SetLineWidth(2.0)
	
	gc.Save()
	gc.SetLineWidth(3.0)
	
	gc.Save()
	gc.SetLineWidth(4.0)
	
	gc.Restore()
	if gc.Current.LineWidth != 3.0 {
		t.Errorf("After 1st Restore, LineWidth = %f, want 3.0", gc.Current.LineWidth)
	}
	
	gc.Restore()
	if gc.Current.LineWidth != 2.0 {
		t.Errorf("After 2nd Restore, LineWidth = %f, want 2.0", gc.Current.LineWidth)
	}
	
	gc.Restore()
	if gc.Current.LineWidth != 1.0 {
		t.Errorf("After 3rd Restore, LineWidth = %f, want 1.0", gc.Current.LineWidth)
	}
}

func TestStackGraphicContext_BeginPath(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(10, 20)
	gc.BeginPath()
	if !gc.IsEmpty() {
		t.Error("BeginPath should clear the path")
	}
}

func TestStackGraphicContext_PathOperations(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(10, 20)
	x, y := gc.LastPoint()
	if x != 10 || y != 20 {
		t.Errorf("LastPoint = (%f, %f), want (10, 20)", x, y)
	}
	gc.LineTo(30, 40)
	x, y = gc.LastPoint()
	if x != 30 || y != 40 {
		t.Errorf("LastPoint after LineTo = (%f, %f), want (30, 40)", x, y)
	}
}

func TestStackGraphicContext_GetPath_ReturnsCopy(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(10, 20)
	p := gc.GetPath()
	// Modify gc's path
	gc.LineTo(30, 40)
	// Returned copy should be unchanged
	if len(p.Components) != 1 {
		t.Error("GetPath should return a copy, not a reference")
	}
}

func TestStackGraphicContext_QuadCurveTo(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(0, 0)
	gc.QuadCurveTo(10, 10, 20, 0)
	x, y := gc.LastPoint()
	if x != 20 || y != 0 {
		t.Errorf("LastPoint after QuadCurveTo = (%f, %f), want (20, 0)", x, y)
	}
}

func TestStackGraphicContext_CubicCurveTo(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(0, 0)
	gc.CubicCurveTo(10, 0, 20, 20, 30, 20)
	x, y := gc.LastPoint()
	if x != 30 || y != 20 {
		t.Errorf("LastPoint after CubicCurveTo = (%f, %f), want (30, 20)", x, y)
	}
}

func TestStackGraphicContext_ArcTo(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(0, 0)
	gc.ArcTo(100, 100, 50, 50, 0, 3.14)
	// Should not crash and should update last point
	x, y := gc.LastPoint()
	if x == 0 && y == 0 {
		t.Error("ArcTo should update last point")
	}
}

func TestStackGraphicContext_Close(t *testing.T) {
	gc := NewStackGraphicContext()
	gc.MoveTo(0, 0)
	gc.LineTo(10, 0)
	gc.Close()
	p := gc.GetPath()
	if len(p.Components) == 0 || p.Components[len(p.Components)-1] != draw2d.CloseCmp {
		t.Error("Close should add CloseCmp")
	}
}
