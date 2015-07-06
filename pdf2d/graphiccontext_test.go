package draw2dpdf

import (
	"fmt"
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/jung-kurt/gofpdf"
)

func ExampleGraphicContext() {
	// Initialize the graphic context on a pdf document
	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	gc := NewGraphicContext(pdf)

	// some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineWidth(5)

	// draw something
	gc.MoveTo(10, 10) // should always be called first for a new path
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()
	fmt.Println(gc.LastPoint())

	// draw2dpdf.SaveToPdfFile("example.pdf", pdf)

	// Output:
	// 10 10
}
