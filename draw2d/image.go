// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
)

type Painter interface {
	raster.Painter
	SetColor(color color.Color)
}

var (
	defaultFontData = FontData{"luxi", FontFamilySans, FontStyleNormal}
)

type ImageGraphicContext struct {
	*StackGraphicContext
	img              draw.Image
	painter          Painter
	fillRasterizer   *raster.Rasterizer
	strokeRasterizer *raster.Rasterizer
	glyphBuf         *truetype.GlyphBuf
	DPI              int
}

/**
 * Create a new Graphic context from an image
 */
func NewGraphicContext(img draw.Image) *ImageGraphicContext {
	var painter Painter
	switch selectImage := img.(type) {
	case *image.RGBA:
		painter = raster.NewRGBAPainter(selectImage)
	default:
		panic("Image type not supported")
	}
	return NewGraphicContextWithPainter(img, painter)
}

// Create a new Graphic context from an image and a Painter (see Freetype-go)
func NewGraphicContextWithPainter(img draw.Image, painter Painter) *ImageGraphicContext {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	dpi := 92
	gc := &ImageGraphicContext{
		NewStackGraphicContext(),
		img,
		painter,
		raster.NewRasterizer(width, height),
		raster.NewRasterizer(width, height),
		truetype.NewGlyphBuf(),
		dpi,
	}
	return gc
}

func (gc *ImageGraphicContext) GetDPI() int {
	return gc.DPI
}

func (gc *ImageGraphicContext) Clear() {
	width, height := gc.img.Bounds().Dx(), gc.img.Bounds().Dy()
	gc.ClearRect(0, 0, width, height)
}

func (gc *ImageGraphicContext) ClearRect(x1, y1, x2, y2 int) {
	imageColor := image.NewUniform(gc.Current.FillColor)
	draw.Draw(gc.img, image.Rect(x1, y1, x2, y2), imageColor, image.ZP, draw.Over)
}

func (gc *ImageGraphicContext) DrawImage(img image.Image) {
	DrawImage(img, gc.img, gc.Current.Tr, draw.Over, BilinearFilter)
}

func (gc *ImageGraphicContext) FillString(text string) (cursor float64) {
	return gc.FillStringAt(text, 0, 0)
}

func (gc *ImageGraphicContext) FillStringAt(text string, x, y float64) (cursor float64) {
	width := gc.CreateStringPath(text, x, y)
	gc.Fill()
	return width
}

func (gc *ImageGraphicContext) StrokeString(text string) (cursor float64) {
	return gc.StrokeStringAt(text, 0, 0)
}

func (gc *ImageGraphicContext) StrokeStringAt(text string, x, y float64) (cursor float64) {
	width := gc.CreateStringPath(text, x, y)
	gc.Stroke()
	return width
}

// CreateStringPath creates a path from the string s at x, y, and returns the string width.
// The text is placed so that the left edge of the em square of the first character of s
// and the baseline intersect at x, y. The majority of the affected pixels will be
// above and to the right of the point, but some may be below or to the left.
// For example, drawing a string that starts with a 'J' in an italic font may
// affect pixels below and left of the point.
func (gc *ImageGraphicContext) CreateStringPath(text string, x, y float64) (width float64) {
	font := GetFont(gc.Current.FontData)
	if font == nil {
		font = GetFont(defaultFontData)
	}
	if font == nil {
		return 0
	}
	gc.SetFont(font)
	gc.SetFontSize(gc.Current.FontSize)
	width, err := gc._createStringPath(text, 0, 0)
	if err != nil {
		log.Println(err)
	}
	return width
}

func fUnitsToFloat64(x int32) float64 {
	scaled := x << 2
	return float64(scaled/256) + float64(scaled%256)/256.0
}

// p is a truetype.Point measured in FUnits and positive Y going upwards.
// The returned value is the same thing measured in floating point and positive Y
// going downwards.
func pointToF64Point(p truetype.Point) (x, y float64) {
	return fUnitsToFloat64(p.X), -fUnitsToFloat64(p.Y)
}

// drawContour draws the given closed contour at the given sub-pixel offset.
func (gc *ImageGraphicContext) drawContour(ps []truetype.Point, dx, dy float64) {
	if len(ps) == 0 {
		return
	}
	startX, startY := pointToF64Point(ps[0])
	gc.MoveTo(startX+dx, startY+dy)
	q0X, q0Y, on0 := startX, startY, true
	for _, p := range ps[1:] {
		qX, qY := pointToF64Point(p)
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				gc.LineTo(qX+dx, qY+dy)
			} else {
				gc.QuadCurveTo(q0X+dx, q0Y+dy, qX+dx, qY+dy)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				midX := (q0X + qX) / 2
				midY := (q0Y + qY) / 2
				gc.QuadCurveTo(q0X+dx, q0Y+dy, midX+dx, midY+dy)
			}
		}
		q0X, q0Y, on0 = qX, qY, on
	}
	// Close the curve.
	if on0 {
		gc.LineTo(startX+dx, startY+dy)
	} else {
		gc.QuadCurveTo(q0X+dx, q0Y+dy, startX+dx, startY+dy)
	}
}

func (gc *ImageGraphicContext) drawGlyph(glyph truetype.Index, dx, dy float64) error {
	if err := gc.glyphBuf.Load(gc.Current.font, gc.Current.scale, glyph, nil); err != nil {
		return err
	}
	e0 := 0
	for _, e1 := range gc.glyphBuf.End {
		gc.drawContour(gc.glyphBuf.Point[e0:e1], dx, dy)
		e0 = e1
	}
	return nil
}

func (gc *ImageGraphicContext) _createStringPath(s string, x, y float64) (float64, error) {
	font := gc.Current.font
	if font == nil {
		return 0.0, errors.New("draw2d: CreateStringPath called with a nil font")
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := font.Index(rune)
		if hasPrev {
			x += fUnitsToFloat64(font.Kerning(gc.Current.scale, prev, index))
		}
		err := gc.drawGlyph(index, x, y)
		if err != nil {
			return startx - x, err
		}
		x += fUnitsToFloat64(font.HMetric(gc.Current.scale, index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x - startx, nil
}

// recalc recalculates scale and bounds values from the font size, screen
// resolution and font metrics, and invalidates the glyph cache.
func (gc *ImageGraphicContext) recalc() {
	gc.Current.scale = int32(gc.Current.FontSize * float64(gc.DPI) * (64.0 / 72.0))
}

// SetDPI sets the screen resolution in dots per inch.
func (gc *ImageGraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.recalc()
}

// SetFont sets the font used to draw text.
func (gc *ImageGraphicContext) SetFont(font *truetype.Font) {
	gc.Current.font = font
}

// SetFontSize sets the font size in points (as in ``a 12 point font'').
func (gc *ImageGraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
	gc.recalc()
}

func (gc *ImageGraphicContext) paint(rasterizer *raster.Rasterizer, color color.Color) {
	gc.painter.SetColor(color)
	rasterizer.Rasterize(gc.painter)
	rasterizer.Clear()
	gc.Current.Path.Clear()
}

/**** second method ****/
func (gc *ImageGraphicContext) Stroke(paths ...*PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.strokeRasterizer.UseNonZeroWinding = true

	stroker := NewLineStroker(gc.Current.Cap, gc.Current.Join, NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2
	var pathConverter *PathConverter
	if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
		dasher := NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		pathConverter = NewPathConverter(dasher)
	} else {
		pathConverter = NewPathConverter(stroker)
	}
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

/**** second method ****/
func (gc *ImageGraphicContext) Fill(paths ...*PathStorage) {
	paths = append(paths, gc.Current.Path)
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()

	/**** first method ****/
	pathConverter := NewPathConverter(NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.fillRasterizer)))
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
}

/* second method */
func (gc *ImageGraphicContext) FillStroke(paths ...*PathStorage) {
	gc.fillRasterizer.UseNonZeroWinding = gc.Current.FillRule.UseNonZeroWinding()
	gc.strokeRasterizer.UseNonZeroWinding = true

	filler := NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.fillRasterizer))

	stroker := NewLineStroker(gc.Current.Cap, gc.Current.Join, NewVertexMatrixTransform(gc.Current.Tr, NewVertexAdder(gc.strokeRasterizer)))
	stroker.HalfLineWidth = gc.Current.LineWidth / 2

	demux := NewDemuxConverter(filler, stroker)
	paths = append(paths, gc.Current.Path)
	pathConverter := NewPathConverter(demux)
	pathConverter.ApproximationScale = gc.Current.Tr.GetScale()
	pathConverter.Convert(paths...)

	gc.paint(gc.fillRasterizer, gc.Current.FillColor)
	gc.paint(gc.strokeRasterizer, gc.Current.StrokeColor)
}

func (f FillRule) UseNonZeroWinding() bool {
	switch f {
	case FillRuleEvenOdd:
		return false
	case FillRuleWinding:
		return true
	}
	return false
}

func (c Cap) Convert() raster.Capper {
	switch c {
	case RoundCap:
		return raster.RoundCapper
	case ButtCap:
		return raster.ButtCapper
	case SquareCap:
		return raster.SquareCapper
	}
	return raster.RoundCapper
}

func (j Join) Convert() raster.Joiner {
	switch j {
	case RoundJoin:
		return raster.RoundJoiner
	case BevelJoin:
		return raster.BevelJoiner
	}
	return raster.RoundJoiner
}
