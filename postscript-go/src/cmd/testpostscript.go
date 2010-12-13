package main


import (
	"fmt"
	"time"
	"log"
	"os"
	"bufio"
	"image"
	"image/png"
	"draw2d.googlecode.com/svn/trunk/draw2d/src/pkg/draw2d"
	"postscript"
	//"postscript-go.googlecode.com/svn/trunk/postscript-go/src/pkg/postscript"
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

func main() {
	i := image.NewRGBA(600, 800)
	gc := draw2d.NewGraphicContext(i)
	gc.Translate(0, 380)
	gc.Scale(1,-1)
	gc.Translate(0, -380)
	lastTime := time.Nanoseconds()
	interpreter := postscript.NewInterpreter(gc)
	interpreter.ExecuteFile("../../test_files/tiger.ps")
	dt := time.Nanoseconds() - lastTime
	fmt.Printf("Draw image: %f ms\n", float(dt)*1e-6)
	saveToPngFile("../../TestPostscript.png", i)
}
