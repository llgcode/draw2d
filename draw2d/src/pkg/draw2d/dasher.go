package draw2d

import(
	"freetype-go.googlecode.com/hg/freetype/raster"	
)

type DashAdder struct {
	adder          raster.Adder
	x, y, distance float
	dash           []float
	currentDash    int
	dashOffset     float
}

func traceDashPath(dash []float, dashOffset float, approximationScale float, adder raster.Adder, paths ...*PathStorage) {
	var dasher DashAdder
	if dash != nil && len(dash) > 0 {
		dasher.dash = dash
	} else {
		dasher.dash = nil
	}
	dasher.currentDash = 0
	dasher.dashOffset = dashOffset
	dasher.adder = adder
	for _, path := range paths {
		path.TraceLine(&dasher, approximationScale)
	}
}

func (dasher *DashAdder) MoveTo(x, y float) {
	dasher.adder.Start(floatToPoint(x, y))
	dasher.x, dasher.y = x, y
	dasher.distance = dasher.dashOffset
	dasher.currentDash = 0
}

func (dasher *DashAdder) LineTo(x, y float) {
	rest := dasher.dash[dasher.currentDash] - dasher.distance
	for rest < 0 {
		dasher.distance = dasher.distance - dasher.dash[dasher.currentDash]
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
		rest = dasher.dash[dasher.currentDash] - dasher.distance
	}
	d := distance(dasher.x, dasher.y, x, y)
	for d >= rest {
		k := rest / d
		lx := dasher.x + k*(x-dasher.x)
		ly := dasher.y + k*(y-dasher.y)
		if dasher.currentDash%2 == 0 {
			// line
			dasher.adder.Add1(floatToPoint(lx, ly))
		} else {
			// gap
			dasher.adder.Start(floatToPoint(lx, ly))
		}
		d = d - rest
		dasher.x, dasher.y = lx, ly
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
		rest = dasher.dash[dasher.currentDash]
	}
	dasher.distance = d
	if dasher.currentDash%2 == 0 {
		dasher.adder.Add1(floatToPoint(x, y))
	} else {
		dasher.adder.Start(floatToPoint(x, y))
	}
	if dasher.distance >= dasher.dash[dasher.currentDash] {
		dasher.distance = dasher.distance - dasher.dash[dasher.currentDash]
		dasher.currentDash = (dasher.currentDash + 1) % len(dasher.dash)
	}
	dasher.x, dasher.y = x, y
}
