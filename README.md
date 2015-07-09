draw2d
======

Package draw2d is a pure [go](http://golang.org) 2D vector graphics library with support for multiple output devices such as [images](http://golang.org/pkg/image) (draw2d), pdf documents (draw2dpdf) and opengl (draw2dopengl), which can also be used on the google app engine. It can be used as a pure go [Cairo](http://www.cairographics.org/) alternative.

See the [documentation](http://godoc.org/github.com/llgcode/draw2d) for more details.

Features
--------

Operations in draw2d include stroking and filling polygons, arcs, Bézier curves, drawing images and text rendering with truetype fonts. All drawing operations can be transformed by affine transformations (scale, rotation, translation).

Installation
------------

Install [golang](http://golang.org/doc/install). To install or update the package draw2d on your system, run:

```
go get -u github.com/llgcode/draw2d
```

Quick Start
-----------

The following Go code generates a simple drawing and saves it to an image file with package draw2d:

```go
// Initialize the graphic context on an RGBA image
dest := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
gc := draw2d.NewGraphicContext(dest)

// Set some properties
gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
gc.SetLineWidth(5)

// Draw a closed shape
gc.MoveTo(10, 10) // should always be called first for a new path
gc.LineTo(100, 50)
gc.QuadCurveTo(100, 10, 10, 10)
gc.Close()
gc.FillStroke()

// Save to file
draw2d.SaveToPngFile(fn, dest)
```

The same Go code can also generate a pdf document with package draw2dpdf:

```go
// Initialize the graphic context on an RGBA image
dest := draw2dpdf.NewPdf("L", "mm", "A4")
gc := draw2d.NewGraphicContext(dest)

// Set some properties
gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
gc.SetLineWidth(5)

// Draw a closed shape
gc.MoveTo(10, 10) // should always be called first for a new path
gc.LineTo(100, 50)
gc.QuadCurveTo(100, 10, 10, 10)
gc.Close()
gc.FillStroke()

// Save to file
draw2dpdf.SaveToPdfFile(fn, dest)
```

There are more examples here: https://github.com/llgcode/draw2d.samples

Drawing on opengl is provided by the draw2dgl package.

Acknowledgments
---------------

[Laurent Le Goff](https://github.com/llgcode) wrote this library, inspired by [Postscript](http://www.tailrecursive.org/postscript) and [HTML5 canvas](http://www.w3.org/TR/2dcontext/). He implemented the image and opengl backend with the [freetype-go](https://code.google.com/p/freetype-go/) package. Also he created a pure go [Postscript interpreter](https://github.com/llgcode/ps), which can read postscript images and draw to a draw2d graphic context. [Stani Michiels](https://github.com/stanim) implemented the pdf backend with the [gofpdf](https://github.com/jung-kurt/gofpdf) package.



Packages using draw2d
---------------------

 - [ps](https://github.com/llgcode/ps): Postscript interpreter written in Go
 - [gonum/plot](https://github.com/gonum/plot): drawing plots in Go
 - [go.uik](https://github.com/skelterjohn/go.uik): a concurrent UI kit written in pure go.
 - [smartcrop](https://github.com/muesli/smartcrop): content aware image cropping
 - [karta](https://github.com/peterhellberg/karta): drawing Voronoi diagrams
 - [chart](https://github.com/vdobler/chart): basic charts in Go

References
---------

 - [antigrain.com](http://www.antigrain.com)
 - [freetype-go](http://code.google.com/p/freetype-go)
