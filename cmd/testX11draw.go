// +build ignore

package main

import (
	"exp/gui"
	"exp/gui/x11"
	"fmt"
	"image"
	"math"

	"github.com/llgcode/draw2d"
)

func main() {
	window, err := x11.NewWindow()
	if err != nil {
		fmt.Printf("Cannot open an x11 window\n")
		return
	}
	screen := window.Screen()
	gc := draw2d.NewGraphicContext(screen)
	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	gc.Clear()
	for i := 0.0; i < 360; i = i + 10 { // Go from 0 to 360 degrees in 10 degree steps
		gc.BeginPath() // Start a new path
		gc.Save()      // Keep rotations temporary
		gc.MoveTo(144, 144)
		gc.Rotate(i * (math.Pi / 180.0)) // Rotate by degrees on stack from 'for'
		gc.RLineTo(72, 0)
		gc.Stroke()
		gc.Restore() // Get back the unrotated state
	}

	window.FlushImage()

	gc.SetLineWidth(3)
	nbclick := 0
	for {

		switch evt := (<-window.EventChan()).(type) {
		case gui.KeyEvent:
			if evt.Key == 'q' {
				window.Close()
			}
		case gui.MouseEvent:
			if evt.Buttons&1 != 0 {
				if nbclick%2 == 0 {
					gc.MoveTo(float64(evt.Loc.X), float64(evt.Loc.Y))
				} else {
					gc.LineTo(float64(evt.Loc.X), float64(evt.Loc.Y))
					gc.Stroke()
					window.FlushImage()
				}
				nbclick = nbclick + 1
			}
		}
	}
}
