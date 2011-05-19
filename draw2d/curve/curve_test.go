package curve

import (
	"testing"
	"log"
	"fmt"
	"os"
	"bufio"
	"image"
	"image/png"
	"exp/draw"
	"draw2d.googlecode.com/hg/draw2d/raster"
)


var (
	testsFloat64 = []CubicCurveFloat64{
		CubicCurveFloat64{100, 100, 200, 100, 100, 200, 200, 200},
		CubicCurveFloat64{100, 100, 300, 200, 200, 200, 300, 100},
		CubicCurveFloat64{100, 100, 0, 300, 200, 0, 300, 300},
	}
)

func init() {
	f, err := os.Create("_test.html")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	log.Printf("Create html viewer")
	f.Write([]byte("<html><body>"))
	for i := 0; i < len(testsFloat64); i++ {
		f.Write([]byte(fmt.Sprintf("<div><img src='_testRec%d.png'/><img src='_test%d.png'/></div>", i, i)))
	}
	f.Write([]byte("</body></html>"))

}

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

func drawPoints(img draw.Image, c image.Color, s ...float64) image.Image {
	for i := 0; i < len(s); i += 2 {
		x, y := int(s[i]+0.5), int(s[i+1]+0.5)
		img.Set(x, y, c)
		img.Set(x, y+1, c)
		img.Set(x, y-1, c)
		img.Set(x+1, y, c)
		img.Set(x+1, y+1, c)
		img.Set(x+1, y-1, c)
		img.Set(x-1, y, c)
		img.Set(x-1, y+1, c)
		img.Set(x-1, y-1, c)
		
	}
	return img
}

func TestCubicCurveRec(t *testing.T) {
	for i, curve := range testsFloat64 {
		d := curve.EstimateDistance()
		log.Printf("Distance estimation: %f\n", d)
		numSegments := int(d * 0.25)
		log.Printf("Max segments estimation: %d\n", numSegments)
		s := make([]float64, 0, numSegments)
		s = curve.SegmentRec(s)
		img := image.NewNRGBA(300, 300)
		raster.PolylineBresenham(img, image.NRGBAColor{0xff, 0, 0, 0xff}, curve.X1, curve.Y1, curve.X2, curve.Y2, curve.X3, curve.Y3, curve.X4, curve.Y4)
		raster.PolylineBresenham(img, image.Black, s...)
		drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve.X1, curve.Y1, curve.X2, curve.Y2, curve.X3, curve.Y3, curve.X4, curve.Y4)
		drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, s...)
		savepng(fmt.Sprintf("_testRec%d.png", i), img)
		log.Printf("Num of points: %d\n", len(s))
	}
}

func TestCubicCurve(t *testing.T) {
	for i, curve := range testsFloat64 {
		d := curve.EstimateDistance()
		log.Printf("Distance estimation: %f\n", d)
		numSegments := int(d * 0.25)
		log.Printf("Max segments estimation: %d\n", numSegments)
		s := make([]float64, 0, numSegments)
		s = curve.Segment(s)
		img := image.NewNRGBA(300, 300)
		raster.PolylineBresenham(img, image.NRGBAColor{0xff, 0, 0, 0xff}, curve.X1, curve.Y1, curve.X2, curve.Y2, curve.X3, curve.Y3, curve.X4, curve.Y4)
		raster.PolylineBresenham(img, image.Black, s...)
		drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve.X1, curve.Y1, curve.X2, curve.Y2, curve.X3, curve.Y3, curve.X4, curve.Y4)
		drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, s...)
		savepng(fmt.Sprintf("_test%d.png", i), img)
		log.Printf("Num of points: %d\n", len(s))
	}
}


func BenchmarkCubicCurveRec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, curve := range testsFloat64 {
			d := curve.EstimateDistance()
			numSegments := int(d * 0.25)
			s := make([]float64, 0, numSegments)
			curve.SegmentRec(s)
		}
	}
}

func BenchmarkCubicCurve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, curve := range testsFloat64 {
			d := curve.EstimateDistance()
			numSegments := int(d * 0.25)
			s := make([]float64, 0, numSegments)
			curve.Segment(s)
		}
	}
}
