// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"testing"
)

// Test related to issue #95: DashVertexConverter state preservation
func TestDashVertexConverter_StatePreservation(t *testing.T) {
	segPath := &SegmentedPath{}
	dash := []float64{5, 5}
	dasher := NewDashConverter(dash, 0, segPath)
	
	dasher.MoveTo(0, 0)
	initialLen := len(segPath.Points)
	
	dasher.LineTo(10, 0)
	afterFirstLen := len(segPath.Points)
	
	dasher.LineTo(20, 0)
	afterSecondLen := len(segPath.Points)
	
	// Second LineTo should add more points
	if afterSecondLen <= afterFirstLen {
		t.Error("Second LineTo should add more points, state may not be preserved")
	}
	if initialLen >= afterFirstLen {
		t.Error("First LineTo should add points")
	}
}

func TestDashVertexConverter_SingleDash(t *testing.T) {
	segPath := &SegmentedPath{}
	dash := []float64{10}
	dasher := NewDashConverter(dash, 0, segPath)
	
	dasher.MoveTo(0, 0)
	dasher.LineTo(50, 0)
	
	// Should produce output
	if len(segPath.Points) == 0 {
		t.Error("Single-element dash array should produce output")
	}
}

func TestDashVertexConverter_DashOffset(t *testing.T) {
	segPath1 := &SegmentedPath{}
	segPath2 := &SegmentedPath{}
	dash := []float64{5, 5}
	
	dasher1 := NewDashConverter(dash, 0, segPath1)
	dasher1.MoveTo(0, 0)
	dasher1.LineTo(50, 0)
	
	dasher2 := NewDashConverter(dash, 2.5, segPath2)
	dasher2.MoveTo(0, 0)
	dasher2.LineTo(50, 0)
	
	// Different offsets should produce different output
	if len(segPath1.Points) == len(segPath2.Points) {
		// Check if points are actually different
		allSame := true
		minLen := len(segPath1.Points)
		if len(segPath2.Points) < minLen {
			minLen = len(segPath2.Points)
		}
		for i := 0; i < minLen; i++ {
			if segPath1.Points[i] != segPath2.Points[i] {
				allSame = false
				break
			}
		}
		if allSame && len(segPath1.Points) > 0 {
			t.Error("Different dash offsets should produce different output")
		}
	}
}

func TestDashVertexConverter_Close(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close panicked: %v", r)
		}
	}()
	segPath := &SegmentedPath{}
	dash := []float64{5, 5}
	dasher := NewDashConverter(dash, 0, segPath)
	dasher.MoveTo(0, 0)
	dasher.LineTo(10, 10)
	dasher.Close()
}

func TestDashVertexConverter_End(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("End panicked: %v", r)
		}
	}()
	segPath := &SegmentedPath{}
	dash := []float64{5, 5}
	dasher := NewDashConverter(dash, 0, segPath)
	dasher.MoveTo(0, 0)
	dasher.LineTo(10, 10)
	dasher.End()
}

func TestDashVertexConverter_MoveTo(t *testing.T) {
	segPath := &SegmentedPath{}
	dash := []float64{5, 5}
	dasher := NewDashConverter(dash, 0, segPath)
	
	dasher.MoveTo(10, 20)
	// Check that position is set correctly
	if dasher.x != 10 || dasher.y != 20 {
		t.Errorf("MoveTo should set position to (10, 20), got (%f, %f)", dasher.x, dasher.y)
	}
	// Check that distance is reset to dashOffset
	if dasher.distance != dasher.dashOffset {
		t.Error("MoveTo should reset distance to dashOffset")
	}
	// Check that currentDash is reset
	if dasher.currentDash != 0 {
		t.Error("MoveTo should reset currentDash to 0")
	}
}
