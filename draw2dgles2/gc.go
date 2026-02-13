// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

package draw2dgles2

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Renderer handles the OpenGL ES 2 rendering
type Renderer struct {
	width, height     int
	program           uint32
	textureProgram    uint32
	vbo               uint32
	projectionUniform int32
	
	// Batching
	vertices []float32
	colors   []float32
	indices  []uint16
}

// NewRenderer creates a new OpenGL ES 2 renderer
func NewRenderer(width, height int) (*Renderer, error) {
	r := &Renderer{
		width:    width,
		height:   height,
		vertices: make([]float32, 0, 4096),
		colors:   make([]float32, 0, 4096),
		indices:  make([]uint16, 0, 2048),
	}

	// Create shader program
	var err error
	r.program, err = createProgram(VertexShader, FragmentShader)
	if err != nil {
		return nil, fmt.Errorf("failed to create shader program: %w", err)
	}

	r.textureProgram, err = createProgram(TextureVertexShader, TextureFragmentShader)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture shader program: %w", err)
	}

	// Get uniform locations
	r.projectionUniform = gl.GetUniformLocation(r.program, gl.Str("projection\x00"))

	// Create VBO
	gl.GenBuffers(1, &r.vbo)

	// Setup projection matrix
	r.setupProjection()

	return r, nil
}

// setupProjection sets up the orthographic projection matrix
func (r *Renderer) setupProjection() {
	gl.UseProgram(r.program)

	// Orthographic projection matrix for screen coordinates
	// Maps (0,0) to top-left, (width, height) to bottom-right
	matrix := [16]float32{
		2.0 / float32(r.width), 0, 0, 0,
		0, -2.0 / float32(r.height), 0, 0,
		0, 0, -1, 0,
		-1, 1, 0, 1,
	}

	gl.UniformMatrix4fv(r.projectionUniform, 1, false, &matrix[0])

	// Also setup for texture program
	gl.UseProgram(r.textureProgram)
	texProjectionUniform := gl.GetUniformLocation(r.textureProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(texProjectionUniform, 1, false, &matrix[0])

	gl.UseProgram(0)
}

// Flush renders all batched primitives
func (r *Renderer) Flush() {
	if len(r.indices) == 0 {
		return
	}

	gl.UseProgram(r.program)

	// Enable attributes
	posAttrib := uint32(gl.GetAttribLocation(r.program, gl.Str("position\x00")))
	colorAttrib := uint32(gl.GetAttribLocation(r.program, gl.Str("color\x00")))

	gl.EnableVertexAttribArray(posAttrib)
	gl.EnableVertexAttribArray(colorAttrib)

	// Upload vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	
	// Interleave position and color data
	vertexSize := 2 + 4 // 2 floats for position, 4 for color
	data := make([]float32, len(r.vertices)/2*vertexSize)
	
	for i := 0; i < len(r.vertices)/2; i++ {
		data[i*vertexSize+0] = r.vertices[i*2+0]
		data[i*vertexSize+1] = r.vertices[i*2+1]
		data[i*vertexSize+2] = r.colors[i*4+0]
		data[i*vertexSize+3] = r.colors[i*4+1]
		data[i*vertexSize+4] = r.colors[i*4+2]
		data[i*vertexSize+5] = r.colors[i*4+3]
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STREAM_DRAW)

	stride := int32(vertexSize * 4)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, stride, gl.PtrOffset(2*4))

	// Draw triangles
	gl.DrawElements(gl.TRIANGLES, int32(len(r.indices)), gl.UNSIGNED_SHORT, gl.Ptr(r.indices))

	gl.DisableVertexAttribArray(posAttrib)
	gl.DisableVertexAttribArray(colorAttrib)

	// Clear buffers
	r.vertices = r.vertices[:0]
	r.colors = r.colors[:0]
	r.indices = r.indices[:0]
}

// AddTriangle adds a triangle to the batch
func (r *Renderer) AddTriangle(x1, y1, x2, y2, x3, y3 float32, c color.Color) {
	baseIdx := uint16(len(r.vertices) / 2)

	// Add vertices
	r.vertices = append(r.vertices, x1, y1, x2, y2, x3, y3)

	// Add colors
	red, green, blue, alpha := c.RGBA()
	rf := float32(red) / 65535.0
	gf := float32(green) / 65535.0
	bf := float32(blue) / 65535.0
	af := float32(alpha) / 65535.0

	for i := 0; i < 3; i++ {
		r.colors = append(r.colors, rf, gf, bf, af)
	}

	// Add indices
	r.indices = append(r.indices, baseIdx, baseIdx+1, baseIdx+2)
}

// AddPolygon adds a filled polygon (will be triangulated)
func (r *Renderer) AddPolygon(vertices []Point2D, c color.Color) {
	if len(vertices) < 3 {
		return
	}

	// Triangulate the polygon
	triangleIndices := Triangulate(vertices)
	if len(triangleIndices) == 0 {
		return
	}

	baseIdx := uint16(len(r.vertices) / 2)

	// Add all vertices
	for _, v := range vertices {
		r.vertices = append(r.vertices, v.X, v.Y)
	}

	// Add colors for all vertices
	red, green, blue, alpha := c.RGBA()
	rf := float32(red) / 65535.0
	gf := float32(green) / 65535.0
	bf := float32(blue) / 65535.0
	af := float32(alpha) / 65535.0

	for range vertices {
		r.colors = append(r.colors, rf, gf, bf, af)
	}

	// Add indices (offset by base index)
	for _, idx := range triangleIndices {
		r.indices = append(r.indices, baseIdx+idx)
	}
}

// Destroy cleans up OpenGL resources
func (r *Renderer) Destroy() {
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
	}
	if r.program != 0 {
		gl.DeleteProgram(r.program)
	}
	if r.textureProgram != 0 {
		gl.DeleteProgram(r.textureProgram)
	}
}

// createProgram creates a shader program from vertex and fragment shader source
func createProgram(vertexSource, fragmentSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)

	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logStr := make([]byte, logLength+1)
		gl.GetProgramInfoLog(program, logLength, nil, &logStr[0])

		return 0, fmt.Errorf("failed to link program: %s", logStr)
	}

	return program, nil
}

// compileShader compiles a shader from source
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logStr := make([]byte, logLength+1)
		gl.GetShaderInfoLog(shader, logLength, nil, &logStr[0])

		return 0, fmt.Errorf("failed to compile shader: %s", logStr)
	}

	return shader, nil
}

// GraphicContext implements the draw2d.GraphicContext interface using OpenGL ES 2
type GraphicContext struct {
	*draw2dbase.StackGraphicContext
	renderer   *Renderer
	FontCache  draw2d.FontCache
	glyphCache draw2dbase.GlyphCache
	glyphBuf   *truetype.GlyphBuf
	DPI        int
}

// NewGraphicContext creates a new OpenGL ES 2 GraphicContext
func NewGraphicContext(width, height int) (*GraphicContext, error) {
	renderer, err := NewRenderer(width, height)
	if err != nil {
		return nil, err
	}

	gc := &GraphicContext{
		StackGraphicContext: draw2dbase.NewStackGraphicContext(),
		renderer:            renderer,
		FontCache:           draw2d.GetGlobalFontCache(),
		glyphCache:          draw2dbase.NewGlyphCache(),
		glyphBuf:            &truetype.GlyphBuf{},
		DPI:                 92,
	}

	return gc, nil
}

// Clear clears the screen
func (gc *GraphicContext) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// ClearRect clears a rectangular region
func (gc *GraphicContext) ClearRect(x1, y1, x2, y2 int) {
	gl.Enable(gl.SCISSOR_TEST)
	gl.Scissor(int32(x1), int32(y1), int32(x2-x1), int32(y2-y1))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Disable(gl.SCISSOR_TEST)
}

// DrawImage draws an image (not yet implemented for ES2 backend)
func (gc *GraphicContext) DrawImage(img image.Image) {
	log.Println("DrawImage not yet implemented for draw2dgles2")
}

// Stroke strokes the current path
func (gc *GraphicContext) Stroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)

	// Convert path to line segments with stroking
	var vertices []Point2D
	for _, path := range paths {
		// Flatten the path to line segments
		flattener := &pathFlattener{vertices: &vertices, transform: gc.Current.Tr}
		stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, flattener)
		stroker.HalfLineWidth = gc.Current.LineWidth / 2

		var liner draw2dbase.Flattener
		if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
			liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		} else {
			liner = stroker
		}

		draw2dbase.Flatten(path, liner, gc.Current.Tr.GetScale())
	}

	if len(vertices) > 0 {
		gc.renderer.AddPolygon(vertices, gc.Current.StrokeColor)
	}

	gc.Current.Path.Clear()
}

// Fill fills the current path
func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)

	// Convert paths to polygons
	for _, path := range paths {
		vertices := gc.pathToVertices(path)
		if len(vertices) > 0 {
			gc.renderer.AddPolygon(vertices, gc.Current.FillColor)
		}
	}

	gc.Current.Path.Clear()
}

// FillStroke fills and strokes the current path
func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {
	// Fill first, then stroke
	gc.Fill(paths...)
	// Re-add the paths since Fill cleared them
	gc.Stroke(paths...)
}

// pathToVertices converts a path to a list of vertices
func (gc *GraphicContext) pathToVertices(path *draw2d.Path) []Point2D {
	var vertices []Point2D
	flattener := &pathFlattener{vertices: &vertices, transform: gc.Current.Tr}
	draw2dbase.Flatten(path, flattener, gc.Current.Tr.GetScale())
	return vertices
}

// pathFlattener implements draw2dbase.Flattener to collect vertices
type pathFlattener struct {
	vertices  *[]Point2D
	transform draw2d.Matrix
	lastX, lastY float64
}

func (pf *pathFlattener) MoveTo(x, y float64) {
	x, y = pf.transform.Transform(x, y)
	pf.lastX, pf.lastY = x, y
}

func (pf *pathFlattener) LineTo(x, y float64) {
	x, y = pf.transform.Transform(x, y)
	*pf.vertices = append(*pf.vertices, Point2D{float32(pf.lastX), float32(pf.lastY)})
	*pf.vertices = append(*pf.vertices, Point2D{float32(x), float32(y)})
	pf.lastX, pf.lastY = x, y
}

func (pf *pathFlattener) LineJoin() {}
func (pf *pathFlattener) Close()    {}
func (pf *pathFlattener) End()      {}

// Flush renders all batched primitives
func (gc *GraphicContext) Flush() {
	gc.renderer.Flush()
}

// Destroy cleans up resources
func (gc *GraphicContext) Destroy() {
	gc.renderer.Destroy()
}

// Font-related methods (simplified for now)

func (gc *GraphicContext) loadCurrentFont() (*truetype.Font, error) {
	font, err := gc.FontCache.Load(gc.Current.FontData)
	if err != nil {
		font, err = gc.FontCache.Load(draw2dbase.DefaultFontData)
	}
	if font != nil {
		gc.SetFont(font)
		gc.SetFontSize(gc.Current.FontSize)
	}
	return font, err
}

func (gc *GraphicContext) SetFont(font *truetype.Font) {
	gc.Current.Font = font
}

func (gc *GraphicContext) SetFontSize(fontSize float64) {
	gc.Current.FontSize = fontSize
	gc.recalc()
}

func (gc *GraphicContext) SetDPI(dpi int) {
	gc.DPI = dpi
	gc.recalc()
}

func (gc *GraphicContext) GetDPI() int {
	return gc.DPI
}

func (gc *GraphicContext) recalc() {
	gc.Current.Scale = gc.Current.FontSize * float64(gc.DPI) * (64.0 / 72.0)
}

// FillString draws filled text
func (gc *GraphicContext) FillString(text string) float64 {
	return gc.FillStringAt(text, 0, 0)
}

// FillStringAt draws filled text at a specific position
func (gc *GraphicContext) FillStringAt(text string, x, y float64) float64 {
	_, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0.0
	}

	// For now, use rasterized glyphs similar to draw2dgl
	// A full implementation would use texture atlases
	startx := x
	prev, hasPrev := truetype.Index(0), false
	fontName := gc.GetFontName()
	
	f := gc.Current.Font
	for _, r := range text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		glyph := gc.glyphCache.Fetch(gc, fontName, r)
		
		// Use draw2dimg's glyph renderer temporarily
		// In a full implementation, this would render to texture atlas
		x += glyph.Fill(gc, x, y)
		
		prev, hasPrev = index, true
	}
	return x - startx
}

// StrokeString draws stroked text
func (gc *GraphicContext) StrokeString(text string) float64 {
	return gc.StrokeStringAt(text, 0, 0)
}

// StrokeStringAt draws stroked text at a specific position
func (gc *GraphicContext) StrokeStringAt(text string, x, y float64) float64 {
	_, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0.0
	}

	startx := x
	prev, hasPrev := truetype.Index(0), false
	fontName := gc.GetFontName()
	
	f := gc.Current.Font
	for _, r := range text {
		index := f.Index(r)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		glyph := gc.glyphCache.Fetch(gc, fontName, r)
		x += glyph.Stroke(gc, x, y)
		prev, hasPrev = index, true
	}
	return x - startx
}

// GetStringBounds returns string bounding box
func (gc *GraphicContext) GetStringBounds(s string) (left, top, right, bottom float64) {
	f, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0
	}
	
	top, left, bottom, right = 10e6, 10e6, -10e6, -10e6
	cursor := 0.0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			cursor += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		if err := gc.glyphBuf.Load(gc.Current.Font, fixed.Int26_6(gc.Current.Scale), index, font.HintingNone); err != nil {
			log.Println(err)
			return 0, 0, 0, 0
		}
		e0 := 0
		for _, e1 := range gc.glyphBuf.Ends {
			ps := gc.glyphBuf.Points[e0:e1]
			for _, p := range ps {
				x, y := pointToF64Point(p)
				top = min(top, y)
				bottom = max(bottom, y)
				left = min(left, x+cursor)
				right = max(right, x+cursor)
			}
			e0 = e1
		}
		cursor += fUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.Current.Scale), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return left, top, right, bottom
}

// CreateStringPath creates a path from string
func (gc *GraphicContext) CreateStringPath(s string, x, y float64) float64 {
	f, err := gc.loadCurrentFont()
	if err != nil {
		log.Println(err)
		return 0.0
	}
	startx := x
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			x += fUnitsToFloat64(f.Kern(fixed.Int26_6(gc.Current.Scale), prev, index))
		}
		err := gc.drawGlyph(index, x, y)
		if err != nil {
			log.Println(err)
			return startx - x
		}
		x += fUnitsToFloat64(f.HMetric(fixed.Int26_6(gc.Current.Scale), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return x - startx
}

func (gc *GraphicContext) drawGlyph(glyph truetype.Index, dx, dy float64) error {
	if err := gc.glyphBuf.Load(gc.Current.Font, fixed.Int26_6(gc.Current.Scale), glyph, font.HintingNone); err != nil {
		return err
	}
	e0 := 0
	for _, e1 := range gc.glyphBuf.Ends {
		drawContour(gc, gc.glyphBuf.Points[e0:e1], dx, dy)
		e0 = e1
	}
	return nil
}

func pointToF64Point(p truetype.Point) (x, y float64) {
	return fUnitsToFloat64(p.X), -fUnitsToFloat64(p.Y)
}

func drawContour(path draw2d.PathBuilder, ps []truetype.Point, dx, dy float64) {
	if len(ps) == 0 {
		return
	}
	startX, startY := pointToF64Point(ps[0])
	var others []truetype.Point
	if ps[0].Flags&0x01 != 0 {
		others = ps[1:]
	} else {
		lastX, lastY := pointToF64Point(ps[len(ps)-1])
		if ps[len(ps)-1].Flags&0x01 != 0 {
			startX, startY = lastX, lastY
			others = ps[:len(ps)-1]
		} else {
			startX = (startX + lastX) / 2
			startY = (startY + lastY) / 2
			others = ps
		}
	}
	path.MoveTo(startX+dx, startY+dy)
	q0X, q0Y, on0 := startX, startY, true
	for _, p := range others {
		qX, qY := pointToF64Point(p)
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				path.LineTo(qX+dx, qY+dy)
			} else {
				path.QuadCurveTo(q0X+dx, q0Y+dy, qX+dx, qY+dy)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				midX := (q0X + qX) / 2
				midY := (q0Y + qY) / 2
				path.QuadCurveTo(q0X+dx, q0Y+dy, midX+dx, midY+dy)
			}
		}
		q0X, q0Y, on0 = qX, qY, on
	}
	// Close the curve.
	if on0 {
		path.LineTo(startX+dx, startY+dy)
	} else {
		path.QuadCurveTo(q0X+dx, q0Y+dy, startX+dx, startY+dy)
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func fUnitsToFloat64(x fixed.Int26_6) float64 {
	scaled := x << 2
	return float64(scaled/256) + float64(scaled%256)/256.0
}

// Make sure the interface is satisfied at compile time
var _ draw2d.GraphicContext = (*GraphicContext)(nil)
