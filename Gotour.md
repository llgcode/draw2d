# How to use draw2d in gotour #


Install the gotour program locally, install draw2d and run the gotour executable:
```
go get code.google.com/p/go-tour/gotour
go get code.google.com/p/draw2d/draw2d
gotour
```

now open the gotour into your browser http://127.0.0.1:3999 and run this snippet of code:

```
package main

import (
	"code.google.com/p/draw2d/draw2d"
        "image"
	"image/color"
	"code.google.com/p/go-tour/pic"
)

func main() {
	i := image.NewRGBA(image.Rect(0, 0, 200, 200))
	gc := draw2d.NewGraphicContext(i)
	gc.Save()
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.White)
	draw2d.Rect(gc, 10, 10, 100, 100)
	gc.FillStroke()
	gc.Restore()
	pic.ShowImage(i)
}
```