// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"math"
	"testing"
)

func TestSubdivideCubic(t *testing.T) {
	c := []float64{0, 0, 10, 20, 30, 40, 50, 50}
	c1 := make([]float64, 8)
	c2 := make([]float64, 8)

	SubdivideCubic(c, c1, c2)

	// Verify that the endpoint of c1 equals the start point of c2
	if c1[6] != c2[0] || c1[7] != c2[1] {
		t.Errorf("SubdivideCubic: c1 endpoint (%v, %v) != c2 start point (%v, %v)",
			c1[6], c1[7], c2[0], c2[1])
	}
}

func TestSubdivideCubic_Endpoints(t *testing.T) {
	c := []float64{10, 20, 30, 40, 50, 60, 70, 80}
	c1 := make([]float64, 8)
	c2 := make([]float64, 8)

	SubdivideCubic(c, c1, c2)

	// Verify first point of c1 equals first point of c
	if c1[0] != c[0] || c1[1] != c[1] {
		t.Errorf("SubdivideCubic: c1 start (%v, %v) != c start (%v, %v)",
			c1[0], c1[1], c[0], c[1])
	}

	// Verify last point of c2 equals last point of c
	if c2[6] != c[6] || c2[7] != c[7] {
		t.Errorf("SubdivideCubic: c2 end (%v, %v) != c end (%v, %v)",
			c2[6], c2[7], c[6], c[7])
	}
}

func TestTraceCubic_ErrorOnShortSlice(t *testing.T) {
	var liner mockLiner
	shortSlice := []float64{0, 0, 10, 10, 20, 20}

	err := TraceCubic(&liner, shortSlice, 0.5)
	if err == nil {
		t.Error("TraceCubic should return error for slice with length < 8")
	}
}

func TestTraceCubic_ValidCurve(t *testing.T) {
	var liner mockLiner
	curve := []float64{0, 0, 10, 20, 30, 40, 50, 50}

	err := TraceCubic(&liner, curve, 0.5)
	if err != nil {
		t.Errorf("TraceCubic returned unexpected error: %v", err)
	}

	if len(liner.points) == 0 {
		t.Error("TraceCubic did not produce any line segments")
	}
}

func TestTraceQuad_ValidCurve(t *testing.T) {
	var liner mockLiner
	curve := []float64{0, 0, 25, 50, 50, 0}

	err := TraceQuad(&liner, curve, 0.5)
	if err != nil {
		t.Errorf("TraceQuad returned unexpected error: %v", err)
	}

	if len(liner.points) == 0 {
		t.Error("TraceQuad did not produce any line segments")
	}
}

func TestTraceArc(t *testing.T) {
	var liner mockLiner

	// Trace a 90-degree arc
	x, y := 100.0, 100.0
	rx, ry := 50.0, 50.0
	start := 0.0
	angle := math.Pi / 2 // 90 degrees
	scale := 1.0

	lastX, lastY := TraceArc(&liner, x, y, rx, ry, start, angle, scale)

	// Verify that TraceArc produces valid endpoint coordinates
	if math.IsNaN(lastX) || math.IsNaN(lastY) {
		t.Error("TraceArc produced NaN coordinates")
	}

	// Verify that some line segments were produced
	if len(liner.points) == 0 {
		t.Error("TraceArc did not produce any line segments")
	}

	// Verify endpoint is approximately where we expect (x, y + ry)
	expectedX := x + math.Cos(start+angle)*rx
	expectedY := y + math.Sin(start+angle)*ry

	tolerance := 0.1
	if math.Abs(lastX-expectedX) > tolerance || math.Abs(lastY-expectedY) > tolerance {
		t.Errorf("TraceArc endpoint (%v, %v) not close to expected (%v, %v)",
			lastX, lastY, expectedX, expectedY)
	}
}

// mockLiner is a simple implementation of Liner for testing
type mockLiner struct {
	points []float64
}

func (m *mockLiner) MoveTo(x, y float64) {
	m.points = append(m.points, x, y)
}

func (m *mockLiner) LineTo(x, y float64) {
	m.points = append(m.points, x, y)
}

func (m *mockLiner) LineJoin() {}

func (m *mockLiner) Close() {}

func (m *mockLiner) End() {}
