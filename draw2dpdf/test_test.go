// Package draw2dpdf_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package draw2dpdf_test

import (
	"os"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
)

func test(t *testing.T, sample draw2d.Sample) {
	// Initialize the graphic context on an pdf document
	dest := draw2dpdf.NewPdf("L", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(dest)
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
	err = draw2dpdf.SaveToPdfFile(fn, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", fn, err)
	}
}
