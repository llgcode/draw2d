package draw2dkit

import (
	"github.com/llgcode/draw2d"
	"math"
)

func Droid(gc draw2d.GraphicContext, x, y float64) {
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineWidth(5)

	// head
	gc.ArcTo(x+80, y+70, 50, 50, 180*(math.Pi/180), 360*(math.Pi/180))
	gc.FillStroke()
	gc.MoveTo(x+60, y+25)
	gc.LineTo(x+50, y+10)
	gc.MoveTo(x+100, y+25)
	gc.LineTo(x+110, y+10)
	gc.Stroke()

	// left eye
	draw2d.Circle(gc, x+60, y+45, 5)
	gc.FillStroke()

	// right eye
	draw2d.Circle(gc, x+100, y+45, 5)
	gc.FillStroke()

	// body
	draw2d.RoundedRectangle(gc, x+30, y+75, x+30+100, y+75+90, 10, 10)
	gc.FillStroke()
	draw2d.Rectangle(gc, x+30, y+75, x+30+100, y+75+80)
	gc.FillStroke()

	// left arm
	draw2d.RoundedRectangle(gc, x+5, y+80, x+5+20, y+80+70, 10, 10)
	gc.FillStroke()

	// right arm
	draw2d.RoundedRectangle(gc, x+135, y+80, x+135+20, y+80+70, 10, 10)
	gc.FillStroke()

	// left leg
	draw2d.RoundedRectangle(gc, x+50, y+150, x+50+20, y+150+50, 10, 10)
	gc.FillStroke()

	// right leg
	draw2d.RoundedRectangle(gc, x+90, y+150, x+90+20, y+150+50, 10, 10)
	gc.FillStroke()
}
