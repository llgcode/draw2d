// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

// +build ignore

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/llgcode/draw2d"
)

const (
	w, h = 512, 512
)

var (
	lastTime int64
	folder   = "../resource/result/"
)

func initGc(w, h int) (image.Image, draw2d.GraphicContext) {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2d.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	//gc.Clear()

	return i, gc
}

func saveToPngFile(TestName string, m image.Image) {
	filePath := folder + TestName + ".png"
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
	fmt.Printf("Wrote %s OK.\n", filePath)
}

/*
  <img src="../test_results/TestPath.png"/>
*/
func TestPath() {
	i, gc := initGc(w, h)
	gc.Translate(10, 10)
	gc.MoveTo(0.0, 0.0)
	gc.LineTo(100.0, 00.0)
	gc.LineTo(100.0, 100.0)
	gc.LineTo(0.0, 100.0)
	gc.LineTo(0.0, 0.0)
	gc.FillStroke()
	saveToPngFile("TestPath", i)
}

/*
  <img src="../test_results/TestDrawArc.png"/>
*/
func TestDrawArc() {
	i, gc := initGc(w, h)
	// draw an arc
	xc, yc := 128.0, 128.0
	radiusX, radiusY := 100.0, 100.0
	startAngle := 45.0 * (math.Pi / 180.0) /* angles are specified */
	angle := 135 * (math.Pi / 180.0)       /* in radians           */
	gc.SetLineWidth(10)
	gc.SetLineCap(draw2d.ButtCap)
	gc.SetStrokeColor(image.Black)
	gc.ArcTo(xc, yc, radiusX, radiusY, startAngle, angle)
	gc.Stroke()
	// fill a circle
	gc.SetStrokeColor(color.NRGBA{255, 0x33, 0x33, 0x80})
	gc.SetFillColor(color.NRGBA{255, 0x33, 0x33, 0x80})
	gc.SetLineWidth(6)

	gc.MoveTo(xc, yc)
	gc.LineTo(xc+math.Cos(startAngle)*radiusX, yc+math.Sin(startAngle)*radiusY)
	gc.MoveTo(xc, yc)
	gc.LineTo(xc-radiusX, yc)
	gc.Stroke()

	gc.ArcTo(xc, yc, 10.0, 10.0, 0, 2*math.Pi)
	gc.Fill()
	saveToPngFile("TestDrawArc", i)
}

/*
  <img src="../test_results/TestDrawArc.png"/>
*/
func TestDrawArcNegative() {
	i, gc := initGc(w, h)
	// draw an arc
	xc, yc := 128.0, 128.0
	radiusX, radiusY := 100.0, 100.0
	startAngle := 45.0 * (math.Pi / 180.0) /* angles are specified */
	angle := -225 * (math.Pi / 180.0)      /* in radians           */
	gc.SetLineWidth(10)
	gc.SetLineCap(draw2d.ButtCap)
	gc.SetStrokeColor(image.Black)

	gc.ArcTo(xc, yc, radiusX, radiusY, startAngle, angle)
	gc.Stroke()
	// fill a circle
	gc.SetStrokeColor(color.NRGBA{255, 0x33, 0x33, 0x80})
	gc.SetFillColor(color.NRGBA{255, 0x33, 0x33, 0x80})
	gc.SetLineWidth(6)

	gc.MoveTo(xc, yc)
	gc.LineTo(xc+math.Cos(startAngle)*radiusX, yc+math.Sin(startAngle)*radiusY)
	gc.MoveTo(xc, yc)
	gc.LineTo(xc-radiusX, yc)
	gc.Stroke()

	gc.ArcTo(xc, yc, 10.0, 10.0, 0, 2*math.Pi)
	gc.Fill()
	saveToPngFile("TestDrawArcNegative", i)
}

func TestCurveRectangle() {
	i, gc := initGc(w, h)

	/* a custom shape that could be wrapped in a function */
	x0, y0 := 25.6, 25.6 /* parameters like cairo_rectangle */
	rect_width, rect_height := 204.8, 204.8
	radius := 102.4 /* and an approximate curvature radius */

	x1 := x0 + rect_width
	y1 := y0 + rect_height
	if rect_width/2 < radius {
		if rect_height/2 < radius {
			gc.MoveTo(x0, (y0+y1)/2)
			gc.CubicCurveTo(x0, y0, x0, y0, (x0+x1)/2, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, (y0+y1)/2)
			gc.CubicCurveTo(x1, y1, x1, y1, (x1+x0)/2, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, (y0+y1)/2)
		} else {
			gc.MoveTo(x0, y0+radius)
			gc.CubicCurveTo(x0, y0, x0, y0, (x0+x1)/2, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, y0+radius)
			gc.LineTo(x1, y1-radius)
			gc.CubicCurveTo(x1, y1, x1, y1, (x1+x0)/2, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, y1-radius)
		}
	} else {
		if rect_height/2 < radius {
			gc.MoveTo(x0, (y0+y1)/2)
			gc.CubicCurveTo(x0, y0, x0, y0, x0+radius, y0)
			gc.LineTo(x1-radius, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, (y0+y1)/2)
			gc.CubicCurveTo(x1, y1, x1, y1, x1-radius, y1)
			gc.LineTo(x0+radius, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, (y0+y1)/2)
		} else {
			gc.MoveTo(x0, y0+radius)
			gc.CubicCurveTo(x0, y0, x0, y0, x0+radius, y0)
			gc.LineTo(x1-radius, y0)
			gc.CubicCurveTo(x1, y0, x1, y0, x1, y0+radius)
			gc.LineTo(x1, y1-radius)
			gc.CubicCurveTo(x1, y1, x1, y1, x1-radius, y1)
			gc.LineTo(x0+radius, y1)
			gc.CubicCurveTo(x0, y1, x0, y1, x0, y1-radius)
		}
	}
	gc.Close()

	gc.SetFillColor(color.NRGBA{0x80, 0x80, 0xFF, 0xFF})
	gc.SetStrokeColor(color.NRGBA{0x80, 0, 0, 0x80})
	gc.SetLineWidth(10.0)
	gc.FillStroke()

	saveToPngFile("TestCurveRectangle", i)
}

/*
  <img src="../test_results/TestDrawCubicCurve.png"/>
*/
func TestDrawCubicCurve() {
	i, gc := initGc(w, h)
	// draw a cubic curve
	x, y := 25.6, 128.0
	x1, y1 := 102.4, 230.4
	x2, y2 := 153.6, 25.6
	x3, y3 := 230.4, 128.0

	gc.SetFillColor(color.NRGBA{0xAA, 0xAA, 0xAA, 0xFF})
	gc.SetLineWidth(10)
	gc.MoveTo(x, y)
	gc.CubicCurveTo(x1, y1, x2, y2, x3, y3)
	gc.Stroke()

	gc.SetStrokeColor(color.NRGBA{0xFF, 0x33, 0x33, 0x88})

	gc.SetLineWidth(6)
	// draw segment of curve
	gc.MoveTo(x, y)
	gc.LineTo(x1, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x3, y3)
	gc.Stroke()
	saveToPngFile("TestDrawCubicCurve", i)
}

/*
  <img src="../test_results/TestDash.png"/>
*/
func TestDash() {
	i, gc := initGc(w, h)
	gc.SetLineDash([]float64{50, 10, 10, 10}, -50.0)
	gc.SetLineCap(draw2d.ButtCap)
	gc.SetLineJoin(draw2d.BevelJoin)
	gc.SetLineWidth(10)

	gc.MoveTo(128.0, 25.6)
	gc.LineTo(128.0, 25.6)
	gc.LineTo(230.4, 230.4)
	gc.RLineTo(-102.4, 0.0)
	gc.CubicCurveTo(51.2, 230.4, 51.2, 128.0, 128.0, 128.0)
	gc.Stroke()
	gc.SetLineDash(nil, 0.0)
	saveToPngFile("TestDash", i)
}

/*
  <img src="../test_results/TestFillStroke.png"/>
*/
func TestFillStroke() {
	i, gc := initGc(w, h)
	gc.MoveTo(128.0, 25.6)
	gc.LineTo(230.4, 230.4)
	gc.RLineTo(-102.4, 0.0)
	gc.CubicCurveTo(51.2, 230.4, 51.2, 128.0, 128.0, 128.0)
	gc.Close()

	gc.MoveTo(64.0, 25.6)
	gc.RLineTo(51.2, 51.2)
	gc.RLineTo(-51.2, 51.2)
	gc.RLineTo(-51.2, -51.2)
	gc.Close()

	gc.SetLineWidth(10.0)
	gc.SetFillColor(color.NRGBA{0, 0, 0xFF, 0xFF})
	gc.SetStrokeColor(image.Black)
	gc.FillStroke()
	saveToPngFile("TestFillStroke", i)
}

/*
  <img src="../test_results/TestFillStyle.png"/>
*/
func TestFillStyle() {
	i, gc := initGc(w, h)
	gc.SetLineWidth(6)

	draw2d.Rect(gc, 12, 12, 244, 70)

	wheel1 := new(draw2d.PathStorage)
	wheel1.ArcTo(64, 64, 40, 40, 0, 2*math.Pi)
	wheel2 := new(draw2d.PathStorage)
	wheel2.ArcTo(192, 64, 40, 40, 0, 2*math.Pi)

	gc.SetFillRule(draw2d.FillRuleEvenOdd)
	gc.SetFillColor(color.NRGBA{0, 0xB2, 0, 0xFF})

	gc.SetStrokeColor(image.Black)
	gc.FillStroke(wheel1, wheel2)

	draw2d.Rect(gc, 12, 140, 244, 198)
	wheel1 = new(draw2d.PathStorage)
	wheel1.ArcTo(64, 192, 40, 40, 0, 2*math.Pi)
	wheel2 = new(draw2d.PathStorage)
	wheel2.ArcTo(192, 192, 40, 40, 0, -2*math.Pi)

	gc.SetFillRule(draw2d.FillRuleWinding)
	gc.SetFillColor(color.NRGBA{0, 0, 0xE5, 0xFF})
	gc.FillStroke(wheel1, wheel2)
	saveToPngFile("TestFillStyle", i)
}

func TestMultiSegmentCaps() {
	i, gc := initGc(w, h)
	gc.MoveTo(50.0, 75.0)
	gc.LineTo(200.0, 75.0)

	gc.MoveTo(50.0, 125.0)
	gc.LineTo(200.0, 125.0)

	gc.MoveTo(50.0, 175.0)
	gc.LineTo(200.0, 175.0)

	gc.SetLineWidth(30.0)
	gc.SetLineCap(draw2d.RoundCap)
	gc.Stroke()
	saveToPngFile("TestMultiSegmentCaps", i)
}

func TestRoundRectangle() {
	i, gc := initGc(w, h)
	/* a custom shape that could be wrapped in a function */
	x, y := 25.6, 25.6
	width, height := 204.8, 204.8
	aspect := 1.0                  /* aspect ratio */
	corner_radius := height / 10.0 /* and corner curvature radius */

	radius := corner_radius / aspect
	degrees := math.Pi / 180.0

	gc.ArcTo(x+width-radius, y+radius, radius, radius, -90*degrees, 90*degrees)
	gc.ArcTo(x+width-radius, y+height-radius, radius, radius, 0*degrees, 90*degrees)
	gc.ArcTo(x+radius, y+height-radius, radius, radius, 90*degrees, 90*degrees)
	gc.ArcTo(x+radius, y+radius, radius, radius, 180*degrees, 90*degrees)
	gc.Close()

	gc.SetFillColor(color.NRGBA{0x80, 0x80, 0xFF, 0xFF})
	gc.SetStrokeColor(color.NRGBA{0x80, 0, 0, 0x80})
	gc.SetLineWidth(10.0)
	gc.FillStroke()

	saveToPngFile("TestRoundRectangle", i)
}

func TestLineCap() {
	i, gc := initGc(w, h)
	gc.SetLineWidth(30.0)
	gc.SetLineCap(draw2d.ButtCap)
	gc.MoveTo(64.0, 50.0)
	gc.LineTo(64.0, 200.0)
	gc.Stroke()
	gc.SetLineCap(draw2d.RoundCap)
	gc.MoveTo(128.0, 50.0)
	gc.LineTo(128.0, 200.0)
	gc.Stroke()
	gc.SetLineCap(draw2d.SquareCap)
	gc.MoveTo(192.0, 50.0)
	gc.LineTo(192.0, 200.0)
	gc.Stroke()

	/* draw helping lines */
	gc.SetStrokeColor(color.NRGBA{0xFF, 0x33, 0x33, 0xFF})
	gc.SetLineWidth(2.56)
	gc.MoveTo(64.0, 50.0)
	gc.LineTo(64.0, 200.0)
	gc.MoveTo(128.0, 50.0)
	gc.LineTo(128.0, 200.0)
	gc.MoveTo(192.0, 50.0)
	gc.LineTo(192.0, 200.0)
	gc.Stroke()
	saveToPngFile("TestLineCap", i)
}
func TestLineJoin() {
	i, gc := initGc(w, h)
	gc.SetLineWidth(40.96)
	gc.MoveTo(76.8, 84.48)
	gc.RLineTo(51.2, -51.2)
	gc.RLineTo(51.2, 51.2)
	gc.SetLineJoin(draw2d.MiterJoin) /* default */
	gc.Stroke()

	gc.MoveTo(76.8, 161.28)
	gc.RLineTo(51.2, -51.2)
	gc.RLineTo(51.2, 51.2)
	gc.SetLineJoin(draw2d.BevelJoin)
	gc.Stroke()

	gc.MoveTo(76.8, 238.08)
	gc.RLineTo(51.2, -51.2)
	gc.RLineTo(51.2, 51.2)
	gc.SetLineJoin(draw2d.RoundJoin)
	gc.Stroke()
	saveToPngFile("TestLineJoin", i)
}

func TestBubble() {
	i, gc := initGc(w, h)
	gc.BeginPath()
	gc.MoveTo(75, 25)
	gc.QuadCurveTo(25, 25, 25, 62.5)
	gc.QuadCurveTo(25, 100, 50, 100)
	gc.QuadCurveTo(50, 120, 30, 125)
	gc.QuadCurveTo(60, 120, 65, 100)
	gc.QuadCurveTo(125, 100, 125, 62.5)
	gc.QuadCurveTo(125, 25, 75, 25)
	gc.Stroke()
	saveToPngFile("TestBubble", i)
}

func TestStar() {
	i, gc := initGc(w, h)
	for i := 0.0; i < 360; i = i + 10 { // Go from 0 to 360 degrees in 10 degree steps
		gc.Save()
		gc.SetLineWidth(5) // Keep rotations temporary
		gc.Translate(144, 144)
		gc.Rotate(i * (math.Pi / 180.0)) // Rotate by degrees on stack from 'for'
		gc.MoveTo(0, 0)
		gc.LineTo(72, 0)
		gc.Stroke()
		gc.Restore()
	}
	saveToPngFile("TestStar", i)
}

func TestTransform() {
	i, gc := initGc(800, 600)

	gc.Save()
	gc.Translate(40, 40) // Set origin to (40, 40)
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.RLineTo(72, 0)
	gc.RLineTo(0, 72)
	gc.RLineTo(-72, 0)
	gc.Close()
	gc.Stroke()
	gc.Restore()

	gc.Save()
	gc.Translate(100, 150)            // Translate origin to (100, 150)
	gc.Rotate(30 * (math.Pi / 180.0)) // Rotate counter-clockwise by 30 degrees
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.RLineTo(72, 0)
	gc.RLineTo(0, 72)
	gc.RLineTo(-72, 0)
	gc.Close() // Draw box...
	gc.Stroke()
	gc.Restore()

	gc.Save()
	gc.Translate(40, 300) // Translate to  (40, 300)
	gc.Scale(0.5, 1)      // Reduce x coord by 1/2, y coord left alone
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.RLineTo(72, 0)
	gc.RLineTo(0, 72)
	gc.RLineTo(-72, 0)
	gc.Close() // Draw box...
	gc.Stroke()
	gc.Restore()

	gc.Save()
	gc.Translate(300, 300)            // Set origin to (300, 300)
	gc.Rotate(45 * (math.Pi / 180.0)) // Rotate coordinates by 45 degrees
	gc.Scale(0.5, 1)                  // Scale coordinates
	gc.BeginPath()
	gc.MoveTo(0, 0)
	gc.RLineTo(72, 0)
	gc.RLineTo(0, 72)
	gc.RLineTo(-72, 0)
	gc.Close() // Draw box
	gc.Stroke()
	gc.Restore()

	saveToPngFile("TestTransform", i)
}

func TestPathTransform() {
	i, gc := initGc(800, 600)
	gc.SetLineWidth(20)
	gc.Scale(1, 4)
	gc.ArcTo(200, 80, 50, 50, 0, math.Pi*2)
	gc.Close()
	gc.Stroke()
	saveToPngFile("TestPathTransform", i)
}

func TestFillString() {
	draw2d.SetFontFolder("../resource/font/")
	i, gc := initGc(100, 100)
	draw2d.RoundRect(gc, 5, 5, 95, 95, 10, 10)
	gc.FillStroke()
	gc.SetFillColor(image.Black)
	gc.SetFontSize(18)
	gc.Translate(6, 52)
	gc.SetFontData(draw2d.FontData{"luxi", draw2d.FontFamilyMono, draw2d.FontStyleBold | draw2d.FontStyleItalic})
	width := gc.FillString("cou")
	gc.Translate(width+1, 0)
	left, top, right, bottom := gc.GetStringBounds("cou")
	gc.SetStrokeColor(color.NRGBA{255, 0x33, 0x33, 0x80})
	draw2d.Rect(gc, left, top, right, bottom)
	gc.SetLineWidth(3.0)
	gc.Stroke()
	gc.SetStrokeColor(image.Black)
	gc.SetLineWidth(1.0)
	gc.StrokeString("cou")
	saveToPngFile("TestFillString", i)
}

func TestBigPicture() {
	i, gc := initGc(w, h)
	gc.SetLineWidth(10)

	draw2d.Rect(gc, 0, 0, w, h)
	gc.Fill()
	saveToPngFile("TestBigPicture", i)
}

func main() {
	TestPath()
	TestDrawArc()
	TestDrawArcNegative()
	TestCurveRectangle()
	TestDrawCubicCurve()
	TestDash()
	TestFillStroke()
	TestFillStyle()
	TestMultiSegmentCaps()
	TestRoundRectangle()
	TestLineCap()
	TestLineJoin()
	TestBubble()
	TestStar()
	TestTransform()
	TestPathTransform()
	TestFillString()
	TestBigPicture()
}
