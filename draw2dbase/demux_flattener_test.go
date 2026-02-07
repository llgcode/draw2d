// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"testing"
)

func TestDemuxFlattener_LineJoin(t *testing.T) {
	mock1 := &mockFlattener{}
	mock2 := &mockFlattener{}

	demux := DemuxFlattener{
		Flatteners: []Flattener{mock1, mock2},
	}

	demux.LineJoin()

	if !mock1.lineJoinCalled {
		t.Error("LineJoin not dispatched to first flattener")
	}

	if !mock2.lineJoinCalled {
		t.Error("LineJoin not dispatched to second flattener")
	}
}

func TestDemuxFlattener_Close(t *testing.T) {
	mock1 := &mockFlattener{}
	mock2 := &mockFlattener{}

	demux := DemuxFlattener{
		Flatteners: []Flattener{mock1, mock2},
	}

	demux.Close()

	if !mock1.closeCalled {
		t.Error("Close not dispatched to first flattener")
	}

	if !mock2.closeCalled {
		t.Error("Close not dispatched to second flattener")
	}
}

func TestDemuxFlattener_End(t *testing.T) {
	mock1 := &mockFlattener{}
	mock2 := &mockFlattener{}

	demux := DemuxFlattener{
		Flatteners: []Flattener{mock1, mock2},
	}

	demux.End()

	if !mock1.endCalled {
		t.Error("End not dispatched to first flattener")
	}

	if !mock2.endCalled {
		t.Error("End not dispatched to second flattener")
	}
}

func TestDemuxFlattener_Empty(t *testing.T) {
	// Create DemuxFlattener with empty slice - should not panic
	demux := DemuxFlattener{
		Flatteners: []Flattener{},
	}

	// These should not panic with empty flatteners
	demux.MoveTo(10, 10)
	demux.LineTo(20, 20)
	demux.LineJoin()
	demux.Close()
	demux.End()
}

// mockFlattener implements Flattener for testing
type mockFlattener struct {
	lineJoinCalled bool
	closeCalled    bool
	endCalled      bool
	moveToX        float64
	moveToY        float64
	lineToX        float64
	lineToY        float64
}

func (m *mockFlattener) MoveTo(x, y float64) {
	m.moveToX = x
	m.moveToY = y
}

func (m *mockFlattener) LineTo(x, y float64) {
	m.lineToX = x
	m.lineToY = y
}

func (m *mockFlattener) LineJoin() {
	m.lineJoinCalled = true
}

func (m *mockFlattener) Close() {
	m.closeCalled = true
}

func (m *mockFlattener) End() {
	m.endCalled = true
}
