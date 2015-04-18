package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/llgcode/draw2d"
)

func saveToPngFile(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		return
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Wrote %s OK.\n", filePath)
}

func loadFromPngFile(filePath string) image.Image {
	f, err := os.OpenFile(filePath, 0, 0)
	if f == nil {
		log.Printf("can't open file; err=%s\n", err)
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

func main() {
	source := loadFromPngFile("../resource/image/TestAndroid.png")
	dest := image.NewRGBA(image.Rect(0, 0, 1024, 768))
	width, height := float64(source.Bounds().Dx()), float64(source.Bounds().Dy())
	tr := draw2d.NewIdentityMatrix()
	tr.Translate(width/2, height/2)
	tr.Rotate(30 * math.Pi / 180)
	//tr.Scale(3, 3)
	tr.Translate(-width/2, -height/2)
	tr.Translate(200, 5)
	draw2d.DrawImage(source, dest, tr, draw.Over, draw2d.BilinearFilter)
	saveToPngFile("../resource/result/TestDrawImage.png", dest)
}
