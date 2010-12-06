package draw2d

import(
	"math"
)


//high level path creation
func Rect(path Path, x1, y1, x2, y2 float) {
	path.MoveTo(x1, y1)
	path.LineTo(x2, y1)
	path.LineTo(x2, y2)
	path.LineTo(x1, y2)
	path.Close()
}

func RoundRect(path Path, x1, y1, x2, y2, arcWidth, arcHeight float) {
	arcWidth = arcWidth/2;
	arcHeight = arcHeight/2;
	path.MoveTo(x1, y1+ arcHeight);
	path.QuadCurveTo(x1, y1, x1 + arcWidth, y1);
	path.LineTo(x2-arcWidth, y1);
	path.QuadCurveTo(x2, y1, x2, y1 + arcHeight);
	path.LineTo(x2, y2-arcHeight);
	path.QuadCurveTo(x2, y2, x2 - arcWidth, y2);
	path.LineTo(x1 + arcWidth, y2);
	path.QuadCurveTo(x1, y2, x1, y2 - arcHeight);
	path.Close()
}

func Ellipse(path Path, cx, cy, rx, ry float) {
	path.ArcTo(cx, cy, rx, ry, 0, -math.Pi * 2)
	path.Close()
}

func Circle(path Path, cx, cy, radius float) {
	path.ArcTo(cx, cy, radius, radius, 0, -math.Pi * 2)
	path.Close()
}