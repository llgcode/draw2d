# draw2dgles2 - OpenGL ES 2.0 Renderer for draw2d

## Overview

`draw2dgles2` is a modern, efficient OpenGL ES 2.0-compatible renderer for the draw2d library. It provides hardware-accelerated vector graphics rendering using shader-based techniques, making it suitable for:

- Modern desktop OpenGL (3.0+)
- OpenGL ES 2.0+ (mobile devices, embedded systems)
- WebGL applications
- Cross-platform GUI applications requiring hardware acceleration

## Why OpenGL ES 2.0?

The original `draw2dgl` backend uses OpenGL 2.1 with the legacy fixed-function pipeline (immediate mode):
- Limited to OpenGL 2.1 contexts
- Uses `gl.EnableClientState` and immediate mode rendering
- Not compatible with modern GPU drivers or mobile devices
- Inefficient: rasterizes everything to horizontal lines then renders them

`draw2dgles2` addresses these limitations by using modern OpenGL features:
- **Shader-based rendering** - Custom GLSL shaders for maximum flexibility
- **Vertex Buffer Objects (VBOs)** - Efficient GPU memory management
- **Triangle-based rendering** - Filled shapes rendered as triangles (not lines)
- **Batching system** - Minimizes draw calls for better performance
- **ES 2.0 compatible** - Works on mobile, web, and desktop

## Architecture

### Rendering Pipeline

```
Path Definition → Flattening → Triangulation → Batching → GPU Rendering
```

1. **Path Definition**: User defines paths using MoveTo, LineTo, CurveTo, etc.
2. **Flattening**: Curves are converted to line segments using adaptive subdivision
3. **Triangulation**: Polygons are converted to triangles using ear-clipping algorithm
4. **Batching**: Triangles are collected in batches to minimize draw calls
5. **GPU Rendering**: Batches are uploaded to GPU via VBOs and rendered with shaders

### Key Components

#### 1. Shader System (`shaders.go`)

Two shader programs:

**Basic Shader** (for filled/stroked shapes):
- Vertex shader: Transforms vertices using projection matrix
- Fragment shader: Applies per-vertex colors

**Texture Shader** (for text/glyphs):
- Vertex shader: Transforms vertices and passes texture coordinates
- Fragment shader: Samples texture and applies color/alpha

#### 2. Triangulation (`triangulate.go`)

Implements the **ear-clipping algorithm** to convert arbitrary polygons into triangles:
- O(n²) worst case, but fast enough for typical GUI polygons
- Handles concave polygons correctly
- Produces minimal triangle count

#### 3. Renderer (`gc.go`)

The `Renderer` struct manages:
- Shader programs and uniform locations
- Vertex Buffer Objects (VBOs)
- Batching buffers (vertices, colors, indices)
- Projection matrix setup

The `GraphicContext` implements `draw2d.GraphicContext`:
- Integrates with draw2dbase for path handling
- Converts paths to triangle batches
- Manages graphics state (colors, transforms, line styles)

### Comparison: draw2dgl vs draw2dgles2

| Feature | draw2dgl (Legacy) | draw2dgles2 (Modern) |
|---------|-------------------|----------------------|
| OpenGL Version | 2.1 (fixed pipeline) | ES 2.0+ (programmable) |
| Rendering Method | Horizontal lines via rasterizer | Triangle-based |
| GPU Memory | Client-side arrays | VBOs |
| Shaders | None (fixed function) | Custom GLSL |
| Mobile Support | ❌ No | ✅ Yes |
| Performance | Low (many draw calls) | High (batching) |
| Flexibility | Limited | High (custom shaders) |

## Usage Example

```go
package main

import (
    "image/color"
    "github.com/llgcode/draw2d/draw2dgles2"
    "github.com/llgcode/draw2d/draw2dkit"
)

func main() {
    // Initialize OpenGL context first (using GLFW, SDL, etc.)
    // ... OpenGL context initialization ...
    
    // Create graphics context
    gc, err := draw2dgles2.NewGraphicContext(800, 600)
    if err != nil {
        panic(err)
    }
    defer gc.Destroy()
    
    // Clear screen
    gc.Clear()
    
    // Draw a filled rectangle
    gc.SetFillColor(color.RGBA{255, 0, 0, 255})
    draw2dkit.Rectangle(gc, 100, 100, 300, 300)
    gc.Fill()
    
    // Draw a stroked circle
    gc.SetStrokeColor(color.RGBA{0, 0, 255, 255})
    gc.SetLineWidth(5)
    draw2dkit.Circle(gc, 400, 400, 100)
    gc.Stroke()
    
    // Flush batched drawing commands
    gc.Flush()
    
    // Swap buffers
    // ... swap buffers ...
}
```

## Implementation Details

### OpenGL State Management

The renderer minimizes OpenGL state changes:
- Uses a single shader program per frame where possible
- Batches primitives with the same shader
- Sets up projection matrix once at initialization

### Coordinate System

Uses screen coordinates with origin at top-left:
- (0, 0) = top-left corner
- (width, height) = bottom-right corner
- Y-axis points downward (standard GUI convention)

The projection matrix transforms these coordinates to OpenGL's normalized device coordinates (-1 to 1).

### Performance Considerations

**Batching**: All drawing operations are batched until `Flush()` is called:
```go
gc.Fill()      // Adds to batch
gc.Stroke()    // Adds to batch
gc.Flush()     // Renders everything
```

**VBO Usage**: Dynamic VBOs with `GL_STREAM_DRAW` for frequent updates:
- Buffers are resized automatically
- Data is uploaded once per frame
- Indexed rendering reduces vertex count

**Triangle Count**: Ear-clipping produces O(n) triangles for n-vertex polygons:
- Simple shapes: minimal triangles
- Complex shapes: more triangles but still efficient
- Curves are adaptively subdivided based on scale

### Text Rendering

Text rendering uses a hybrid approach:
1. Glyphs are rasterized to textures (similar to draw2dimg)
2. Textures are cached in GPU memory
3. Text is rendered as textured quads

For production use, consider implementing:
- Texture atlas for glyph caching
- SDF (Signed Distance Field) for scalable text
- Subpixel antialiasing

## Limitations and Future Work

### Current Limitations

1. **Text Rendering**: Uses rasterized glyphs (not GPU-accelerated)
2. **Image Drawing**: `DrawImage()` not yet implemented
3. **Antialiasing**: Relies on OpenGL's MSAA (no custom antialiasing)
4. **Shader Effects**: No advanced shader effects yet

### Planned Improvements

1. **GPU Text Rendering**:
   - Texture atlas for glyph caching
   - SDF rendering for resolution-independent text
   - Better performance for dynamic text

2. **Advanced Features**:
   - Gradient fills (linear, radial)
   - Pattern fills
   - Image texturing
   - Shadow effects
   - Blur effects

3. **Optimizations**:
   - Persistent VBOs for static geometry
   - Instanced rendering for repeated shapes
   - Frustum culling for large scenes
   - GPU-based curve tessellation

## API Compatibility

`draw2dgles2` implements the `draw2d.GraphicContext` interface, making it a drop-in replacement for other backends:

```go
var gc draw2d.GraphicContext

// Can use any backend:
gc = draw2dimg.NewGraphicContext(img)     // CPU rasterizer
gc = draw2dpdf.NewPdf(...)                 // PDF output
gc, _ = draw2dgles2.NewGraphicContext(...) // GPU accelerated
```

All backends support the same drawing operations:
- Path operations (MoveTo, LineTo, CurveTo, etc.)
- Stroke/Fill/FillStroke
- Text rendering
- Transformations
- Graphics state management

## Requirements

- OpenGL ES 2.0+ or OpenGL 3.0+
- Go 1.20+
- OpenGL binding library (e.g., go-gl)

## Building

```bash
go get github.com/llgcode/draw2d/draw2dgles2
go build your-app.go
```

Note: OpenGL library must be installed on your system.

## License

Same as draw2d (BSD-style license)

## References

- [OpenGL ES 2.0 Specification](https://www.khronos.org/opengles/2_X/)
- [Learn OpenGL](https://learnopengl.com/)
- [GPU Gems 3 - Chapter 25: Rendering Vector Art on the GPU](https://developer.nvidia.com/gpugems/gpugems3/part-iv-image-effects/chapter-25-rendering-vector-art-gpu)
- [Loop-Blinn Algorithm](http://research.microsoft.com/en-us/um/people/cloop/loopblinn05.pdf) (for advanced curve rendering)

## Contributing

Contributions welcome! Areas for improvement:
- GPU-accelerated text rendering
- Advanced shader effects
- Performance optimizations
- Additional platform support
- Test coverage
