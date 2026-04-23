// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

// Package draw2dgles2 provides an efficient graphic context that can draw vector
// graphics and text on OpenGL ES 2.0+ using modern shader-based rendering.
//
// This package provides a more efficient alternative to draw2dgl by using:
//   - Shader-based rendering instead of legacy fixed-function pipeline
//   - Vertex Buffer Objects (VBOs) for better performance
//   - Triangle-based rendering for filled shapes (using ear-clipping triangulation)
//   - Efficient batching to minimize draw calls
//   - Texture atlases for glyph caching
//
// The implementation is compatible with OpenGL ES 2.0, OpenGL 3.0+, and WebGL.
package draw2dgles2
