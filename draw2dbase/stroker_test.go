// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"math"
	"testing"

	"github.com/llgcode/draw2d"
)

func TestLineStroker_BasicLine(t *testing.T) {
	segPath := &SegmentedPath{}
	stroker := NewLineStroker(draw2d.RoundCap, draw2d.RoundJoin, segPath)
	stroker.HalfLineWidth = 1.0
	
	stroker.MoveTo(0, 0)
	stroker.LineTo(10, 0)
	stroker.End()
	
	// Should produce output in the inner flattener
	if len(segPath.Points) == 0 {
		t.Error("LineStroker should produce output")
	}
}

func TestLineStroker_HalfLineWidth(t *testing.T) {
	segPath1 := &SegmentedPath{}
	stroker1 := NewLineStroker(draw2d.RoundCap, draw2d.RoundJoin, segPath1)
	stroker1.HalfLineWidth = 1.0
	stroker1.MoveTo(0, 0)
	stroker1.LineTo(10, 0)
	stroker1.End()
	
	segPath2 := &SegmentedPath{}
	stroker2 := NewLineStroker(draw2d.RoundCap, draw2d.RoundJoin, segPath2)
	stroker2.HalfLineWidth = 2.0
	stroker2.MoveTo(0, 0)
	stroker2.LineTo(10, 0)
	stroker2.End()
	
	// Different line widths should produce different output
	if len(segPath1.Points) == len(segPath2.Points) {
		// Check if points are actually different
		allSame := true
		for i := range segPath1.Points {
			if segPath1.Points[i] != segPath2.Points[i] {
				allSame = false
				break
			}
		}
		if allSame && len(segPath1.Points) > 0 {
			t.Error("Different HalfLineWidth should produce different output")
		}
	}
}

func TestLineStroker_End(t *testing.T) {
	segPath := &SegmentedPath{}
	stroker := NewLineStroker(draw2d.RoundCap, draw2d.RoundJoin, segPath)
	stroker.HalfLineWidth = 1.0
	
	stroker.MoveTo(0, 0)
	stroker.LineTo(10, 0)
	stroker.LineTo(10, 10)
	
	initialLen := len(segPath.Points)
	stroker.End()
	afterEndLen := len(segPath.Points)
	
	// End should flush output
	if afterEndLen <= initialLen {
		t.Error("End should flush and add points to output")
	}
	
	// After End, internal state should be reset
	if len(stroker.vertices) != 0 || len(stroker.rewind) != 0 {
		t.Error("End should reset internal vertices and rewind")
	}
}

func TestVectorDistance(t *testing.T) {
	tests := []struct {
		name     string
		dx, dy   float64
		expected float64
	}{
		{"horizontal", 3, 0, 3},
		{"vertical", 0, 4, 4},
		{"diagonal", 3, 4, 5},
		{"zero", 0, 0, 0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vectorDistance(tt.dx, tt.dy)
			if math.Abs(result-tt.expected) > 1e-6 {
				t.Errorf("vectorDistance(%f, %f) = %f, want %f", tt.dx, tt.dy, result, tt.expected)
			}
		})
	}
}
