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

	"github.com/llgcode/draw2d/draw2d"
)

func saveToPngFile(filePath string, m image.Image) {
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
func main() {
	file, err := os.Open("android.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	a, _, err := image.Decode(file)

	//load go icon image
	file2, err := os.Open("go.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	g, _, err := image.Decode(file2)

	if err != nil {
		log.Fatal(err)
	}

	ar := a.Bounds()
	w, h, x := ar.Dx(), ar.Dy(), 30.0
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(i, ar, a, ar.Min, draw.Src)

	tr := draw2d.NewRotationMatrix(x * (math.Pi / 180.0))
	draw2d.DrawImage(g, i, tr, draw.Over, draw2d.LinearFilter)
	saveToPngFile("Test2.png", i)
}
