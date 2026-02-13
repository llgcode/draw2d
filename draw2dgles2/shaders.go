// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 11/02/2026 by Copilot

package draw2dgles2

// VertexShader is the vertex shader for rendering primitives
const VertexShader = `
#version 120

attribute vec2 position;
attribute vec4 color;

uniform mat4 projection;

varying vec4 v_color;

void main() {
    gl_Position = projection * vec4(position, 0.0, 1.0);
    v_color = color;
}
`

// FragmentShader is the fragment shader for rendering primitives
const FragmentShader = `
#version 120

varying vec4 v_color;

void main() {
    gl_FragColor = v_color;
}
`

// TextureVertexShader is the vertex shader for textured rendering (text glyphs)
const TextureVertexShader = `
#version 120

attribute vec2 position;
attribute vec2 texCoord;
attribute vec4 color;

uniform mat4 projection;

varying vec2 v_texCoord;
varying vec4 v_color;

void main() {
    gl_Position = projection * vec4(position, 0.0, 1.0);
    v_texCoord = texCoord;
    v_color = color;
}
`

// TextureFragmentShader is the fragment shader for textured rendering
const TextureFragmentShader = `
#version 120

varying vec2 v_texCoord;
varying vec4 v_color;

uniform sampler2D texture;

void main() {
    float alpha = texture2D(texture, v_texCoord).r;
    gl_FragColor = vec4(v_color.rgb, v_color.a * alpha);
}
`
