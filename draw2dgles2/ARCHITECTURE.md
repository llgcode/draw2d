# Where are the OpenGL Draw Calls? Understanding draw2dgl Architecture

## The Original Question

> "Where is the code that actually calls OpenGL? All I can find are a few lines in gc.go that render lines. But where are the triangles getting rendered?"

## The Answer

The original `draw2dgl` backend **does not render triangles**. Instead, it uses an unusual approach:

### How draw2dgl Works

1. **Paths are rasterized to horizontal scanlines** using the freetype rasterizer (`github.com/golang/freetype/raster`)
2. **Scanlines are converted to OpenGL lines** in the `Painter` struct
3. **Lines are rendered using legacy OpenGL** with client-side arrays

#### The Actual OpenGL Calls (draw2dgl/gc.go, lines 82-95)

```go
func (p *Painter) Flush() {
    if len(p.vertices) != 0 {
        // Enable legacy client-side arrays (deprecated in OpenGL 3.0+)
        gl.EnableClientState(gl.COLOR_ARRAY)
        gl.EnableClientState(gl.VERTEX_ARRAY)
        
        // Set up pointers to vertex and color data
        gl.ColorPointer(4, gl.UNSIGNED_BYTE, 0, gl.Ptr(p.colors))
        gl.VertexPointer(2, gl.INT, 0, gl.Ptr(p.vertices))

        // THIS IS THE ACTUAL DRAW CALL
        // Renders horizontal lines from the rasterized spans
        gl.DrawArrays(gl.LINES, 0, int32(len(p.vertices)/2))
        
        gl.DisableClientState(gl.VERTEX_ARRAY)
        gl.DisableClientState(gl.COLOR_ARRAY)
        
        // Clear buffers for next batch
        p.vertices = p.vertices[0:0]
        p.colors = p.colors[0:0]
    }
}
```

#### The Rendering Pipeline

```
Vector Path → Stroke/Fill → Rasterizer → Spans → Lines → OpenGL
```

Detailed flow:

1. **User defines path**: `gc.MoveTo()`, `gc.LineTo()`, `gc.CubicCurveTo()`, etc.
2. **Path is stroked or filled**: `gc.Stroke()` or `gc.Fill()` is called
3. **Path is flattened**: Curves converted to line segments
4. **Rasterizer processes path**: Freetype's rasterizer converts to coverage spans
5. **Painter receives spans**: Each span is a horizontal line with alpha coverage
6. **Spans converted to GL vertices**: `Painter.Paint()` adds vertices to buffer
7. **GL renders the lines**: `Painter.Flush()` calls `gl.DrawArrays(gl.LINES, ...)`

### Why No Triangles?

The original implementation chose to reuse the CPU rasterizer from freetype rather than implementing GPU-based triangulation. This means:

- ❌ **No triangle rendering**
- ❌ **No GPU-accelerated rasterization**
- ✅ CPU does all the heavy lifting
- ✅ OpenGL just displays the pre-rasterized result as lines

This approach has several problems:
- Limited to OpenGL 2.1 (client-side arrays are deprecated)
- Inefficient (rasterizing on CPU, then uploading lines to GPU)
- Not compatible with OpenGL ES or modern contexts
- Many draw calls (one per span)

## The Modern Solution: draw2dgles2

The new `draw2dgles2` package addresses these limitations with proper GPU rendering:

### Modern Rendering Pipeline

```
Vector Path → Flatten → Triangulate → Batch → GPU Shaders
```

1. **Paths are flattened** to line segments
2. **Polygons are triangulated** using ear-clipping algorithm
3. **Triangles are batched** in GPU memory
4. **Custom shaders render** the triangles

### Where Are the Triangles in draw2dgles2?

#### Triangulation (draw2dgles2/triangulate.go)

The `Triangulate()` function converts polygons to triangles:

```go
// Ear-clipping algorithm
func Triangulate(vertices []Point2D) []uint16 {
    // Returns triangle indices
    // Each triplet of indices forms one triangle
}
```

#### Triangle Rendering (draw2dgles2/gc.go)

```go
func (r *Renderer) AddPolygon(vertices []Point2D, c color.Color) {
    // 1. Triangulate the polygon
    triangleIndices := Triangulate(vertices)
    
    // 2. Add vertices to batch
    for _, v := range vertices {
        r.vertices = append(r.vertices, v.X, v.Y)
    }
    
    // 3. Add colors
    for range vertices {
        r.colors = append(r.colors, rf, gf, bf, af)
    }
    
    // 4. Add triangle indices
    for _, idx := range triangleIndices {
        r.indices = append(r.indices, baseIdx+idx)
    }
}

func (r *Renderer) Flush() {
    // Upload to GPU via VBO
    gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STREAM_DRAW)
    
    // THIS IS THE ACTUAL TRIANGLE DRAW CALL
    gl.DrawElements(gl.TRIANGLES, int32(len(r.indices)), gl.UNSIGNED_SHORT, gl.Ptr(r.indices))
}
```

### Comparison Table

| Aspect | draw2dgl (Old) | draw2dgles2 (New) |
|--------|----------------|-------------------|
| **Rasterization** | CPU (freetype) | GPU (triangles) |
| **Primitives** | Horizontal lines | Triangles |
| **OpenGL Calls** | `DrawArrays(LINES)` | `DrawElements(TRIANGLES)` |
| **Memory** | Client arrays | VBOs |
| **Shaders** | None (fixed-function) | Custom GLSL |
| **Batching** | Per-span | All shapes |
| **OpenGL Version** | 2.1 only | ES 2.0 / 3.0+ / WebGL |
| **Performance** | Low | High |

## Code Locations

### draw2dgl (Legacy)

- **Main file**: `draw2dgl/gc.go`
- **OpenGL draw call**: Line 90: `gl.DrawArrays(gl.LINES, ...)`
- **Painter**: Lines 26-120 (converts rasterizer spans to lines)
- **Flush method**: Lines 82-96 (the actual rendering)

### draw2dgles2 (Modern)

- **Main file**: `draw2dgles2/gc.go`
- **OpenGL draw call**: Line ~168: `gl.DrawElements(gl.TRIANGLES, ...)`
- **Triangulation**: `draw2dgles2/triangulate.go`
- **Shaders**: `draw2dgles2/shaders.go`
- **Renderer**: Lines 18-284 (manages GPU resources)

## Key Insights

1. **draw2dgl doesn't use triangles** - it renders rasterized spans as horizontal lines
2. **The "trick" is in the Painter** - it receives coverage spans from the rasterizer and converts them to OpenGL lines
3. **Modern OpenGL requires triangles** - which is why draw2dgles2 was created
4. **Triangulation is necessary** for GPU rendering - the ear-clipping algorithm handles this

## Further Reading

- **Rasterization vs GPU Rendering**: 
  - CPU: compute pixel coverage → upload to GPU
  - GPU: upload geometry → GPU computes coverage
  
- **Why Triangles**:
  - GPUs are optimized for triangle rasterization
  - All modern graphics APIs use triangles as the fundamental primitive
  - Efficient hardware implementation

- **Modern Approaches**:
  - NV_path_rendering (NVIDIA extension for vector graphics)
  - Loop-Blinn algorithm (curve rendering via shaders)
  - Stencil-and-cover (two-pass rendering)

## Conclusion

The original draw2dgl's OpenGL calls are minimal because it offloads rasterization to the CPU. The new draw2dgles2 backend provides true GPU-accelerated rendering with triangle-based primitives and modern shader support, making it suitable for OpenGL ES 2.0 and beyond.
