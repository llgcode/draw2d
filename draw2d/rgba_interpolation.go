package draw2d

import (
	"exp/draw"
	"image"
	"math"
)

type ImageFilter int

const (
	LinearFilter ImageFilter = iota
	BilinearFilter
	BicubicFilter
)

//see http://pippin.gimp.org/image_processing/chap_resampling.html
func getColorLinear(img image.Image, x, y float64) image.Color {
	return img.At(int(x), int(y))
}

func getColorBilinear(img image.Image, x, y float64) image.Color {
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	dx := x - x0
	dy := y - y0

	color0 := img.At(int(x0), int(y0))
	color1 := img.At(int(x0+1), int(y0))
	color2 := img.At(int(x0+1), int(y0+1))
	color3 := img.At(int(x0), int(y0+1))

	return lerp(lerp(color0, color1, dx), lerp(color3, color2, dx), dy)
}
/**
-- LERP
-- /lerp/, vi.,n.
--
-- Quasi-acronym for Linear Interpolation, used as a verb or noun for
-- the operation. "Bresenham's algorithm lerps incrementally between the
-- two endpoints of the line." (From Jargon File (4.4.4, 14 Aug 2003)
*/
func lerp(c1, c2 image.Color, ratio float64) image.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r := int(float64(r1)*(1-ratio) + float64(r2)*ratio)
	g := int(float64(g1)*(1-ratio) + float64(g2)*ratio)
	b := int(float64(b1)*(1-ratio) + float64(b2)*ratio)
	a := int(float64(a1)*(1-ratio) + float64(a2)*ratio)
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}


func getColorCubicRow(img image.Image, x, y, offset float64) image.Color {
	c0 := img.At(int(x), int(y))
	c1 := img.At(int(x+1), int(y))
	c2 := img.At(int(x+2), int(y))
	c3 := img.At(int(x+3), int(y))
	rt, gt, bt, at := c0.RGBA()
	r0, g0, b0, a0 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c1.RGBA()
	r1, g1, b1, a1 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c2.RGBA()
	r2, g2, b2, a2 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c3.RGBA()
	r3, g3, b3, a3 := float64(rt), float64(gt), float64(bt), float64(at)
	r, g, b, a := cubic(offset, r0, r1, r2, r3), cubic(offset, g0, g1, g2, g3), cubic(offset, b0, b1, b2, b3), cubic(offset, a0, a1, a2, a3)
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func getColorBicubic(img image.Image, x, y float64) image.Color {
	x0 := math.Floor(x)
	y0 := math.Floor(y)
	dx := x - x0
	dy := y - y0
	c0 := getColorCubicRow(img, x0-1, y0-1, dx)
	c1 := getColorCubicRow(img, x0-1, y0, dx)
	c2 := getColorCubicRow(img, x0-1, y0+1, dx)
	c3 := getColorCubicRow(img, x0-1, y0+2, dx)
	rt, gt, bt, at := c0.RGBA()
	r0, g0, b0, a0 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c1.RGBA()
	r1, g1, b1, a1 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c2.RGBA()
	r2, g2, b2, a2 := float64(rt), float64(gt), float64(bt), float64(at)
	rt, gt, bt, at = c3.RGBA()
	r3, g3, b3, a3 := float64(rt), float64(gt), float64(bt), float64(at)
	r, g, b, a := cubic(dy, r0, r1, r2, r3), cubic(dy, g0, g1, g2, g3), cubic(dy, b0, b1, b2, b3), cubic(dy, a0, a1, a2, a3)
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func cubic(offset, v0, v1, v2, v3 float64) uint32 {
	// offset is the offset of the sampled value between v1 and v2
	return uint32(((((-7*v0+21*v1-21*v2+7*v3)*offset+
		(15*v0-36*v1+27*v2-6*v3))*offset+
		(-9*v0+9*v2))*offset + (v0 + 16*v1 + v2)) / 18.0)
}

func compose(c1, c2 image.Color) image.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	ia := M - a2
	r := ((r1 * ia) / M) + r2
	g := ((g1 * ia) / M) + g2
	b := ((b1 * ia) / M) + b2
	a := ((a1 * ia) / M) + a2
	return image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}


func DrawImage(src image.Image, dest draw.Image, tr MatrixTransform, op draw.Op, filter ImageFilter) {
	b := src.Bounds()
	x0, y0, x1, y1 := float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y)
	tr.TransformRectangle(&x0, &y0, &x1, &y1)
	var x, y, u, v float64
	for x = x0; x < x1; x++ {
		for y = y0; y < y1; y++ {
			u = x
			v = y
			tr.InverseTransform(&u, &v)
			c1 := dest.At(int(x), int(y))
			var c2 image.Color
			switch filter {
			case LinearFilter:
				c2 = src.At(int(u), int(v))
			case BilinearFilter:
				c2 = getColorBilinear(src, u, v)
			case BicubicFilter:
				c2 = getColorBicubic(src, u, v)
			}
			var cr image.Color
			switch op {
			case draw.Over:
				r1, g1, b1, a1 := c1.RGBA()
				r2, g2, b2, a2 := c2.RGBA()
				ia := M - a2
				r := ((r1 * ia) / M) + r2
				g := ((g1 * ia) / M) + g2
				b := ((b1 * ia) / M) + b2
				a := ((a1 * ia) / M) + a2
				cr = image.RGBAColor{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			default:
				cr = c2
			}
			dest.Set(int(x), int(y), cr)
		}
	}
}
