package glmenu

import (
	"github.com/go-gl/mathgl/mgl32"
)

type MenuImage struct {
	MenuTexture      *MenuTexture
	MenuTextureIndex int

	// final position on screen
	finalPosition mgl32.Vec2

	// text color
	color mgl32.Vec3

	// general opengl values
	vao           uint32
	vbo           uint32
	ebo           uint32
	vboData       []float32
	vboIndexCount int
	eboData       []int32
	eboIndexCount int

	// X1, X2: the lower left and upper right points of a box that bounds the text with a center point (0,0)
	// lower left
	X1 Point
	// upper right
	X2 Point

	// Screen position away from center
	Position mgl32.Vec2
}

// NewMenuImage creates a new image using the subimage specified by index
/*
func NewMenuImage(mt *MenuTexture, index int) (mi *MenuImage) {
	mi = &MenuImage{}
	mi.MenuTexture = mt
	mi.MenuTextureIndex = index

	glfloat_size := int32(4)

	// stride of the buffered data
	xy_count := int32(2)
	stride := xy_count + int32(2)

	gl.GenVertexArrays(1, &mi.vao)
	gl.GenBuffers(1, &mi.vbo)
	gl.GenBuffers(1, &mi.ebo)

	// vao
	gl.BindVertexArray(mi.vao)

	gl.BindTexture(gl.TEXTURE_2D, mi.MenuTexture.textureID)

	// vbo
	// specify the buffer for which the VertexAttribPointer calls apply
	gl.BindBuffer(gl.ARRAY_BUFFER, mi.vbo)

	gl.EnableVertexAttribArray(mi.MenuTexture.centeredPositionAttribute)
	gl.VertexAttribPointer(
		mi.MenuTexture.centeredPositionAttribute,
		2,
		gl.FLOAT,
		false,
		glfloat_size*stride,
		gl.PtrOffset(0),
	)

	gl.EnableVertexAttribArray(mi.MenuTexture.uvAttribute)
	gl.VertexAttribPointer(
		mi.MenuTexture.uvAttribute,
		2,
		gl.FLOAT,
		false,
		glfloat_size*stride,
		gl.PtrOffset(int(glfloat_size*xy_count)),
	)

	// ebo
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mi.ebo)

	// i am guessing that order is important here
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	mi.vboIndexCount = 4 * 2 * 2 // 4 indexes per image (containing 2 position + 2 texture)
	mi.eboIndexCount = 6         // each image requires 6 triangle indices for a quad
	mi.vboData = make([]float32, mi.vboIndexCount, mi.vboIndexCount)
	mi.eboData = make([]int32, mi.eboIndexCount, mi.eboIndexCount)

	// generate the basic vbo data and bounding box
	// center the vbo data around the orthographic (0,0) point
	mi.X1 = Point{0, 0}
	mi.X2 = Point{0, 0}
	mi.makeBufferData(indices)
	mi.centerTheData(t.getLowerLeft())

	gl.BindVertexArray(mi.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, mi.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER, int(glfloat_size)*mi.vboIndexCount, gl.Ptr(mi.vboData), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mi.ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, int(glfloat_size)*mi.eboIndexCount, gl.Ptr(mi.eboData), gl.DYNAMIC_DRAW)
	gl.BindVertexArray(0)

	// possibly not necesssary?
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return mi
}
*/
