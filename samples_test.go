// See also test_test.go

package draw2d_test

import (
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d.samples"
	"github.com/llgcode/draw2d.samples/android"
	"github.com/llgcode/draw2d.samples/frameimage"
	"github.com/llgcode/draw2d.samples/gopher"
	"github.com/llgcode/draw2d.samples/helloworld"
	"github.com/llgcode/draw2d.samples/line"
	"github.com/llgcode/draw2d.samples/linecapjoin"
	"github.com/llgcode/draw2d.samples/postscript"
)

func TestSampleAndroid(t *testing.T) {
	test(t, android.Main)
}

func TestSampleGopher(t *testing.T) {
	test(t, gopher.Main)
}

func TestSampleHelloWorld(t *testing.T) {
	// Set the global folder for searching fonts
	draw2d.SetFontFolder(samples.Dir("helloworld", ""))
	test(t, helloworld.Main)
}

func TestSampleFrameImage(t *testing.T) {
	test(t, frameimage.Main)
}

func TestSampleLine(t *testing.T) {
	test(t, line.Main)
}

func TestSampleLineCap(t *testing.T) {
	test(t, linecapjoin.Main)
}

func TestSamplePostscript(t *testing.T) {
	test(t, postscript.Main)
}
