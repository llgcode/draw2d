// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

// Draw an android avatar to android.png
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"
	"github.com/stanim/gofpdf"
)

func main() {
	// Initialize the graphic context on a pdf document
	pdf := gofpdf.New("P", "mm", "A4", "../font")
	pdf.AddPage()
	gc := pdf2d.NewGraphicContext(pdf)

	// set the fill and stroke color of the droid
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})

	// Draw the droid
	DrawDroid(gc, 10, 10)

	// Save to pdf
	pdf2d.SaveToPdfFile("android.pdf", pdf)
}

func DrawDroid(gc draw2d.GraphicContext, x, y float64) {
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineWidth(5)

	fmt.Println("\nhead")
	gc.MoveTo(x+30, y+70)
	gc.ArcTo(x+80, y+70, 50, 50, 180*(math.Pi/180), 180*(math.Pi/180))
	gc.Close()
	gc.FillStroke()
	gc.MoveTo(x+60, y+25)
	gc.LineTo(x+50, y+10)
	gc.MoveTo(x+100, y+25)
	gc.LineTo(x+110, y+10)
	gc.Stroke()

	fmt.Println("\nleft eye")
	draw2d.Circle(gc, x+60, y+45, 5)
	gc.FillStroke()

	fmt.Println("\nright eye")
	draw2d.Circle(gc, x+100, y+45, 5)
	gc.FillStroke()

	fmt.Println("\nbody")
	draw2d.RoundRect(gc, x+30, y+75, x+30+100, y+75+90, 10, 10)
	gc.FillStroke()
	draw2d.Rect(gc, x+30, y+75, x+30+100, y+75+80)
	gc.FillStroke()

	fmt.Println("\nleft arm")
	draw2d.RoundRect(gc, x+5, y+80, x+5+20, y+80+70, 10, 10)
	gc.FillStroke()

	fmt.Println("\nright arm")
	draw2d.RoundRect(gc, x+135, y+80, x+135+20, y+80+70, 10, 10)
	gc.FillStroke()

	fmt.Println("\nleft leg")
	draw2d.RoundRect(gc, x+50, y+150, x+50+20, y+150+50, 10, 10)
	gc.FillStroke()

	fmt.Println("\nright leg")
	draw2d.RoundRect(gc, x+90, y+150, x+90+20, y+150+50, 10, 10)
	gc.FillStroke()

}
