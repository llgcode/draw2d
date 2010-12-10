package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"math"
	"image"
	"time"
	"image/png"
	"draw2d"
	//"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
)


func saveToPngFile(filePath string, m image.Image) {
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

func loadFromPngFile(filePath string) image.Image {
	f, err := os.Open(filePath, 0, 0) 
	if f == nil {
		log.Printf("can't open file; err=%s\n", err.String())
		return nil
	}
	defer f.Close()
	b := bufio.NewReader(f)
	i, err := png.Decode(b)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Read %s OK.\n", filePath)
	return i
}


func testBubble(gc * draw2d.GraphicContext) {
	gc.BeginPath()
	gc.MoveTo(75, 25)
	gc.QuadCurveTo(25, 25, 25, 62.5)
	gc.QuadCurveTo(25, 100, 50, 100)
	gc.QuadCurveTo(50, 120, 30, 125)
	gc.QuadCurveTo(60, 120, 65, 100)
	gc.QuadCurveTo(125, 100, 125, 62.5)
	gc.QuadCurveTo(125, 25, 75, 25)
	gc.Stroke()
}

func main() {


	source := loadFromPngFile("../../Varna_Railway_Station_HDR.png")
	i := image.NewRGBA(1024, 768)
	gc := draw2d.NewGraphicContext(i)
	gc.Scale(2, 0.5)
	//gc.Translate(75, 25)
	gc.Rotate(30 * math.Pi/180)
	lastTime := time.Nanoseconds()
	gc.DrawImage(source)
	dt := time.Nanoseconds() - lastTime
	fmt.Printf("Draw image: %f ms\n", float(dt)*1e-6)
	saveToPngFile("../../TestDrawImage.png", i)
}
