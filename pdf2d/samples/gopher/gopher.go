// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff, Stani Michiels

// Draw a gopher avatar to gopher.png translating from this svg https://github.com/golang-samples/gopher-vector/
package main

import (
	"image/color"

	"github.com/stanim/draw2d"
	"github.com/stanim/draw2d/pdf2d"
)

func main() {
	// Initialize the graphic context on an RGBA image
	dest := pdf2d.NewPdf("P", "mm", "A4")
	gc := pdf2d.NewGraphicContext(dest)

	// Draw a gopher
	Gopher(gc, 48, 48, 240, 72)

	// Save to pdf
	pdf2d.SaveToPdfFile("gopher.pdf", dest)
}

// Gopher draw a gopher using a gc thanks to https://github.com/golang-samples/gopher-vector/
func Gopher(gc draw2d.GraphicContext, x, y, w, h float64) {
	// Initialize Stroke Attribute
	gc.SetLineWidth(3)
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetStrokeColor(color.Black)

	// Left hand
	// <path fill-rule="evenodd" clip-rule="evenodd" fill="#F6D2A2" stroke="#000000" stroke-width="3" stroke-linecap="round" d="
	// M10.634,300.493c0.764,15.751,16.499,8.463,23.626,3.539c6.765-4.675,8.743-0.789,9.337-10.015
	// c0.389-6.064,1.088-12.128,0.744-18.216c-10.23-0.927-21.357,1.509-29.744,7.602C10.277,286.542,2.177,296.561,10.634,300.493"/>
	gc.SetFillColor(color.RGBA{0xF6, 0xD2, 0xA2, 0xff})
	gc.MoveTo(10.634, 300.493)
	gc.RCubicCurveTo(0.764, 15.751, 16.499, 8.463, 23.626, 3.539)
	gc.RCubicCurveTo(6.765, -4.675, 8.743, -0.789, 9.337, -10.015)
	gc.RCubicCurveTo(0.389, -6.064, 1.088, -12.128, 0.744, -18.216)
	gc.RCubicCurveTo(-10.23, -0.927, -21.357, 1.509, -29.744, 7.602)
	gc.CubicCurveTo(10.277, 286.542, 2.177, 296.561, 10.634, 300.493)
	gc.FillStroke()

	// <path fill-rule="evenodd" clip-rule="evenodd" fill="#C6B198" stroke="#000000" stroke-width="3" stroke-linecap="round" d="
	// M10.634,300.493c2.29-0.852,4.717-1.457,6.271-3.528"/>
	gc.MoveTo(10.634, 300.493)
	gc.RCubicCurveTo(2.29, -0.852, 4.717, -1.457, 6.271, -3.528)
	gc.Stroke()

	// Left Ear
	// <path fill-rule="evenodd" clip-rule="evenodd" fill="#6AD7E5" stroke="#000000" stroke-width="3" stroke-linecap="round" d="
	// M46.997,112.853C-13.3,95.897,31.536,19.189,79.956,50.74L46.997,112.853z"/>
	gc.MoveTo(46.997, 112.853)
	gc.CubicCurveTo(-13.3, 95.897, 31.536, 19.189, 79.956, 50.74)
	gc.LineTo(46.997, 112.853)
	gc.Close()
	gc.Stroke()
}
