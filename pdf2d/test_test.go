// Package pdf2d_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package pdf2d_test

import (
	"testing"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"
)

func test(t *testing.T, sample draw2d.Sample) {
	// Initialize the graphic context on an pdf document
	dest := pdf2d.NewPdf("L", "mm", "A4")
	gc := pdf2d.NewGraphicContext(dest)
	// Draw Android logo
	fn, err := sample(gc, "pdf")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", fn, err)
		return
	}
	// Save to png
	err = pdf2d.SaveToPdfFile(fn, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", fn, err)
	}
}
