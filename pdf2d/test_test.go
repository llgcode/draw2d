// Package pdf2d_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package pdf2d_test

import (
	"os"
	"testing"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"
)

func test(t *testing.T, sample draw2d.Sample) {
	// Initialize the graphic context on an pdf document
	dest := pdf2d.NewPdf("L", "mm", "A4")
	gc := pdf2d.NewGraphicContext(dest)
	// Draw sample
	fn, err := sample(gc, "pdf")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", fn, err)
		return
	}
	// Save to pdf only if it doesn't exist because of git
	if _, err = os.Stat(fn); err == nil {
		t.Skipf("Saving %q skipped, as it exists already. (Git would consider it modified.)", fn)
		return
	}
	err = pdf2d.SaveToPdfFile(fn, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", fn, err)
	}
}
