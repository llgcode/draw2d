// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

func minMax(x, y float64) (min, max float64) {
	if x > y {
		return y, x
	}
	return x, y
}
