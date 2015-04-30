package draw2dkit

import (
	"github.com/llgcode/draw2d"
	"math"
)

// Droid draws a droid at specified position
func Droid(drawer draw2d.Drawer, x, y float64, fillStyle draw2d.FillStyle, strokeStyle draw2d.StrokeStyle) {
	strokeStyle.LineCap = draw2d.RoundCap
	strokeStyle.Width = 5

	path := &draw2d.Path{}

	// head
	path.ArcTo(x+80, y+70, 50, 50, 180*(math.Pi/180), 360*(math.Pi/180))
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	path.Clear()
	path.MoveTo(x+60, y+25)
	path.LineTo(x+50, y+10)
	path.MoveTo(x+100, y+25)
	path.LineTo(x+110, y+10)
	drawer.Stroke(path, strokeStyle)

	// left eye
	path.Clear()
	Circle(path, x+60, y+45, 5)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// right eye
	path.Clear()
	Circle(path, x+100, y+45, 5)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// body
	path.Clear()
	RoundedRectangle(path, x+30, y+75, x+30+100, y+75+90, 10, 10)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	path.Clear()
	Rectangle(path, x+30, y+75, x+30+100, y+75+80)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// left arm
	path.Clear()
	RoundedRectangle(path, x+5, y+80, x+5+20, y+80+70, 10, 10)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// right arm
	path.Clear()
	RoundedRectangle(path, x+135, y+80, x+135+20, y+80+70, 10, 10)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// left leg
	path.Clear()
	RoundedRectangle(path, x+50, y+150, x+50+20, y+150+50, 10, 10)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)

	// right leg
	path.Clear()
	RoundedRectangle(path, x+90, y+150, x+90+20, y+150+50, 10, 10)
	drawer.Fill(path, fillStyle)
	drawer.Stroke(path, strokeStyle)
}
