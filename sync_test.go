// See also test_test.go

package draw2d_test

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"image"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"
)

func TestSync(t *testing.T) {
	ch := make(chan int)
	limit := 200
	for i := 0; i < limit; i++ {
		go Draw(i, ch)
	}

	for i := 0; i < limit; i++ {
		counter := <-ch
		t.Logf("Goroutine %d returned\n", counter)
	}
}

func Draw(i int, ch chan<- int) {
	draw2d.SetFontCache(testCache)

	// Draw a rounded rectangle using default colors
	dest := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
	gc := draw2dimg.NewGraphicContext(dest)

	draw2dkit.RoundedRectangle(gc, 5, 5, 135, 95, 10, 10)
	gc.FillStroke()

	// Set the fill text color to black
	gc.SetFillColor(image.Black)
	gc.SetFontSize(14)

	// Display Hello World dimensions
	x1, y1, x2, y2 := gc.GetStringBounds("Hello world")
	gc.FillStringAt(fmt.Sprintf("%.2f %.2f %.2f %.2f", x1, y1, x2, y2), 0, 0)

	ch <- i
}

//testFontCache closely follows draw2d's defaultFontCache
type testFontCache struct {
	fonts  map[string]*truetype.Font
	folder string
	namer  draw2d.FontFileNamer
}

func (cache *testFontCache) Load(fontData draw2d.FontData) (font *truetype.Font, err error) {
	if font = cache.fonts[cache.namer(fontData)]; font != nil {
		return font, nil
	}

	var data []byte
	var file = cache.namer(fontData)

	if data, err = ioutil.ReadFile(filepath.Join(cache.folder, file)); err != nil {
		return
	}

	if font, err = truetype.Parse(data); err != nil {
		return
	}

	var mu sync.Mutex
	mu.Lock()
	cache.fonts[file] = font
	mu.Unlock()
	return
}

func (cache *testFontCache) Store(fontData draw2d.FontData, font *truetype.Font) {
	var mu sync.Mutex
	mu.Lock()
	cache.fonts[cache.namer(fontData)] = font
	mu.Unlock()
}

var (
	testFonts = &testFontCache{
		fonts:  make(map[string]*truetype.Font),
		folder: "./resource/font",
		namer:  draw2d.FontFileName,
	}

	testCache draw2d.FontCache = testFonts
)
