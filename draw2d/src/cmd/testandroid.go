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
	//"draw2d"
	"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
)

const (
	width, height = 178, 224
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

func android(gc *draw2d.GraphicContext, x, y float) {
	gc.SetLineCap(draw2d.RoundCap)
	gc.SetLineWidth(5)
	gc.ArcTo(x+80, y+70, 50, 50, 180 * (math.Pi/180), 360 * (math.Pi/180))                     // head
	gc.FillStroke()
	gc.MoveTo(x+60, y+25)
	gc.LineTo(x+50, y+10)
	gc.MoveTo(x+100, y+25)
	gc.LineTo( x+110, y+10)
	gc.Stroke()
	draw2d.Circle(gc, x+60, y+45, 5)                                      // left eye
	gc.FillStroke()
	draw2d.Circle(gc, x+100, y+45, 5)                                   // right eye
	gc.FillStroke()
	draw2d.RoundRect(gc, x+30, y+75, x+30+100, y+75+90, 10, 10)                                   // body
	gc.FillStroke()
	draw2d.Rect(gc, x+30, y+75, x+30+100, y+75+80)
	gc.FillStroke()
	draw2d.RoundRect(gc, x+5, y+80, x+5+20, y+80+70, 10, 10)   // left arm
	gc.FillStroke()
	draw2d.RoundRect(gc, x+135, y+80, x+135+20, y+80+70, 10, 10) // right arm
	gc.FillStroke()
	draw2d.RoundRect(gc, x+50, y+150, x+50+20, y+150+50, 10, 10) // left leg
	gc.FillStroke()
	draw2d.RoundRect(gc, x+90, y+150, x+90+20, y+150+50, 10, 10) // right leg
	gc.FillStroke()
}



func main() {
	i, gc := initGc(width, height)
	gc.SetFillColor(image.RGBAColor{0xff, 0x44, 0x44, 0xff})
	gc.SetStrokeColor(image.RGBAColor{0x44, 0x44, 0x44, 0xff})
	android(gc, 10, 10)
	saveToPngFile("TestAndroid", i)
}
