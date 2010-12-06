// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package main

import (
	"draw2d"
	//"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
	"fmt"
)

func main() {
	path := new(draw2d.Path)
	path.MoveTo(2.0, 3.0)
	path.LineTo(2.0, 3.0)
	path.QuadCurveTo(2.0, 3.0, 10, 20)
	path.CubicCurveTo(2.0, 3.0, 10, 20, 13, 23)
	path.ArcTo(2.0, 3.0, 100, 200, 200, 300)
	fmt.Printf("%v\n", path)
}
