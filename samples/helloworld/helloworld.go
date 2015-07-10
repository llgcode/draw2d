// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff, Stani Michiels

// Package helloworld displays "Hello World" twice (one rotated) in a
// rounded rectangle.
package helloworld

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/samples"
)

// Main draws "Hello World" and returns the filename. This should only be
// used during testing.
func Main(gc draw2d.GraphicContext, ext string) (string, error) {
	// Draw hello world
	Draw(gc)

	// Return the output filename
	return samples.Output("helloworld", ext), nil
}

// Draw "Hello World"
func Draw(gc draw2d.GraphicContext) {
	// Draw a rounded rectangle using default colors
	draw2d.RoundRect(gc, 5, 5, 292, 205, 10, 10)
	gc.FillStroke()

	// Set the font luximbi.ttf
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic})
	// Set the fill text color to black
	gc.SetFillColor(image.Black)
	gc.SetDPI(72)
	gc.SetFontSize(14)
	// Display Hello World
	gc.FillStringAt("Hello World", 8, 52)

	gc.Save()
	gc.SetFillColor(color.NRGBA{0xFF, 0x33, 0x33, 0xFF})
	gc.Translate(145, 105)
	gc.Rotate(math.Pi / 4)
	gc.FillStringAt("Hello World", 0, 0)
	gc.Restore()
}
