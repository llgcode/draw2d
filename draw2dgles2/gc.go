// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

package draw2dgles2

import (
	"fmt"
	"image"
	"image/color"
	"log"

	gl "github.com/go-gl/gl/v3.1/gles2"
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
	vao               uint32
	vbo               uint32
	ebo               uint32
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

	// Create VAO (required for OpenGL 3.2+ core profile contexts)
	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)

	// Create VBO for interleaved vertex data
	gl.GenBuffers(1, &r.vbo)

	// Create EBO for index data (required for core profile)
	gl.GenBuffers(1, &r.ebo)

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

	vertexCount := len(r.vertices) / 2

	gl.UseProgram(r.program)

	// Bind VAO (required for core profile contexts)
	gl.BindVertexArray(r.vao)

	// Get attribute locations
	posAttrib := uint32(gl.GetAttribLocation(r.program, gl.Str("position\x00")))
	colorAttrib := uint32(gl.GetAttribLocation(r.program, gl.Str("color\x00")))

	// Interleave position and color data into a single buffer
	vertexSize := 2 + 4 // 2 floats for position, 4 for color
	data := make([]float32, vertexCount*vertexSize)

	for i := 0; i < vertexCount; i++ {
		data[i*vertexSize+0] = r.vertices[i*2+0]
		data[i*vertexSize+1] = r.vertices[i*2+1]
		data[i*vertexSize+2] = r.colors[i*4+0]
		data[i*vertexSize+3] = r.colors[i*4+1]
		data[i*vertexSize+4] = r.colors[i*4+2]
		data[i*vertexSize+5] = r.colors[i*4+3]
	}

	// Upload vertex data to VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STREAM_DRAW)

	// Setup vertex attribute pointers
	stride := int32(vertexSize * 4)
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, stride, gl.PtrOffset(2*4))

	// Upload index data to EBO (required for core profile; client-side indices don't work)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(r.indices)*2, gl.Ptr(r.indices), gl.STREAM_DRAW)

	// Draw triangles using the element buffer
	gl.DrawElements(gl.TRIANGLES, int32(len(r.indices)), gl.UNSIGNED_SHORT, gl.PtrOffset(0))

	gl.DisableVertexAttribArray(posAttrib)
	gl.DisableVertexAttribArray(colorAttrib)
	gl.BindVertexArray(0)

	// Clear batching buffers
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

// AddTriangleStrip renders a triangle strip from matched outer/inner vertex arrays.
// This is used for stroke rendering where the stroke outline forms a strip
// between the outer and inner edges of the path.
func (r *Renderer) AddTriangleStrip(outer, inner []Point2D, clr color.Color) {
	minLen := len(outer)
	if len(inner) < minLen {
		minLen = len(inner)
	}
	if minLen < 2 {
		return
	}

	c := color.RGBAModel.Convert(clr).(color.RGBA)
	red, green, blue, alpha := c.RGBA()
	rf := float32(red) / 65535.0
	gf := float32(green) / 65535.0
	bf := float32(blue) / 65535.0
	af := float32(alpha) / 65535.0

	baseIdx := uint16(len(r.vertices) / 2)

	// Add outer vertices
	for i := 0; i < minLen; i++ {
		r.vertices = append(r.vertices, outer[i].X, outer[i].Y)
		r.colors = append(r.colors, rf, gf, bf, af)
	}
	// Add inner vertices
	for i := 0; i < minLen; i++ {
		r.vertices = append(r.vertices, inner[i].X, inner[i].Y)
		r.colors = append(r.colors, rf, gf, bf, af)
	}

	// Create quads connecting outer[i]-outer[i+1]-inner[i+1]-inner[i]
	for i := 0; i < minLen-1; i++ {
		o0 := baseIdx + uint16(i)
		o1 := baseIdx + uint16(i+1)
		i0 := baseIdx + uint16(minLen+i)
		i1 := baseIdx + uint16(minLen+i+1)

		r.indices = append(r.indices, o0, o1, i0)
		r.indices = append(r.indices, o1, i1, i0)
	}

	// Close the strip (last quad connects back to first)
	o0 := baseIdx + uint16(minLen-1)
	o1 := baseIdx // first outer
	i0 := baseIdx + uint16(2*minLen-1)
	i1 := baseIdx + uint16(minLen) // first inner

	r.indices = append(r.indices, o0, o1, i0)
	r.indices = append(r.indices, o1, i1, i0)
}

// Destroy cleans up OpenGL resources
func (r *Renderer) Destroy() {
	if r.vao != 0 {
		gl.DeleteVertexArrays(1, &r.vao)
	}
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
	}
	if r.ebo != 0 {
		gl.DeleteBuffers(1, &r.ebo)
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

	for _, path := range paths {
		sf := &strokeFlattener{
			renderer:  gc.renderer,
			color:     gc.Current.StrokeColor,
			transform: gc.Current.Tr,
		}

		stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, sf)
		stroker.HalfLineWidth = gc.Current.LineWidth / 2

		var liner draw2dbase.Flattener
		if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
			liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		} else {
			liner = stroker
		}

		draw2dbase.Flatten(path, liner, gc.Current.Tr.GetScale())
		sf.flush()
	}

	gc.Current.Path.Clear()
}

// Fill fills the current path
func (gc *GraphicContext) Fill(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)

	for _, path := range paths {
		for _, polygon := range gc.pathToPolygons(path) {
			gc.renderer.AddPolygon(polygon, gc.Current.FillColor)
		}
	}

	gc.Current.Path.Clear()
}

// FillStroke fills and strokes the current path
func (gc *GraphicContext) FillStroke(paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)

	for _, path := range paths {
		// Collect polygons for filling
		var fillPolygons [][]Point2D
		fillFlattener := &pathFlattener{polygons: &fillPolygons, transform: gc.Current.Tr}

		// Stroke via triangle strip
		sf := &strokeFlattener{
			renderer:  gc.renderer,
			color:     gc.Current.StrokeColor,
			transform: gc.Current.Tr,
		}

		stroker := draw2dbase.NewLineStroker(gc.Current.Cap, gc.Current.Join, sf)
		stroker.HalfLineWidth = gc.Current.LineWidth / 2

		var liner draw2dbase.Flattener
		if gc.Current.Dash != nil && len(gc.Current.Dash) > 0 {
			liner = draw2dbase.NewDashConverter(gc.Current.Dash, gc.Current.DashOffset, stroker)
		} else {
			liner = stroker
		}

		// Use DemuxFlattener to send path to both fill and stroke
		demux := draw2dbase.DemuxFlattener{Flatteners: []draw2dbase.Flattener{fillFlattener, liner}}
		draw2dbase.Flatten(path, demux, gc.Current.Tr.GetScale())
		fillFlattener.flushCurrent()
		sf.flush()

		for _, polygon := range fillPolygons {
			gc.renderer.AddPolygon(polygon, gc.Current.FillColor)
		}
	}

	gc.Current.Path.Clear()
}

// pathToPolygons converts a path to a list of polygons (one per sub-path)
func (gc *GraphicContext) pathToPolygons(path *draw2d.Path) [][]Point2D {
	var polygons [][]Point2D
	flattener := &pathFlattener{polygons: &polygons, transform: gc.Current.Tr}
	draw2dbase.Flatten(path, flattener, gc.Current.Tr.GetScale())
	flattener.flushCurrent()
	return polygons
}

// pathFlattener implements draw2dbase.Flattener to collect vertices
// organized into separate polygons per sub-path.
type pathFlattener struct {
	polygons     *[][]Point2D
	current      []Point2D
	transform    draw2d.Matrix
	lastX, lastY float64
	started      bool
}

// flushCurrent saves the current polygon (if valid) and resets for a new sub-path.
func (pf *pathFlattener) flushCurrent() {
	if len(pf.current) >= 3 {
		*pf.polygons = append(*pf.polygons, pf.current)
	}
	pf.current = nil
	pf.started = false
}

func (pf *pathFlattener) MoveTo(x, y float64) {
	// Flush previous sub-path if any
	pf.flushCurrent()
	x, y = pf.transform.TransformPoint(x, y)
	pf.lastX, pf.lastY = x, y
	pf.started = false
}

func (pf *pathFlattener) LineTo(x, y float64) {
	x, y = pf.transform.TransformPoint(x, y)

	// Add the starting point on the first LineTo after MoveTo
	if !pf.started {
		pf.current = append(pf.current, Point2D{float32(pf.lastX), float32(pf.lastY)})
		pf.started = true
	}

	// Add the current point to form the polygon
	pf.current = append(pf.current, Point2D{float32(x), float32(y)})
	pf.lastX, pf.lastY = x, y
}

func (pf *pathFlattener) LineJoin() {}

func (pf *pathFlattener) Close() {
	pf.flushCurrent()
}

func (pf *pathFlattener) End() {
	pf.flushCurrent()
}

// strokeFlattener receives stroke outline vertices from the LineStroker
// and renders them as a triangle strip between the outer and inner edges.
// The LineStroker outputs vertices in order: outer edge forward, then
// inner edge reversed, then back to start. This flattener splits them
// at the midpoint and creates a quad strip.
type strokeFlattener struct {
	renderer     *Renderer
	color        color.Color
	transform    draw2d.Matrix
	current      []Point2D
	lastX, lastY float64
	started      bool
}

func (sf *strokeFlattener) MoveTo(x, y float64) {
	sf.flush()
	x, y = sf.transform.TransformPoint(x, y)
	sf.lastX, sf.lastY = x, y
	sf.started = false
}

func (sf *strokeFlattener) LineTo(x, y float64) {
	x, y = sf.transform.TransformPoint(x, y)
	if !sf.started {
		sf.current = append(sf.current, Point2D{float32(sf.lastX), float32(sf.lastY)})
		sf.started = true
	}
	sf.current = append(sf.current, Point2D{float32(x), float32(y)})
	sf.lastX, sf.lastY = x, y
}

func (sf *strokeFlattener) LineJoin() {}
func (sf *strokeFlattener) Close()   { sf.flush() }
func (sf *strokeFlattener) End()     { sf.flush() }

func (sf *strokeFlattener) flush() {
	verts := sf.current
	sf.current = nil
	sf.started = false

	if len(verts) < 6 {
		return
	}

	// Remove trailing vertices that duplicate the first vertex
	for len(verts) > 4 {
		last := len(verts) - 1
		dx := verts[last].X - verts[0].X
		dy := verts[last].Y - verts[0].Y
		if dx*dx+dy*dy < 0.5 {
			verts = verts[:last]
		} else {
			break
		}
	}

	n := len(verts)
	if n < 6 {
		return
	}

	mid := n / 2
	outer := verts[:mid]
	inner := make([]Point2D, n-mid)
	copy(inner, verts[mid:])

	// Reverse inner to match outer's direction
	for i, j := 0, len(inner)-1; i < j; i, j = i+1, j-1 {
		inner[i], inner[j] = inner[j], inner[i]
	}

	sf.renderer.AddTriangleStrip(outer, inner, sf.color)
}

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
