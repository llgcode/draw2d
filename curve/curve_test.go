package curve

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"testing"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/raster"
)

var (
	flattening_threshold float64 = 0.5
	testsCubicFloat64            = []CubicCurveFloat64{
		CubicCurveFloat64{100, 100, 200, 100, 100, 200, 200, 200},
		CubicCurveFloat64{100, 100, 300, 200, 200, 200, 300, 100},
		CubicCurveFloat64{100, 100, 0, 300, 200, 0, 300, 300},
		CubicCurveFloat64{150, 290, 10, 10, 290, 10, 150, 290},
		CubicCurveFloat64{10, 290, 10, 10, 290, 10, 290, 290},
		CubicCurveFloat64{100, 290, 290, 10, 10, 10, 200, 290},
	}
	testsQuadFloat64 = []QuadCurveFloat64{
		QuadCurveFloat64{100, 100, 200, 100, 200, 200},
		QuadCurveFloat64{100, 100, 290, 200, 290, 100},
		QuadCurveFloat64{100, 100, 0, 290, 200, 290},
		QuadCurveFloat64{150, 290, 10, 10, 290, 290},
		QuadCurveFloat64{10, 290, 10, 10, 290, 290},
		QuadCurveFloat64{100, 290, 290, 10, 120, 290},
	}
)

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

func init() {
	os.Mkdir("test_results", 0666)
	f, err := os.Create("test_results/_test.html")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	log.Printf("Create html viewer")
	f.Write([]byte("<html><body>"))
	for i := 0; i < len(testsCubicFloat64); i++ {
		f.Write([]byte(fmt.Sprintf("<div><img src='_test%d.png'/></div>\n", i)))
	}
	for i := 0; i < len(testsQuadFloat64); i++ {
		f.Write([]byte(fmt.Sprintf("<div><img src='_testQuad%d.png'/>\n</div>\n", i)))
	}
	f.Write([]byte("</body></html>"))

}

func drawPoints(img draw.Image, c color.Color, s ...float64) image.Image {
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

func TestCubicCurve(t *testing.T) {
	for i, curve := range testsCubicFloat64 {
		var p Path
		p.LineTo(curve[0], curve[1])
		curve.Trace(&p, flattening_threshold)
		img := image.NewNRGBA(image.Rect(0, 0, 300, 300))
		raster.PolylineBresenham(img, color.NRGBA{0xff, 0, 0, 0xff}, curve[:]...)
		raster.PolylineBresenham(img, image.Black, p.points...)
		//drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve[:]...)
		drawPoints(img, color.NRGBA{0, 0, 0, 0xff}, p.points...)
		draw2d.SaveToPngFile(fmt.Sprintf("test_results/_test%d.png", i), img)
		log.Printf("Num of points: %d\n", len(p.points))
	}
	fmt.Println()
}

func TestQuadCurve(t *testing.T) {
	for i, curve := range testsQuadFloat64 {
		var p Path
		p.LineTo(curve[0], curve[1])
		curve.Trace(&p, flattening_threshold)
		img := image.NewNRGBA(image.Rect(0, 0, 300, 300))
		raster.PolylineBresenham(img, color.NRGBA{0xff, 0, 0, 0xff}, curve[:]...)
		raster.PolylineBresenham(img, image.Black, p.points...)
		//drawPoints(img, image.NRGBAColor{0, 0, 0, 0xff}, curve[:]...)
		drawPoints(img, color.NRGBA{0, 0, 0, 0xff}, p.points...)
		draw2d.SaveToPngFile(fmt.Sprintf("test_results/_testQuad%d.png", i), img)
		log.Printf("Num of points: %d\n", len(p.points))
	}
	fmt.Println()
}

func BenchmarkCubicCurve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, curve := range testsCubicFloat64 {
			p := Path{make([]float64, 0, 32)}
			p.LineTo(curve[0], curve[1])
			curve.Trace(&p, flattening_threshold)
		}
	}
}
