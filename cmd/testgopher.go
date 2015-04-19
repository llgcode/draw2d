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
	width, height = 300, 200
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

func gordon(gc draw2d.GraphicContext, x, y, w, h float64) {
	h23 := (h * 2) / 3

	blf := color.RGBA{0, 0, 0, 0xff}
	wf := color.RGBA{0xff, 0xff, 0xff, 0xff}
	nf := color.RGBA{0x8B, 0x45, 0x13, 0xff}
	brf := color.RGBA{0x8B, 0x45, 0x13, 0x99}
	brb := color.RGBA{0x8B, 0x45, 0x13, 0xBB}

	gc.MoveTo(x, y+h)
	gc.CubicCurveTo(x, y+h, x+w/2, y-h, x+w, y+h)
	gc.Close()
	gc.SetFillColor(brb)
	gc.Fill()
	draw2d.RoundRect(gc, x, y+h, x+w, y+h+h, 10, 10)
	gc.Fill()
	draw2d.Circle(gc, x, y+h, w/12) // left ear
	gc.SetFillColor(brf)
	gc.Fill()
	draw2d.Circle(gc, x, y+h, w/12-10)
	gc.SetFillColor(nf)
	gc.Fill()

	draw2d.Circle(gc, x+w, y+h, w/12) // right ear
	gc.SetFillColor(brf)
	gc.Fill()
	draw2d.Circle(gc, x+w, y+h, w/12-10)
	gc.SetFillColor(nf)
	gc.Fill()

	draw2d.Circle(gc, x+w/3, y+h23, w/9) // left eye
	gc.SetFillColor(wf)
	gc.Fill()
	draw2d.Circle(gc, x+w/3+10, y+h23, w/10-10)
	gc.SetFillColor(blf)
	gc.Fill()
	draw2d.Circle(gc, x+w/3+15, y+h23, 5)
	gc.SetFillColor(wf)
	gc.Fill()

	draw2d.Circle(gc, x+w-w/3, y+h23, w/9) // right eye
	gc.Fill()
	draw2d.Circle(gc, x+w-w/3+10, y+h23, w/10-10)
	gc.SetFillColor(blf)
	gc.Fill()
	draw2d.Circle(gc, x+w-(w/3)+15, y+h23, 5)
	gc.SetFillColor(wf)
	gc.Fill()

	gc.SetFillColor(wf)
	draw2d.RoundRect(gc, x+w/2-w/8, y+h+30, x+w/2-w/8+w/8, y+h+30+w/6, 5, 5) // left tooth
	gc.Fill()
	draw2d.RoundRect(gc, x+w/2, y+h+30, x+w/2+w/8, y+h+30+w/6, 5, 5) // right tooth
	gc.Fill()

	draw2d.Ellipse(gc, x+(w/2), y+h+30, w/6, w/12) // snout
	gc.SetFillColor(nf)
	gc.Fill()
	draw2d.Ellipse(gc, x+(w/2), y+h+10, w/10, w/12) // nose
	gc.SetFillColor(blf)
	gc.Fill()

}

func main() {
	i, gc := initGc(width, height)
	gc.Clear()
	gc.Translate(-75, 58)
	gc.Rotate(-30 * (math.Pi / 180.0))
	gordon(gc, 48, 48, 240, 72)
	saveToPngFile("TestGopher", i)
}
