package main

import (
	"bufio"
	"github.com/llgcode/draw2d/draw2d"
	"github.com/llgcode/draw2d/postscript"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	i := image.NewRGBA(image.Rect(0, 0, 600, 800))
	gc := draw2d.NewGraphicContext(i)
	gc.Translate(0, 380)
	gc.Scale(1, -1)
	gc.Translate(0, -380)
	src, err := os.OpenFile("../resource/postscript/tiger.ps", 0, 0)
	if err != nil {
		return
	}
	defer src.Close()
	bytes, err := ioutil.ReadAll(src)
	reader := strings.NewReader(string(bytes))
	interpreter := postscript.NewInterpreter(gc)
	interpreter.Execute(reader)
	saveToPngFile("../resource/result/TestPostscript.png", i)
}
