// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2d

import (
	"image/color"
	"testing"
)

func TestLineCap_String(t *testing.T) {
	tests := []struct {
		name     string
		cap      LineCap
		expected string
	}{
		{"RoundCap", RoundCap, "round"},
		{"ButtCap", ButtCap, "cap"},
		{"SquareCap", SquareCap, "square"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cap.String(); got != tt.expected {
				t.Errorf("LineCap.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLineJoin_String(t *testing.T) {
	tests := []struct {
		name     string
		join     LineJoin
		expected string
	}{
		{"BevelJoin", BevelJoin, "bevel"},
		{"RoundJoin", RoundJoin, "round"},
		{"MiterJoin", MiterJoin, "miter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.join.String(); got != tt.expected {
				t.Errorf("LineJoin.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFillRule_Constants(t *testing.T) {
	tests := []struct {
		name     string
		rule     FillRule
		expected int
	}{
		{"FillRuleEvenOdd", FillRuleEvenOdd, 0},
		{"FillRuleWinding", FillRuleWinding, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.rule) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, int(tt.rule), tt.expected)
			}
		})
	}
}

func TestValign_Constants(t *testing.T) {
	tests := []struct {
		name     string
		valign   Valign
		expected int
	}{
		{"ValignTop", ValignTop, 0},
		{"ValignCenter", ValignCenter, 1},
		{"ValignBottom", ValignBottom, 2},
		{"ValignBaseline", ValignBaseline, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.valign) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, int(tt.valign), tt.expected)
			}
		})
	}
}

func TestHalign_Constants(t *testing.T) {
	tests := []struct {
		name     string
		halign   Halign
		expected int
	}{
		{"HalignLeft", HalignLeft, 0},
		{"HalignCenter", HalignCenter, 1},
		{"HalignRight", HalignRight, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.halign) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, int(tt.halign), tt.expected)
			}
		})
	}
}

func TestStrokeStyle_Defaults(t *testing.T) {
	// Create a StrokeStyle and verify fields can be set and read
	style := StrokeStyle{
		Color:      color.RGBA{255, 0, 0, 255},
		Width:      5.0,
		LineCap:    RoundCap,
		LineJoin:   MiterJoin,
		DashOffset: 2.5,
		Dash:       []float64{10, 5},
	}

	if style.Width != 5.0 {
		t.Errorf("StrokeStyle.Width = %v, want %v", style.Width, 5.0)
	}

	if style.LineCap != RoundCap {
		t.Errorf("StrokeStyle.LineCap = %v, want %v", style.LineCap, RoundCap)
	}

	if style.LineJoin != MiterJoin {
		t.Errorf("StrokeStyle.LineJoin = %v, want %v", style.LineJoin, MiterJoin)
	}

	if style.DashOffset != 2.5 {
		t.Errorf("StrokeStyle.DashOffset = %v, want %v", style.DashOffset, 2.5)
	}

	if len(style.Dash) != 2 || style.Dash[0] != 10 || style.Dash[1] != 5 {
		t.Errorf("StrokeStyle.Dash = %v, want %v", style.Dash, []float64{10, 5})
	}

	r, g, b, a := style.Color.RGBA()
	if r != 65535 || g != 0 || b != 0 || a != 65535 {
		t.Errorf("StrokeStyle.Color RGBA values incorrect")
	}
}

func TestSolidFillStyle_Defaults(t *testing.T) {
	// Create a SolidFillStyle and verify fields can be set and read
	style := SolidFillStyle{
		Color:    color.RGBA{0, 255, 0, 255},
		FillRule: FillRuleWinding,
	}

	if style.FillRule != FillRuleWinding {
		t.Errorf("SolidFillStyle.FillRule = %v, want %v", style.FillRule, FillRuleWinding)
	}

	r, g, b, a := style.Color.RGBA()
	if r != 0 || g != 65535 || b != 0 || a != 65535 {
		t.Errorf("SolidFillStyle.Color RGBA values incorrect")
	}
}
