// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2d

import (
	"math"
	"testing"
)

func TestNewIdentityMatrix(t *testing.T) {
	m := NewIdentityMatrix()
	if !m.IsIdentity() {
		t.Error("NewIdentityMatrix should create an identity matrix")
	}
	expected := Matrix{1, 0, 0, 1, 0, 0}
	if !m.Equals(expected) {
		t.Errorf("NewIdentityMatrix = %v, want %v", m, expected)
	}
}

func TestNewTranslationMatrix(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	x, y := m.GetTranslation()
	if !fequals(x, 5) || !fequals(y, 10) {
		t.Errorf("GetTranslation() = (%f, %f), want (5, 10)", x, y)
	}
	if !m.IsTranslation() {
		t.Error("Translation matrix should return true for IsTranslation()")
	}
}

func TestNewScaleMatrix(t *testing.T) {
	m := NewScaleMatrix(2, 3)
	sx, sy := m.GetScaling()
	if !fequals(sx, 2) || !fequals(sy, 3) {
		t.Errorf("GetScaling() = (%f, %f), want (2, 3)", sx, sy)
	}
}

func TestNewRotationMatrix(t *testing.T) {
	m := NewRotationMatrix(math.Pi / 2)
	// Rotating (1,0) by π/2 should give (0,1)
	x, y := m.TransformPoint(1, 0)
	if !fequals(x, 0) || !fequals(y, 1) {
		t.Errorf("Rotating (1,0) by π/2 = (%f, %f), want (0, 1)", x, y)
	}
}

func TestMatrixDeterminant(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want float64
	}{
		{"identity", NewIdentityMatrix(), 1},
		{"scale(2,3)", NewScaleMatrix(2, 3), 6},
		{"translation", NewTranslationMatrix(5, 10), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Determinant()
			if !fequals(got, tt.want) {
				t.Errorf("Determinant() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestMatrixTransformPoint(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	x, y := m.TransformPoint(1, 2)
	if !fequals(x, 6) || !fequals(y, 12) {
		t.Errorf("TransformPoint(1, 2) = (%f, %f), want (6, 12)", x, y)
	}
}

func TestMatrixTransform(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	points := []float64{1, 2, 3, 4}
	m.Transform(points)
	expected := []float64{6, 12, 8, 14}
	for i := range points {
		if !fequals(points[i], expected[i]) {
			t.Errorf("Transform() point[%d] = %f, want %f", i, points[i], expected[i])
		}
	}
}

func TestMatrixInverseTransformPoint(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	x, y := m.InverseTransformPoint(6, 12)
	if !fequals(x, 1) || !fequals(y, 2) {
		t.Errorf("InverseTransformPoint(6, 12) = (%f, %f), want (1, 2)", x, y)
	}
}

func TestMatrixInverseTransform(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	points := []float64{6, 12, 8, 14}
	m.InverseTransform(points)
	expected := []float64{1, 2, 3, 4}
	for i := range points {
		if !fequals(points[i], expected[i]) {
			t.Errorf("InverseTransform() point[%d] = %f, want %f", i, points[i], expected[i])
		}
	}
}

func TestMatrixInverse(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	m.Inverse()
	x, y := m.GetTranslation()
	if !fequals(x, -5) || !fequals(y, -10) {
		t.Errorf("Inverse translation = (%f, %f), want (-5, -10)", x, y)
	}
}

func TestMatrixInverse_RoundTrip(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	origX, origY := 3.0, 7.0
	// Transform
	x, y := m.TransformPoint(origX, origY)
	// Inverse transform
	x2, y2 := m.InverseTransformPoint(x, y)
	if !fequals(x2, origX) || !fequals(y2, origY) {
		t.Errorf("Round trip: got (%f, %f), want (%f, %f)", x2, y2, origX, origY)
	}
}

func TestMatrixCompose(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	scale := NewScaleMatrix(2, 3)
	m.Compose(scale)
	// Composed matrix should scale then translate
	x, y := m.TransformPoint(1, 1)
	expected_x := 1.0*2 + 5
	expected_y := 1.0*3 + 10
	if !fequals(x, expected_x) || !fequals(y, expected_y) {
		t.Errorf("Compose transform = (%f, %f), want (%f, %f)", x, y, expected_x, expected_y)
	}
}

func TestMatrixCopy(t *testing.T) {
	m1 := NewTranslationMatrix(5, 10)
	m2 := m1.Copy()
	// Modify m2
	m2.Translate(1, 1)
	// m1 should be unchanged
	x, y := m1.GetTranslation()
	if !fequals(x, 5) || !fequals(y, 10) {
		t.Error("Copy should be independent of original")
	}
}

func TestMatrixTransformRectangle(t *testing.T) {
	m := NewIdentityMatrix()
	x0, y0, x2, y2 := m.TransformRectangle(10, 10, 20, 20)
	if !fequals(x0, 10) || !fequals(y0, 10) || !fequals(x2, 20) || !fequals(y2, 20) {
		t.Errorf("Identity transform rectangle = (%f, %f, %f, %f), want (10, 10, 20, 20)", x0, y0, x2, y2)
	}
}

func TestMatrixTransformRectangle_WithScale(t *testing.T) {
	m := NewScaleMatrix(2, 3)
	x0, y0, x2, y2 := m.TransformRectangle(10, 10, 20, 20)
	if !fequals(x0, 20) || !fequals(y0, 30) || !fequals(x2, 40) || !fequals(y2, 60) {
		t.Errorf("Scale transform rectangle = (%f, %f, %f, %f), want (20, 30, 40, 60)", x0, y0, x2, y2)
	}
}

func TestMatrixVectorTransform(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	points := []float64{1, 2}
	m.VectorTransform(points)
	// Translation should be ignored in VectorTransform
	if !fequals(points[0], 1) || !fequals(points[1], 2) {
		t.Errorf("VectorTransform should ignore translation: got (%f, %f), want (1, 2)", points[0], points[1])
	}
}

func TestMatrixEquals(t *testing.T) {
	m1 := NewIdentityMatrix()
	m2 := NewIdentityMatrix()
	if !m1.Equals(m2) {
		t.Error("Two identity matrices should be equal")
	}
	m3 := NewTranslationMatrix(1, 1)
	if m1.Equals(m3) {
		t.Error("Identity and translation matrices should not be equal")
	}
}

func TestMatrixIsIdentity(t *testing.T) {
	m := NewIdentityMatrix()
	if !m.IsIdentity() {
		t.Error("Identity matrix should return true for IsIdentity()")
	}
	m.Translate(1, 1)
	if m.IsIdentity() {
		t.Error("Translated matrix should not be identity")
	}
}

func TestMatrixIsTranslation(t *testing.T) {
	m := NewTranslationMatrix(5, 10)
	if !m.IsTranslation() {
		t.Error("Translation matrix should return true for IsTranslation()")
	}
	m2 := NewScaleMatrix(2, 2)
	if m2.IsTranslation() {
		t.Error("Scale matrix should not be a translation")
	}
}

func TestMatrixScale(t *testing.T) {
	m := NewIdentityMatrix()
	m.Scale(2, 3)
	sx, sy := m.GetScaling()
	if !fequals(sx, 2) || !fequals(sy, 3) {
		t.Errorf("Scale() result GetScaling() = (%f, %f), want (2, 3)", sx, sy)
	}
}

func TestMatrixTranslate(t *testing.T) {
	m := NewIdentityMatrix()
	m.Translate(5, 10)
	x, y := m.GetTranslation()
	if !fequals(x, 5) || !fequals(y, 10) {
		t.Errorf("Translate() result GetTranslation() = (%f, %f), want (5, 10)", x, y)
	}
}

func TestMatrixRotate(t *testing.T) {
	m := NewIdentityMatrix()
	m.Rotate(math.Pi)
	// Rotating (1,0) by π should give (-1,0)
	x, y := m.TransformPoint(1, 0)
	if !fequals(x, -1) || !fequals(y, 0) {
		t.Errorf("Rotate(π) on (1,0) = (%f, %f), want (-1, 0)", x, y)
	}
}

func TestNewMatrixFromRects(t *testing.T) {
	rect1 := [4]float64{0, 0, 10, 10}
	rect2 := [4]float64{0, 0, 20, 20}
	m := NewMatrixFromRects(rect1, rect2)
	// Midpoint (5,5) of rect1 should map to midpoint (10,10) of rect2
	x, y := m.TransformPoint(5, 5)
	if !fequals(x, 10) || !fequals(y, 10) {
		t.Errorf("NewMatrixFromRects transform midpoint = (%f, %f), want (10, 10)", x, y)
	}
}

func TestMatrixGetScale(t *testing.T) {
	m := NewIdentityMatrix()
	scale := m.GetScale()
	if !fequals(scale, 1.0) {
		t.Errorf("GetScale() on identity = %f, want ~1.0", scale)
	}
}

func TestMatrixGetTranslation(t *testing.T) {
	m := NewTranslationMatrix(3, 7)
	x, y := m.GetTranslation()
	if !fequals(x, 3) || !fequals(y, 7) {
		t.Errorf("GetTranslation() = (%f, %f), want (3, 7)", x, y)
	}
}

func TestMatrixGetScaling(t *testing.T) {
	m := NewScaleMatrix(4, 5)
	sx, sy := m.GetScaling()
	if !fequals(sx, 4) || !fequals(sy, 5) {
		t.Errorf("GetScaling() = (%f, %f), want (4, 5)", sx, sy)
	}
}
