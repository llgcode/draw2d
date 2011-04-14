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
	"exp/draw"
	"draw2d.googlecode.com/hg/draw2d"
)


func saveToPngFile(filePath string, m image.Image) {
	f, err := os.Open(filePath, os.O_CREAT|os.O_WRONLY, 0600)
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


func main() {
	source := loadFromPngFile("../resource/image/TestAndroid.png")
	dest := image.NewRGBA(1024, 768)
	width, height := float64(source.Bounds().Dx()), float64(source.Bounds().Dy())
	tr := draw2d.NewIdentityMatrix()
	tr.Translate(width/2, height/2)
	tr.Rotate(30 * math.Pi / 180)
	tr.Translate(-width/2, -height/2)
	tr.Translate(75, 25)
	lastTime := time.Nanoseconds()
	draw2d.DrawImage(source, dest, tr, draw.Over, draw2d.BilinearFilter)
	dt := time.Nanoseconds() - lastTime
	fmt.Printf("Draw image: %f ms\n", float64(dt)*1e-6)
	saveToPngFile("../resource/result/TestDrawImage.png", dest)
}
