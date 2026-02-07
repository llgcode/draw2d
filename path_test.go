// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2d

import (
	"math"
	"strings"
	"testing"
)

func TestPathMoveTo(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	x, y := p.LastPoint()
	if x != 10 || y != 20 {
		t.Errorf("LastPoint() = (%f, %f), want (10, 20)", x, y)
	}
	if len(p.Components) != 1 || p.Components[0] != MoveToCmp {
		t.Error("MoveTo should add a MoveToCmp component")
	}
	if len(p.Points) != 2 || p.Points[0] != 10 || p.Points[1] != 20 {
		t.Error("MoveTo should add correct points")
	}
}

func TestPathLineTo(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	p.LineTo(30, 40)
	if len(p.Components) != 2 || p.Components[1] != LineToCmp {
		t.Error("LineTo should add a LineToCmp component")
	}
	x, y := p.LastPoint()
	if x != 30 || y != 40 {
		t.Errorf("LastPoint after LineTo = (%f, %f), want (30, 40)", x, y)
	}
}

func TestPathLineToWithoutMoveTo(t *testing.T) {
	p := new(Path)
	p.LineTo(10, 20)
	// Should auto-create MoveTo
	if len(p.Components) != 1 || p.Components[0] != MoveToCmp {
		t.Error("LineTo without MoveTo should auto-create MoveTo")
	}
}

func TestPathQuadCurveTo(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.QuadCurveTo(10, 10, 20, 0)
	if len(p.Components) != 2 || p.Components[1] != QuadCurveToCmp {
		t.Error("QuadCurveTo should add QuadCurveToCmp component")
	}
	// MoveTo adds 2 points, QuadCurveTo adds 4 points
	if len(p.Points) != 6 {
		t.Errorf("Points count = %d, want 6", len(p.Points))
	}
}

func TestPathQuadCurveToWithoutMoveTo(t *testing.T) {
	p := new(Path)
	p.QuadCurveTo(10, 10, 20, 0)
	// Should auto-create MoveTo
	if len(p.Components) != 1 || p.Components[0] != MoveToCmp {
		t.Error("QuadCurveTo without MoveTo should auto-create MoveTo")
	}
}

func TestPathCubicCurveTo(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.CubicCurveTo(10, 0, 20, 20, 30, 20)
	if len(p.Components) != 2 || p.Components[1] != CubicCurveToCmp {
		t.Error("CubicCurveTo should add CubicCurveToCmp component")
	}
	// MoveTo adds 2 points, CubicCurveTo adds 6 points
	if len(p.Points) != 8 {
		t.Errorf("Points count = %d, want 8", len(p.Points))
	}
}

func TestPathCubicCurveToWithoutMoveTo(t *testing.T) {
	p := new(Path)
	p.CubicCurveTo(10, 0, 20, 20, 30, 20)
	// Should auto-create MoveTo
	if len(p.Components) != 1 || p.Components[0] != MoveToCmp {
		t.Error("CubicCurveTo without MoveTo should auto-create MoveTo")
	}
}

func TestPathArcTo(t *testing.T) {
	p := new(Path)
	p.MoveTo(100, 100)
	// Full circle
	p.ArcTo(100, 100, 50, 50, 0, 2*math.Pi)
	x, y := p.LastPoint()
	// After full circle, should return near start point
	if math.Abs(x-150) > epsilon || math.Abs(y-100) > epsilon {
		t.Errorf("ArcTo full circle end point = (%f, %f), want (~150, ~100)", x, y)
	}
}

func TestPathArcTo_EmptyPath(t *testing.T) {
	p := new(Path)
	p.ArcTo(100, 100, 50, 50, 0, math.Pi)
	// Should start with MoveTo
	if len(p.Components) == 0 || p.Components[0] != MoveToCmp {
		t.Error("ArcTo on empty path should start with MoveTo")
	}
}

func TestPathArcTo_ExistingPath(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.ArcTo(100, 100, 50, 50, 0, math.Pi)
	// Should have MoveTo, LineTo, ArcTo
	if len(p.Components) < 3 {
		t.Errorf("ArcTo on existing path should add LineTo and ArcTo, got %d components", len(p.Components))
	}
	foundArc := false
	for _, cmp := range p.Components {
		if cmp == ArcToCmp {
			foundArc = true
		}
	}
	if !foundArc {
		t.Error("ArcTo should add ArcToCmp component")
	}
}

func TestPathClose(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.LineTo(10, 0)
	p.Close()
	if len(p.Components) == 0 || p.Components[len(p.Components)-1] != CloseCmp {
		t.Error("Close should add CloseCmp component")
	}
}

func TestPathCopy(t *testing.T) {
	p1 := new(Path)
	p1.MoveTo(10, 20)
	p1.LineTo(30, 40)
	p2 := p1.Copy()
	// Modify p2
	p2.LineTo(50, 60)
	// p1 should be unchanged
	if len(p1.Components) != 2 {
		t.Error("Copy should be independent of original")
	}
}

func TestPathClear(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	p.LineTo(30, 40)
	p.Clear()
	if !p.IsEmpty() {
		t.Error("Clear should make path empty")
	}
	if len(p.Components) != 0 || len(p.Points) != 0 {
		t.Error("Clear should remove all components and points")
	}
}

func TestPathIsEmpty(t *testing.T) {
	p := new(Path)
	if !p.IsEmpty() {
		t.Error("New path should be empty")
	}
	p.MoveTo(10, 20)
	if p.IsEmpty() {
		t.Error("Path with MoveTo should not be empty")
	}
}

func TestPathString(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	p.LineTo(30, 40)
	p.Close()
	s := p.String()
	if !strings.Contains(s, "MoveTo") {
		t.Error("String should contain 'MoveTo'")
	}
	if !strings.Contains(s, "LineTo") {
		t.Error("String should contain 'LineTo'")
	}
	if !strings.Contains(s, "Close") {
		t.Error("String should contain 'Close'")
	}
}

func TestPathString_AllComponents(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.LineTo(10, 10)
	p.QuadCurveTo(20, 20, 30, 10)
	p.CubicCurveTo(40, 0, 50, 0, 60, 10)
	p.ArcTo(70, 10, 5, 5, 0, math.Pi)
	p.Close()
	s := p.String()
	keywords := []string{"MoveTo", "LineTo", "QuadCurveTo", "CubicCurveTo", "ArcTo", "Close"}
	for _, kw := range keywords {
		if !strings.Contains(s, kw) {
			t.Errorf("String should contain '%s'", kw)
		}
	}
}

func TestPathVerticalFlip(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	p.LineTo(30, 40)
	p2 := p.VerticalFlip()
	// Y coordinates should be negated
	if p2.Points[1] != -20 || p2.Points[3] != -40 {
		t.Error("VerticalFlip should negate Y coordinates")
	}
	// Original should be unchanged
	if p.Points[1] != 20 || p.Points[3] != 40 {
		t.Error("VerticalFlip should not modify original")
	}
}

func TestPathVerticalFlip_LastPoint(t *testing.T) {
	p := new(Path)
	p.MoveTo(10, 20)
	p2 := p.VerticalFlip()
	_, y := p2.LastPoint()
	if y != -20 {
		t.Errorf("Flipped LastPoint Y = %f, want -20", y)
	}
}

func TestPathMultipleSubpaths(t *testing.T) {
	p := new(Path)
	p.MoveTo(0, 0)
	p.LineTo(10, 10)
	p.Close()
	p.MoveTo(20, 20)
	p.LineTo(30, 30)
	p.Close()
	// Should have 2 MoveTo + 2 LineTo + 2 Close = 6 components
	if len(p.Components) != 6 {
		t.Errorf("Multiple subpaths: got %d components, want 6", len(p.Components))
	}
}
