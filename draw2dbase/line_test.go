// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 07/02/2026 by draw2d contributors

package draw2dbase

import (
	"image"
	"image/color"
	"testing"
)

func TestBresenham_Horizontal(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{255, 0, 0, 255}

	// Draw horizontal line from (5, 10) to (15, 10)
	Bresenham(img, c, 5, 10, 15, 10)

	// Verify pixels along the line are set
	for x := 5; x <= 15; x++ {
		pixel := img.At(x, 10)
		if pixel != c {
			t.Errorf("Pixel at (%d, 10) not set correctly", x)
		}
	}

	// Verify a pixel off the line is not set
	pixel := img.At(5, 5)
	if pixel == c {
		t.Error("Pixel off the line should not be set")
	}
}

func TestBresenham_Vertical(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{0, 255, 0, 255}

	// Draw vertical line from (10, 5) to (10, 15)
	Bresenham(img, c, 10, 5, 10, 15)

	// Verify pixels along the line are set
	for y := 5; y <= 15; y++ {
		pixel := img.At(10, y)
		if pixel != c {
			t.Errorf("Pixel at (10, %d) not set correctly", y)
		}
	}

	// Verify a pixel off the line is not set
	pixel := img.At(5, 10)
	if pixel == c {
		t.Error("Pixel off the line should not be set")
	}
}

func TestBresenham_Diagonal(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{0, 0, 255, 255}

	// Draw diagonal line from (5, 5) to (15, 15)
	Bresenham(img, c, 5, 5, 15, 15)

	// Verify start and end pixels are set
	if img.At(5, 5) != c {
		t.Error("Start pixel (5, 5) not set")
	}

	if img.At(15, 15) != c {
		t.Error("End pixel (15, 15) not set")
	}

	// Verify a point along the diagonal is set
	if img.At(10, 10) != c {
		t.Error("Middle pixel (10, 10) not set")
	}
}

func TestBresenham_SinglePoint(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{255, 255, 0, 255}

	// Draw from (5, 5) to (5, 5) - single point
	Bresenham(img, c, 5, 5, 5, 5)

	// Verify the single pixel is set
	if img.At(5, 5) != c {
		t.Error("Single point (5, 5) not set")
	}

	// Verify adjacent pixels are not set
	if img.At(6, 5) == c {
		t.Error("Adjacent pixel should not be set")
	}
}

func TestBresenham_ReverseDirection(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{255, 0, 255, 255}

	// Draw from (10, 10) to (0, 0) - reverse direction
	Bresenham(img, c, 10, 10, 0, 0)

	// Verify start and end pixels are set
	if img.At(10, 10) != c {
		t.Error("Start pixel (10, 10) not set")
	}

	if img.At(0, 0) != c {
		t.Error("End pixel (0, 0) not set")
	}

	// Verify a point along the line is set
	if img.At(5, 5) != c {
		t.Error("Middle pixel (5, 5) not set")
	}
}

func TestPolylineBresenham(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 30, 30))
	c := color.RGBA{128, 128, 128, 255}

	// Draw polyline with three segments: (5,5) -> (15,5) -> (15,15) -> (5,15)
	points := []float64{5, 5, 15, 5, 15, 15, 5, 15}
	PolylineBresenham(img, c, points...)

	// Verify key pixels are set
	if img.At(5, 5) != c {
		t.Error("Start pixel (5, 5) not set")
	}

	if img.At(15, 5) != c {
		t.Error("Corner pixel (15, 5) not set")
	}

	if img.At(15, 15) != c {
		t.Error("Corner pixel (15, 15) not set")
	}

	if img.At(5, 15) != c {
		t.Error("End pixel (5, 15) not set")
	}

	// Verify a pixel along the first segment is set
	if img.At(10, 5) != c {
		t.Error("Pixel (10, 5) along first segment not set")
	}
}

func TestPolylineBresenham_TwoPoints(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	c := color.RGBA{64, 64, 64, 255}

	// Draw single segment via polyline: (5,5) -> (10,10)
	points := []float64{5, 5, 10, 10}
	PolylineBresenham(img, c, points...)

	// Verify endpoints are set
	if img.At(5, 5) != c {
		t.Error("Start pixel (5, 5) not set")
	}

	if img.At(10, 10) != c {
		t.Error("End pixel (10, 10) not set")
	}
}
