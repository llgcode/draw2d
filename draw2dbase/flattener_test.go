// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"testing"

	"github.com/llgcode/draw2d"
)

func TestFlatten_EmptyPath(t *testing.T) {
	p := new(draw2d.Path)
	segPath := &SegmentedPath{}
	Flatten(p, segPath, 1.0)
	if len(segPath.Points) != 0 {
		t.Error("Empty path should produce no points")
	}
}

func TestFlatten_MoveTo(t *testing.T) {
	p := new(draw2d.Path)
	p.MoveTo(10, 20)
	segPath := &SegmentedPath{}
	Flatten(p, segPath, 1.0)
	if len(segPath.Points) < 2 {
		t.Error("MoveTo should add points to segmented path")
	}
	if segPath.Points[0] != 10 || segPath.Points[1] != 20 {
		t.Errorf("MoveTo point = (%f, %f), want (10, 20)", segPath.Points[0], segPath.Points[1])
	}
}

func TestFlatten_LineSegments(t *testing.T) {
	p := new(draw2d.Path)
	p.MoveTo(0, 0)
	p.LineTo(10, 10)
	segPath := &SegmentedPath{}
	Flatten(p, segPath, 1.0)
	// Should have at least 4 points (MoveTo + LineTo)
	if len(segPath.Points) < 4 {
		t.Errorf("MoveTo + LineTo should have at least 4 points, got %d", len(segPath.Points))
	}
}

func TestFlatten_WithClose(t *testing.T) {
	p := new(draw2d.Path)
	p.MoveTo(0, 0)
	p.LineTo(10, 0)
	p.LineTo(10, 10)
	p.Close()
	segPath := &SegmentedPath{}
	Flatten(p, segPath, 1.0)
	// Close should add a line back to start
	lastIdx := len(segPath.Points) - 2
	if lastIdx >= 0 {
		lastX, lastY := segPath.Points[lastIdx], segPath.Points[lastIdx+1]
		// Should be back at start (0, 0)
		if lastX != 0 || lastY != 0 {
			t.Errorf("After Close, last point should be (0, 0), got (%f, %f)", lastX, lastY)
		}
	}
}

func TestTransformer_Identity(t *testing.T) {
	segPath := &SegmentedPath{}
	tr := Transformer{
		Tr:        draw2d.NewIdentityMatrix(),
		Flattener: segPath,
	}
	tr.MoveTo(10, 20)
	tr.LineTo(30, 40)
	// Identity transform should pass through
	if segPath.Points[0] != 10 || segPath.Points[1] != 20 {
		t.Error("Identity transform should pass through points")
	}
	if segPath.Points[2] != 30 || segPath.Points[3] != 40 {
		t.Error("Identity transform should pass through points")
	}
}

func TestTransformer_Translation(t *testing.T) {
	segPath := &SegmentedPath{}
	tr := Transformer{
		Tr:        draw2d.NewTranslationMatrix(5, 10),
		Flattener: segPath,
	}
	tr.MoveTo(10, 20)
	// Should be translated to (15, 30)
	if segPath.Points[0] != 15 || segPath.Points[1] != 30 {
		t.Errorf("Translation transform: point = (%f, %f), want (15, 30)", segPath.Points[0], segPath.Points[1])
	}
}

func TestSegmentedPath_MoveTo(t *testing.T) {
	segPath := &SegmentedPath{}
	segPath.MoveTo(10, 20)
	if len(segPath.Points) != 2 {
		t.Error("MoveTo should append 2 points")
	}
	if segPath.Points[0] != 10 || segPath.Points[1] != 20 {
		t.Error("MoveTo should append correct coordinates")
	}
}

func TestSegmentedPath_LineTo(t *testing.T) {
	segPath := &SegmentedPath{}
	segPath.MoveTo(0, 0)
	segPath.LineTo(10, 10)
	if len(segPath.Points) != 4 {
		t.Error("MoveTo + LineTo should have 4 points")
	}
	if segPath.Points[2] != 10 || segPath.Points[3] != 10 {
		t.Error("LineTo should append correct coordinates")
	}
}

func TestDemuxFlattener(t *testing.T) {
	segPath1 := &SegmentedPath{}
	segPath2 := &SegmentedPath{}
	demux := DemuxFlattener{
		Flatteners: []Flattener{segPath1, segPath2},
	}
	demux.MoveTo(10, 20)
	demux.LineTo(30, 40)
	
	// Both flatteners should receive the calls
	if len(segPath1.Points) != 4 || len(segPath2.Points) != 4 {
		t.Error("DemuxFlattener should dispatch to all flatteners")
	}
	if segPath1.Points[0] != 10 || segPath2.Points[0] != 10 {
		t.Error("DemuxFlattener should dispatch correct values")
	}
}
