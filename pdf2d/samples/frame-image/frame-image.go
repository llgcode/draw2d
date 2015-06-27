// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff, Stani Michiels

// Load a png image and rotate it
package main

import (
	"image/color"
	"math"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"
	"github.com/stanim/gofpdf"
)

func main() {
	// Margin between the image and the frame
	const margin = 30
	// Line width od the frame
	const lineWidth = 3

	// Initialize the graphic context on an RGBA image
	dest := gofpdf.New("P", "mm", "A4", "../font")
	dest.AddPage()
	// Size of destination image
	dw, dh := dest.GetPageSize()
	gc := pdf2d.NewGraphicContext(dest)
	// Draw frame
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw2d.RoundRect(gc, lineWidth, lineWidth, dw-lineWidth, dh-lineWidth, 100, 100)
	gc.SetLineWidth(lineWidth)
	gc.FillStroke()

	// load the source image
	source, err := draw2d.LoadFromPngFile("gopher.png")
	if err != nil {
		panic(err)
	}
	// Size of source image
	sw, sh := float64(source.Bounds().Dx()), float64(source.Bounds().Dy())
	// Draw image to fit in the frame
	// TODO Seems to have a transform bug here on draw image
	scale := math.Min((dw-margin*2)/sw, (dh-margin*2)/sh)
	gc.Translate(margin, margin)
	gc.Scale(scale, scale)

	gc.DrawImage(source)

	// Save to pdf
	pdf2d.SaveToPdfFile("frame-image.pdf", dest)
}
