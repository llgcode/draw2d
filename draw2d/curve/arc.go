// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"math"
)

func SegmentArc(t LineTracer, x, y, rx, ry, start, angle, scale float64) {
	end := start + angle
	clockWise := true
	if angle < 0 {
		clockWise = false
	}
	ra := (math.Abs(rx) + math.Abs(ry)) / 2
	da := math.Acos(ra/(ra+0.125/scale)) * 2
	//normalize
	if !clockWise {
		da = -da
	}
	angle = start + da
	var curX, curY float64
	for {
		if (angle < end-da/4) != clockWise {
			curX = x + math.Cos(end)*rx
			curY = y + math.Sin(end)*ry
			return curX, curY
		}
		curX = x + math.Cos(angle)*rx
		curY = y + math.Sin(angle)*ry

		angle += da
		t.LineTo(curX, curY)
	}
	t.LineTo(curX, curY)
}
