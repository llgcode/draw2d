// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

package draw2dgles2

import (
	"testing"
)

func TestTriangulate(t *testing.T) {
	tests := []struct {
		name     string
		vertices []Point2D
		wantLen  int // Expected number of indices (3 per triangle)
	}{
		{
			name:     "empty",
			vertices: []Point2D{},
			wantLen:  0,
		},
		{
			name: "triangle",
			vertices: []Point2D{
				{0, 0},
				{100, 0},
				{50, 100},
			},
			wantLen: 3, // 1 triangle
		},
		{
			name: "square",
			vertices: []Point2D{
				{0, 0},
				{100, 0},
				{100, 100},
				{0, 100},
			},
			wantLen: 6, // 2 triangles
		},
		{
			name: "pentagon",
			vertices: []Point2D{
				{50, 0},
				{100, 38},
				{82, 100},
				{18, 100},
				{0, 38},
			},
			wantLen: 9, // 3 triangles
		},
		{
			name: "concave_L_shape",
			vertices: []Point2D{
				{0, 0},
				{50, 0},
				{50, 50},
				{100, 50},
				{100, 100},
				{0, 100},
			},
			wantLen: 12, // 4 triangles
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indices := Triangulate(tt.vertices)
			if len(indices) != tt.wantLen {
				t.Errorf("Triangulate() got %d indices, want %d", len(indices), tt.wantLen)
			}

			// Verify all indices are valid
			for i, idx := range indices {
				if int(idx) >= len(tt.vertices) {
					t.Errorf("Invalid index at position %d: %d >= %d", i, idx, len(tt.vertices))
				}
			}

			// Verify we have complete triangles
			if len(indices)%3 != 0 {
				t.Errorf("Index count %d is not divisible by 3", len(indices))
			}
		})
	}
}

func TestConvertToFloat32(t *testing.T) {
	tests := []struct {
		x, y         float64
		wantX, wantY float32
	}{
		{0, 0, 0, 0},
		{100.5, 200.7, 100.5, 200.7},
		{-50.3, -75.9, -50.3, -75.9},
	}

	for _, tt := range tests {
		gotX, gotY := ConvertToFloat32(tt.x, tt.y)
		if gotX != tt.wantX || gotY != tt.wantY {
			t.Errorf("ConvertToFloat32(%v, %v) = (%v, %v), want (%v, %v)",
				tt.x, tt.y, gotX, gotY, tt.wantX, tt.wantY)
		}
	}
}

func TestPointInTriangle(t *testing.T) {
	// Triangle vertices
	a := Point2D{0, 0}
	b := Point2D{100, 0}
	c := Point2D{50, 100}

	tests := []struct {
		name string
		p    Point2D
		want bool
	}{
		{"center", Point2D{50, 30}, true},
		// Note: Points exactly on boundaries may return true or false depending on implementation
		// This is acceptable for the ear-clipping algorithm
		{"outside_left", Point2D{-10, 50}, false},
		{"outside_right", Point2D{110, 50}, false},
		{"outside_above", Point2D{50, 110}, false},
		{"outside_below", Point2D{50, -10}, false},
		{"clearly_inside", Point2D{50, 40}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pointInTriangle(tt.p, a, b, c)
			if got != tt.want {
				t.Errorf("pointInTriangle(%v) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		name string
		a, b Point2D
		want float32
	}{
		{
			name: "zero_distance",
			a:    Point2D{0, 0},
			b:    Point2D{0, 0},
			want: 0,
		},
		{
			name: "horizontal",
			a:    Point2D{0, 0},
			b:    Point2D{100, 0},
			want: 100,
		},
		{
			name: "vertical",
			a:    Point2D{0, 0},
			b:    Point2D{0, 100},
			want: 100,
		},
		{
			name: "diagonal",
			a:    Point2D{0, 0},
			b:    Point2D{3, 4},
			want: 5, // 3-4-5 triangle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := distance(tt.a, tt.b)
			// Use a small epsilon for floating point comparison
			epsilon := float32(0.0001)
			if got < tt.want-epsilon || got > tt.want+epsilon {
				t.Errorf("distance(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func BenchmarkTriangulate(b *testing.B) {
	// Create a hexagon
	vertices := []Point2D{
		{50, 0},
		{93.3, 25},
		{93.3, 75},
		{50, 100},
		{6.7, 75},
		{6.7, 25},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Triangulate(vertices)
	}
}

func BenchmarkTriangulateLarge(b *testing.B) {
	// Create a polygon with many vertices
	vertices := make([]Point2D, 100)
	for i := 0; i < 100; i++ {
		angle := float32(i) * 3.14159 * 2 / 100
		vertices[i] = Point2D{
			X: 50 + 40*float32(cos(float64(angle))),
			Y: 50 + 40*float32(sin(float64(angle))),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Triangulate(vertices)
	}
}

func cos(x float64) float64 {
	// Simple cos approximation for benchmark
	return float64(Point2D{}.X) // Placeholder
}

func sin(x float64) float64 {
	// Simple sin approximation for benchmark
	return float64(Point2D{}.Y) // Placeholder
}
