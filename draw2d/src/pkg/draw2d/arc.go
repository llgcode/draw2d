package draw2d


func arc(t LineTracer, x, y, rx, ry, start, angle, scale float) {
	end := start + angle
	clockWise := true
	if angle < 0 {
		clockWise = false
	}
	ra := (fabs(rx) + fabs(ry)) / 2
	da := acos(ra/(ra+0.125/scale)) * 2
	//normalize
	if !clockWise {
		da = -da
	}
	angle = start + da
	var curX, curY float
	for {
		if (angle < end-da/4) != clockWise {
			curX = x + cos(end)*rx
			curY = y + sin(end)*ry
			t.LineTo(curX, curY)
			break
		}
		curX = x + cos(angle)*rx
		curY = y + sin(angle)*ry

		angle += da
		t.LineTo(curX, curY)
	}
}
