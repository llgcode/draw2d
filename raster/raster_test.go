package raster

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"

	"code.google.com/p/freetype-go/freetype/raster"
	"gopkg.in/llgcode/draw2d.v1/curve"
)

var flatteningThreshold = 0.5

func savepng(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type Path struct {
	points []float64
}

func (p *Path) LineTo(x, y float64) {
	if len(p.points)+2 > cap(p.points) {
		points := make([]float64, len(p.points)+2, len(p.points)+32)
		copy(points, p.points)
		p.points = points
	} else {
		p.points = p.points[0 : len(p.points)+2]
	}
	p.points[len(p.points)-2] = x
	p.points[len(p.points)-1] = y
}

func TestFreetype(t *testing.T) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	rasterizer := raster.NewRasterizer(200, 200)
	rasterizer.UseNonZeroWinding = false
	rasterizer.Start(raster.Point{
		X: raster.Fix32(10 * 256),
		Y: raster.Fix32(190 * 256)})
	for j := 0; j < len(poly); j = j + 2 {
		rasterizer.Add1(raster.Point{
			X: raster.Fix32(poly[j] * 256),
			Y: raster.Fix32(poly[j+1] * 256)})
	}
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	rasterizer.Rasterize(painter)

	savepng("../output/raster/TestFreetype.png", img)
}

func TestFreetypeNonZeroWinding(t *testing.T) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	rasterizer := raster.NewRasterizer(200, 200)
	rasterizer.UseNonZeroWinding = true
	rasterizer.Start(raster.Point{
		X: raster.Fix32(10 * 256),
		Y: raster.Fix32(190 * 256)})
	for j := 0; j < len(poly); j = j + 2 {
		rasterizer.Add1(raster.Point{
			X: raster.Fix32(poly[j] * 256),
			Y: raster.Fix32(poly[j+1] * 256)})
	}
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	rasterizer.Rasterize(painter)

	savepng("../output/raster/TestFreetypeNonZeroWinding.png", img)
}

func TestRasterizer(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}
	tr := [6]float64{1, 0, 0, 1, 0, 0}
	r := NewRasterizer8BitsSample(200, 200)
	//PolylineBresenham(img, image.Black, poly...)

	r.RenderEvenOdd(img, &color, &poly, tr)
	savepng("../output/raster/TestRasterizer.png", img)
}

func TestRasterizerNonZeroWinding(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}
	tr := [6]float64{1, 0, 0, 1, 0, 0}
	r := NewRasterizer8BitsSample(200, 200)
	//PolylineBresenham(img, image.Black, poly...)

	r.RenderNonZeroWinding(img, &color, &poly, tr)
	savepng("../output/raster/TestRasterizerNonZeroWinding.png", img)
}

func BenchmarkFreetype(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := raster.NewRasterizer(200, 200)
		rasterizer.UseNonZeroWinding = false
		rasterizer.Start(raster.Point{
			X: raster.Fix32(10 * 256),
			Y: raster.Fix32(190 * 256)})
		for j := 0; j < len(poly); j = j + 2 {
			rasterizer.Add1(raster.Point{
				X: raster.Fix32(poly[j] * 256),
				Y: raster.Fix32(poly[j+1] * 256)})
		}
		painter := raster.NewRGBAPainter(img)
		painter.SetColor(color)
		rasterizer.Rasterize(painter)
	}
}

func BenchmarkFreetypeNonZeroWinding(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}

	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := raster.NewRasterizer(200, 200)
		rasterizer.UseNonZeroWinding = true
		rasterizer.Start(raster.Point{
			X: raster.Fix32(10 * 256),
			Y: raster.Fix32(190 * 256)})
		for j := 0; j < len(poly); j = j + 2 {
			rasterizer.Add1(raster.Point{
				X: raster.Fix32(poly[j] * 256),
				Y: raster.Fix32(poly[j+1] * 256)})
		}
		painter := raster.NewRGBAPainter(img)
		painter.SetColor(color)
		rasterizer.Rasterize(painter)
	}
}

func BenchmarkRasterizerNonZeroWinding(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}
	tr := [6]float64{1, 0, 0, 1, 0, 0}
	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := NewRasterizer8BitsSample(200, 200)
		rasterizer.RenderNonZeroWinding(img, &color, &poly, tr)
	}
}

func BenchmarkRasterizer(b *testing.B) {
	var p Path
	p.LineTo(10, 190)
	c := curve.CubicCurveFloat64{10, 190, 10, 10, 190, 10, 190, 190}
	c.Segment(&p, flatteningThreshold)
	poly := Polygon(p.points)
	color := color.RGBA{0, 0, 0, 0xff}
	tr := [6]float64{1, 0, 0, 1, 0, 0}
	for i := 0; i < b.N; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 200, 200))
		rasterizer := NewRasterizer8BitsSample(200, 200)
		rasterizer.RenderEvenOdd(img, &color, &poly, tr)
	}
}
