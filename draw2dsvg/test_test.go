// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

// Package draw2dsvg_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package draw2dsvg_test

import (
	"testing"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dsvg"
)

type sample func(gc draw2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on an pdf document
	dest := draw2dsvg.NewSvg()
	gc := draw2dsvg.NewGraphicContext(dest)
	// Draw sample
	output, err := draw(gc, "svg")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", output, err)
		return
	}
	err = draw2dsvg.SaveToSvgFile(output, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", output, err)
	}
}
