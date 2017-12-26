// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednářpackage draw2dsvg

package draw2dsvg

import (
	"fmt"
	"github.com/llgcode/draw2d"
	"image/color"
	"math"
	"strings"
)

func toSvgRGBA(c color.Color) string {
	r, g, b, a := c.RGBA()
	r, g, b, a = r>>8, g>>8, b>>8, a>>8
	if a == 255 {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return fmt.Sprintf("rgba(%v,%v,%v,%.3f)", r, g, b, float64(a)/255)
}

func toSvgLength(l float64) string {
	return fmt.Sprintf("%.4f", l)
}

func toSvgArray(nums []float64) string {
	arr := make([]string, len(nums))
	for i, num := range nums {
		arr[i] = fmt.Sprintf("%.4f", num)
	}
	return strings.Join(arr, ",")
}

func toSvgFillRule(rule draw2d.FillRule) string {
	return map[draw2d.FillRule]string{
		draw2d.FillRuleEvenOdd: "evenodd",
		draw2d.FillRuleWinding: "nonzero",
	}[rule]
}

func toSvgPathDesc(p *draw2d.Path) string {
	parts := make([]string, len(p.Components))
	ps := p.Points
	for i, cmp := range p.Components {
		switch cmp {
		case draw2d.MoveToCmp:
			parts[i] = fmt.Sprintf("M %.4f,%.4f", ps[0], ps[1])
			ps = ps[2:]
		case draw2d.LineToCmp:
			parts[i] = fmt.Sprintf("L %.4f,%.4f", ps[0], ps[1])
			ps = ps[2:]
		case draw2d.QuadCurveToCmp:
			parts[i] = fmt.Sprintf("Q %.4f,%.4f %.4f,%.4f", ps[0], ps[1], ps[2], ps[3])
			ps = ps[4:]
		case draw2d.CubicCurveToCmp:
			parts[i] = fmt.Sprintf("C %.4f,%.4f %.4f,%.4f %.4f,%.4f", ps[0], ps[1], ps[2], ps[3], ps[4], ps[5])
			ps = ps[6:]
		case draw2d.ArcToCmp:
			cx, cy := ps[0], ps[1] // center
			rx, ry := ps[2], ps[3] // radii
			fi := ps[4] + ps[5]    // startAngle + angle

			// compute endpoint
			sinfi, cosfi := math.Sincos(fi)
			nom := math.Hypot(ry*cosfi, rx*sinfi)
			x := cx + (rx*ry*cosfi)/nom
			y := cy + (rx*ry*sinfi)/nom

			// compute large and sweep flags
			large := 0
			sweep := 0
			if math.Abs(ps[5]) > math.Pi {
				large = 1
			}
			if !math.Signbit(ps[5]) {
				sweep = 1
			}
			// dirty hack to ensure whole arc is drawn
			// if start point equals end point
			if sweep == 1 {
				x += 0.001 * sinfi
				y += 0.001 * -cosfi
			} else {
				x += 0.001 * sinfi
				y += 0.001 * cosfi
			}

			// rx ry x-axis-rotation large-arc-flag sweep-flag x y
			parts[i] = fmt.Sprintf("A %.4f %.4f %v %v %v %.4f %.4f",
				rx, ry, 0, large, sweep, x, y,
			)
			ps = ps[6:]
		case draw2d.CloseCmp:
			parts[i] = "Z"
		}
	}
	return strings.Join(parts, " ")
}

func toSvgTransform(mat draw2d.Matrix) string {
	if mat.IsIdentity() {
		return ""
	}
	if mat.IsTranslation() {
		x, y := mat.GetTranslation()
		return fmt.Sprintf("translate(%f,%f)", x, y)
	}
	return fmt.Sprintf("matrix(%f,%f,%f,%f,%f,%f)",
		mat[0], mat[1], mat[2], mat[3], mat[4], mat[5],
	)
}
