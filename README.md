draw2d
======

Package draw2d is a pure [go](http://golang.org) 2D vector graphics library with support for multiple output devices such as [images](http://golang.org/pkg/image) (draw2d), pdf documents (draw2dpdf) and opengl (draw2dopengl), which can also be used on the google app engine.
This library is inspired by [postscript](http://www.tailrecursive.org/postscript) and [HTML5 canvas](http://www.w3.org/TR/2dcontext/).

See the [documentation](http://godoc.org/github.com/llgcode/draw2d) for more details.

The package depends on [freetype-go](http://code.google.com/p/freetype-go) package for its rasterization algorithm.

Installation
------------

Install [golang](http://golang.org/doc/install). To install or update the package draw2d on your system, run:

```
go get -u github.com/llgcode/draw2d
```

and start coding using one of the [Samples](https://github.com/llgcode/draw2d.samples).


Softwares and Packages using draw2d
-----------------------------------

 - [golang postscript interpreter](https://github.com/llgcode/ps)
 - [gonum plot](https://github.com/gonum/plot)

References
---------

 - [antigrain.com](http://www.antigrain.com)
 - [freetype-go](http://code.google.com/p/freetype-go)
