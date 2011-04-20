package main


import (
	"fmt"
	"time"
	"log"
	"os"
	"io/ioutil"
	"bufio"
	"strings"
	"image"
	"image/png"
	"draw2d.googlecode.com/hg/draw2d"
	"draw2d.googlecode.com/hg/postscript"
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
	i := image.NewRGBA(600, 800)
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
	lastTime := time.Nanoseconds()
	interpreter.Execute(reader)
	dt := time.Nanoseconds() - lastTime
	fmt.Printf("Draw image: %f ms\n", float64(dt)*1e-6)
	saveToPngFile("../resource/result/TestPostscript.png", i)
}
