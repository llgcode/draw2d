// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

// Benchmark tests for draw2d performance issues
// Run with: go test -bench=. -benchmem

package draw2d_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

// BenchmarkFillStrokeRectangle benchmarks the performance of drawing a filled and stroked rectangle.
// Issue #147 reports that draw2d is 10-30x slower than Cairo for similar operations.
// This benchmark helps quantify the performance characteristics.
func BenchmarkFillStrokeRectangle(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	ctx := draw2dimg.NewGraphicContext(img)
	
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
		ctx.SetFillColor(color.RGBA{0x4d, 0x4d, 0x4d, 0xff})
		ctx.SetLineWidth(2)
		ctx.MoveTo(1, 1)
		ctx.LineTo(499, 1)
		ctx.LineTo(499, 499)
		ctx.LineTo(1, 499)
		ctx.Close()
		ctx.FillStroke()
	}
}

// BenchmarkStrokeSimpleLine benchmarks a simple line stroke operation.
func BenchmarkStrokeSimpleLine(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	ctx := draw2dimg.NewGraphicContext(img)
	
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
		ctx.SetLineWidth(2)
		ctx.MoveTo(10, 10)
		ctx.LineTo(490, 490)
		ctx.Stroke()
	}
}

// BenchmarkFillCircle benchmarks filling a circle.
func BenchmarkFillCircle(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	ctx := draw2dimg.NewGraphicContext(img)
	
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx.SetFillColor(color.RGBA{0x00, 0xff, 0x00, 0xff})
		ctx.ArcTo(250, 250, 100, 100, 0, -6.28318530718) // full circle
		ctx.Close()
		ctx.Fill()
	}
}

// BenchmarkMatrixTransform benchmarks matrix transformation operations.
func BenchmarkMatrixTransform(b *testing.B) {
	m := draw2d.NewTranslationMatrix(10, 20)
	points := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Transform(points)
	}
}

// BenchmarkPathConstruction benchmarks path building operations.
func BenchmarkPathConstruction(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p := new(draw2d.Path)
		p.MoveTo(0, 0)
		p.LineTo(100, 0)
		p.LineTo(100, 100)
		p.LineTo(0, 100)
		p.Close()
	}
}
