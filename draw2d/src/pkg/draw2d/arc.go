// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff
package draw2d

import (
	"freetype-go.googlecode.com/hg/freetype/raster"
)

func arc(t VertexConverter, x, y, rx, ry, start, angle, scale float) (lastX, lastY float) {
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
			return curX, curY
		}
		curX = x + cos(angle)*rx
		curY = y + sin(angle)*ry

		angle += da
		t.Vertex(curX, curY)
	}
	return curX, curY
}


func arcAdder(adder raster.Adder, x, y, rx, ry, start, angle, scale float) raster.Point {
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
			return floatToPoint(curX, curY)
		}
		curX = x + cos(angle)*rx
		curY = y + sin(angle)*ry

		angle += da
		adder.Add1(floatToPoint(curX, curY))
	}
	return floatToPoint(curX, curY)
}
