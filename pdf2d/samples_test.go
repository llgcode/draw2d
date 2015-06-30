package pdf2d_test

import (
	"testing"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"

	"github.com/stanim/draw2d.samples"
	"github.com/stanim/draw2d.samples/android"
	"github.com/stanim/draw2d.samples/frameimage"
	"github.com/stanim/draw2d.samples/gopher"
	"github.com/stanim/draw2d.samples/helloworld"
)

func test(t *testing.T, sample draw2d.Sample) {
	// Initialize the graphic context on an RGBA image
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

func TestSampleAndroid(t *testing.T) {
	test(t, android.Main)
}

func TestSampleGopher(t *testing.T) {
	test(t, gopher.Main)
}

func TestSampleHelloWorld(t *testing.T) {
	// Set the global folder for searching fonts
	// The pdf backend needs for every ttf file its corresponding json
	// file which is generated by gofpdf/makefont.
	draw2d.SetFontFolder(samples.Dir("helloworld", "../"))
	test(t, helloworld.Main)
}

func TestSampleFrameImage(t *testing.T) {
	test(t, frameimage.Main)
}