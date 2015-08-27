package draw2dkit

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/stanim/draw2d"
)

func TestCircle(t *testing.T) {
	width := 200
	height := 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(img)

	gc.SetStrokeColor(color.NRGBA{255, 255, 255, 255})
	gc.SetFillColor(color.NRGBA{255, 255, 255, 255})
	gc.Clear()

	gc.SetStrokeColor(color.NRGBA{255, 0, 0, 255})
	gc.SetLineWidth(1)

	// Draw a circle
	Circle(gc, 100, 100, 50)
	gc.Stroke()

	draw2d.SaveToPngFile("../output/draw2dkit/TestCircle.png", img)
}
