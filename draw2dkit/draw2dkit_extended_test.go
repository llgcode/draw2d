// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dkit

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

func newTestGC(t *testing.T) *draw2dimg.GraphicContext {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	return draw2dimg.NewGraphicContext(img)
}

func TestRectangle(t *testing.T) {
	p := new(draw2d.Path)
	Rectangle(p, 10, 10, 100, 100)
	if len(p.Components) != 5 {
		t.Errorf("Rectangle should have 5 components, got %d", len(p.Components))
	}
	// Should be: MoveTo + 3 LineTo + Close
	if p.Components[0] != draw2d.MoveToCmp {
		t.Error("First component should be MoveTo")
	}
	if p.Components[len(p.Components)-1] != draw2d.CloseCmp {
		t.Error("Last component should be Close")
	}
}

func TestRoundedRectangle(t *testing.T) {
	p := new(draw2d.Path)
	RoundedRectangle(p, 10, 10, 100, 100, 10, 10)
	if p.IsEmpty() {
		t.Error("RoundedRectangle should not be empty")
	}
	// Should contain QuadCurveToCmp
	hasQuadCurve := false
	for _, cmp := range p.Components {
		if cmp == draw2d.QuadCurveToCmp {
			hasQuadCurve = true
			break
		}
	}
	if !hasQuadCurve {
		t.Error("RoundedRectangle should contain QuadCurveToCmp")
	}
}

func TestEllipse(t *testing.T) {
	p := new(draw2d.Path)
	Ellipse(p, 100, 100, 50, 30)
	if p.IsEmpty() {
		t.Error("Ellipse should not be empty")
	}
	// Should contain ArcToCmp
	hasArc := false
	for _, cmp := range p.Components {
		if cmp == draw2d.ArcToCmp {
			hasArc = true
			break
		}
	}
	if !hasArc {
		t.Error("Ellipse should contain ArcToCmp")
	}
}

func TestCircle_PathComponents(t *testing.T) {
	p := new(draw2d.Path)
	Circle(p, 100, 100, 50)
	// Should contain ArcTo and Close
	hasArc := false
	hasClose := false
	for _, cmp := range p.Components {
		if cmp == draw2d.ArcToCmp {
			hasArc = true
		}
		if cmp == draw2d.CloseCmp {
			hasClose = true
		}
	}
	if !hasArc {
		t.Error("Circle should contain ArcToCmp")
	}
	if !hasClose {
		t.Error("Circle should contain CloseCmp")
	}
}

func TestCircle_StrokeDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Circle Stroke panicked: %v", r)
		}
	}()
	gc := newTestGC(t)
	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(1)
	Circle(gc, 100, 100, 50)
	gc.Stroke()
}

func TestCircle_FillDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Circle Fill panicked: %v", r)
		}
	}()
	gc := newTestGC(t)
	gc.SetFillColor(color.NRGBA{0, 0, 255, 255})
	Circle(gc, 100, 100, 50)
	gc.Fill()
}

func TestRectangle_FillStrokeDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Rectangle FillStroke panicked: %v", r)
		}
	}()
	gc := newTestGC(t)
	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetFillColor(color.NRGBA{0, 255, 0, 255})
	gc.SetLineWidth(2)
	Rectangle(gc, 20, 20, 180, 180)
	gc.FillStroke()
}

func TestEllipse_DifferentRadii(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Ellipse with different radii panicked: %v", r)
		}
	}()
	gc := newTestGC(t)
	gc.SetStrokeColor(color.NRGBA{255, 0, 255, 255})
	gc.SetLineWidth(1)
	Ellipse(gc, 100, 100, 80, 40)
	gc.Stroke()
}

func TestRoundedRectangle_FillStrokeDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RoundedRectangle FillStroke panicked: %v", r)
		}
	}()
	gc := newTestGC(t)
	gc.SetStrokeColor(color.NRGBA{0, 0, 0, 255})
	gc.SetFillColor(color.NRGBA{255, 255, 0, 255})
	gc.SetLineWidth(2)
	RoundedRectangle(gc, 20, 20, 180, 180, 20, 20)
	gc.FillStroke()
}

func TestCircle_FullCircle(t *testing.T) {
	p := new(draw2d.Path)
	Circle(p, 100, 100, 50)
	// Arc should be close to 2π
	foundArc := false
	for i, cmp := range p.Components {
		if cmp == draw2d.ArcToCmp {
			foundArc = true
			// ArcTo has 6 parameters, angle is the last one
			pointIdx := 0
			for j := 0; j < i; j++ {
				switch p.Components[j] {
				case draw2d.MoveToCmp, draw2d.LineToCmp:
					pointIdx += 2
				case draw2d.QuadCurveToCmp:
					pointIdx += 4
				case draw2d.CubicCurveToCmp, draw2d.ArcToCmp:
					pointIdx += 6
				}
			}
			if pointIdx+5 < len(p.Points) {
				angle := p.Points[pointIdx+5]
				// Should be close to -2π
				if math.Abs(angle+2*math.Pi) > 0.01 {
					t.Errorf("Circle angle = %f, want ~-2π", angle)
				}
			}
			break
		}
	}
	if !foundArc {
		t.Error("Circle should contain an arc")
	}
}
