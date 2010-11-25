// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"math"
)

func fabs(x float) float {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly fabs(-0)
	}
	return x
}

func cos(f float) float {
	return float(math.Cos(float64(f)))
}
func sin(f float) float {
	return float(math.Sin(float64(f)))
}
func acos(f float) float {
	return float(math.Acos(float64(f)))
}

func atan2(x, y float) float {
	return float(math.Atan2(float64(x), float64(y)))
}

func distance(x1, y1, x2, y2 float) float {
	dx := x2 - x1
	dy := y2 - y1
	return float(math.Sqrt(float64(dx*dx + dy*dy)))
}

func squareDistance(x1, y1, x2, y2 float) float {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}
