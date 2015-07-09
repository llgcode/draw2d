// Package draw2d_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
package draw2d_test

import (
	"image"
	"testing"

	"github.com/llgcode/draw2d"
)

type sample func(gc draw2d.GraphicContext, ext string) (string, error)

func test(t *testing.T, draw sample) {
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
	gc := draw2d.NewGraphicContext(dest)
	// Draw Android logo
	fn, err := draw(gc, "png")
	if err != nil {
		t.Errorf("Drawing %q failed: %v", fn, err)
		return
	}
	// Save to png
	err = draw2d.SaveToPngFile(fn, dest)
	if err != nil {
		t.Errorf("Saving %q failed: %v", fn, err)
	}
}
