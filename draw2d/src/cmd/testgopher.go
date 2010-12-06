package main


import (
	"fmt"
	"log"
	"os"
	"bufio"
	"time"
	"math"

	"image"
	"image/png"
	"draw2d"
	//"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
)

const (
	width, height = 500, 300
)

var (
	lastTime int64
	folder   = "../../../../wiki/test_results/"
)

func initGc(w, h int) (image.Image, *draw2d.GraphicContext) {
	i := image.NewRGBA(w, h)
	gc := draw2d.NewGraphicContext(i)
	lastTime = time.Nanoseconds()

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background 
	//gc.Clear()

	return i, gc
}

func saveToPngFile(TestName string, m image.Image) {
	dt := time.Nanoseconds() - lastTime
	fmt.Printf("%s during: %f ms\n", TestName, float(dt)*10e-6)
	filePath := folder + TestName + ".png"
	f, err := os.Open(filePath, os.O_CREAT|os.O_WRONLY, 0600)
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

func gordon(gc *draw2d.GraphicContext, x, y, w, h float) {
	h23 := (h * 2) / 3

	blf := image.RGBAColor{0, 0, 0,  0xff}
	wf := image.RGBAColor{0xff, 0xff, 0xff, 0xff}
	nf := image.RGBAColor{0x8B, 0x45, 0x13, 0xff}
	brf := image.RGBAColor{0x8B, 0x45, 0x13, 0x99}
	brb := image.RGBAColor{0x8B, 0x45, 0x13, 0xBB}
	
	gc.MoveTo(x, y+h)
	gc.CubicCurveTo(x, y+h, x+w/2, y-h, x+w, y+h)
	gc.ClosePath()
	gc.SetFillColor(brb)
	gc.Fill()
	gc.RoundRect(x, y+h, x+ w, y+h+h, 10, 10)
	gc.Fill()
	gc.Circle(x, y+h, w/12) // left ear
	gc.SetFillColor(brf)
	gc.Fill()
	gc.Circle(x, y+h, w/12-10)
	gc.SetFillColor(nf)
	gc.Fill()
	
	gc.Circle(x+w, y+h, w/12) // right ear
	gc.SetFillColor(brf)
	gc.Fill()
	gc.Circle(x+w, y+h, w/12-10)
	gc.SetFillColor(nf)
	gc.Fill()

	gc.Circle(x+w/3, y+h23, w/9) // left eye
	gc.SetFillColor(wf)
	gc.Fill()
	gc.Circle(x+w/3+10, y+h23, w / 10 - 10)
	gc.SetFillColor(blf)
	gc.Fill()
	gc.Circle(x+w/3+15, y+h23, 5)
	gc.SetFillColor(wf)
	gc.Fill()

	gc.Circle(x+w-w/3, y+h23, w/9) // right eye
	gc.Fill()
	gc.Circle(x+w-w/3+10, y+h23, w / 10 - 10)
	gc.SetFillColor(blf)
	gc.Fill()
	gc.Circle(x+w-(w/3)+15, y+h23, 5)
	gc.SetFillColor(wf)
	gc.Fill()

	gc.SetFillColor(wf)
	gc.RoundRect(x+w/2-w/8, y+h+30, x+w/2-w/8 + w/8, y+h+30 + w/6, 5, 5) // left tooth
	gc.Fill()
	gc.RoundRect(x+w/2, y+h+30, x+w/2+w/8, y+h+30+w/6, 5, 5)    // right tooth
	gc.Fill()


	gc.Ellipse(x+(w/2), y+h+30, w/6, w/12)   // snout
	gc.SetFillColor(nf)
	gc.Fill()
	gc.Ellipse(x+(w/2), y+h+10, w/10, w/12) // nose
	gc.SetFillColor(blf)
	gc.Fill()
	
}

func main() {
	i, gc := initGc(width, height)
	gc.Clear()
	gc.Translate(100, 100)
	gc.Rotate(-30 * (math.Pi / 180.0))
	gordon(gc, 48, 48, 240, 72)
	saveToPngFile("TestGopher", i)
}