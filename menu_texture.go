// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glmenu

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var menuTextureVertexShaderSource = `
#version 330

uniform mat4 orthographic_matrix;
uniform vec2 final_position;

in vec4 centered_position;
in vec2 uv;

out vec2 fragment_uv;

void main() {
  fragment_uv = uv;
  vec4 placed = orthographic_matrix * centered_position;
  gl_Position = vec4(placed.x + final_position.x, placed.y + final_position.y, placed.z, placed.w);
}
` + "\x00"

var menuTextureFragmentShaderSource = `
#version 330

uniform sampler2D fragment_texture;
uniform float text_lowerbound;

in vec2 fragment_uv;
out vec4 fragment_color;

void main() {
  vec4 color     = texture(fragment_texture, fragment_uv);
  fragment_color = color;
}
` + "\x00"

// MenuTexture - not yet complete
type MenuTexture struct {
	Menu *Menu

	textureID uint32 // Holds the glyph texture id.
	program   uint32 // program compiled from shaders

	// attributes
	centeredPositionAttribute uint32 // vertex centered_position required for scaling around the orthographic projections center
	uvAttribute               uint32 // texture position

	// The final screen position post-scaling
	finalPositionUniform int32

	// Position of the shaders fragment texture variable
	fragmentTextureUniform int32

	Dimensions mgl32.Vec2

	orthographicMatrixUniform int32
	OrthographicMatrix        mgl32.Mat4

	textureWidth  float32
	textureHeight float32
	WindowWidth   float32
	WindowHeight  float32

	Image *image.NRGBA
}

// ResizeWindow handles a window resize
func (mt *MenuTexture) ResizeWindow(width float32, height float32) {
	mt.WindowWidth = width
	mt.WindowHeight = height
	mt.OrthographicMatrix = mgl32.Ortho2D(-mt.WindowWidth/2, mt.WindowWidth/2, -mt.WindowHeight/2, mt.WindowHeight/2)
}

// Release opengl objects
func (mt *MenuTexture) Release() {
	gl.DeleteTextures(1, &mt.textureID)
}
