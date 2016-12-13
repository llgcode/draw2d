// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

// Package draw2d is a pure go 2D vector graphics library with support
// for multiple output devices such as images (draw2d), pdf documents
// (draw2dpdf) and opengl (draw2dgl), which can also be used on the
// google app engine. It can be used as a pure go Cairo alternative.
// draw2d is released under the BSD license.
//
// Features
//
// Operations in draw2d include stroking and filling polygons, arcs,
// Bézier curves, drawing images and text rendering with truetype fonts.
// All drawing operations can be transformed by affine transformations
// (scale, rotation, translation).
//
// Package draw2d follows the conventions of http://www.w3.org/TR/2dcontext for coordinate system, angles, etc...
//
// Installation
//
// To install or update the package draw2d on your system, run:
//   go get -u gopkg.in/llgcode/draw2d.v1
//
// Quick Start
//
// Package draw2d itself provides a graphic context that can draw vector
// graphics and text on an image canvas. The following Go code
// generates a simple drawing and saves it to an image file:
//   // Initialize the graphic context on an RGBA image
//   dest := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
//   gc := draw2d.NewGraphicContext(dest)
//
//   // Set some properties
//   gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
//   gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
//   gc.SetLineWidth(5)
//
//   // Draw a closed shape
//   gc.MoveTo(10, 10) // should always be called first for a new path
//   gc.LineTo(100, 50)
//   gc.QuadCurveTo(100, 10, 10, 10)
//   gc.Close()
//   gc.FillStroke()
//
//   // Save to file
//   draw2d.SaveToPngFile("hello.png", dest)
//
// There are more examples here:
// https://gopkg.in/llgcode/draw2d.v1/tree/master/samples
//
// Drawing on pdf documents is provided by the draw2dpdf package.
// Drawing on opengl is provided by the draw2dgl package.
// See subdirectories at the bottom of this page.
//
// Testing
//
// The samples are run as tests from the root package folder `draw2d` by:
//   go test ./...
//
// Or if you want to run with test coverage:
//   go test -cover ./... | grep -v "no test"
//
// This will generate output by the different backends in the output folder.
//
// Acknowledgments
//
// Laurent Le Goff wrote this library, inspired by Postscript and
// HTML5 canvas. He implemented the image and opengl backend with the
// freetype-go package. Also he created a pure go Postscript
// interpreter, which can read postscript images and draw to a draw2d
// graphic context (https://github.com/llgcode/ps). Stani Michiels
// implemented the pdf backend with the gofpdf package.
//
// Packages using draw2d
//
// - https://github.com/llgcode/ps: Postscript interpreter written in Go
//
// - https://github.com/gonum/plot: drawing plots in Go
//
// - https://github.com/muesli/smartcrop: content aware image cropping
//
// - https://github.com/peterhellberg/karta: drawing Voronoi diagrams
//
// - https://github.com/vdobler/chart: basic charts in Go
package draw2d
