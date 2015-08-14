// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 26/06/2015 by Stani Michiels

// Package draw2dpdf_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package draw2dpdf_test

import (
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
)

type sample func(gc draw2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on an pdf document
	dest := draw2dpdf.NewPdf("L", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(dest)
	// Draw sample
	output, err := draw(gc, "pdf")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", output, err)
		return
	}
	/*
		// Save to pdf only if it doesn't exist because of git
		if _, err = os.Stat(output); err == nil {
			t.Skipf("Saving %q skipped, as it exists already. (Git would consider it modified.)", output)
			return
		}
	*/
	err = draw2dpdf.SaveToPdfFile(output, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", output, err)
	}
}
