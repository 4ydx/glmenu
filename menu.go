package menu

import (
	gltext "github.com/4ydx/gltext"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glow/gl-core/3.3/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Point struct {
	X, Y float32
}

var vertexShaderSource string = `
#version 330

uniform mat4 matrix;

in vec4 position;

void main() {
  gl_Position = matrix * position;
}
` + "\x00"

var fragmentShaderSource string = `
#version 330

out vec4 fragment_color;

void main() {
  fragment_color = vec4(1,1,1,1);
}
` + "\x00"

type Menu struct {
	// options
	Visible      bool
	ShowOn       glfw.Key
	Height       float32
	Width        float32
	IsAutoCenter bool
	LowerLeft    Point

	// objects
	Text []Text
	//Fields []Field

	// opengl oriented
	windowWidth   float32
	windowHeight  float32
	program       uint32 // shader program
	glMatrix      int32  // ortho matrix
	position      uint32 // index location
	vao           uint32
	vbo           uint32
	ebo           uint32
	ortho         mgl32.Mat4
	vboData       []float32
	vboIndexCount int
	eboData       []int32
	eboIndexCount int
}

func (menu *Menu) Toggle() {
	menu.Visible = !menu.Visible
}

func (menu *Menu) SetDimension(w, h float32) {
	menu.Width = w
	menu.Height = h
}

func (menu *Menu) Load(lowerLeft Point) error {
	var err error
	glfloat_size := 4
	glint_size := 4

	menu.Visible = false
	menu.ShowOn = glfw.KeyM
	menu.LowerLeft = lowerLeft

	// create shader program and define attributes and uniforms
	menu.program, err = gltext.NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}
	menu.glMatrix = gl.GetUniformLocation(menu.program, gl.Str("matrix\x00"))
	menu.position = uint32(gl.GetAttribLocation(menu.program, gl.Str("position\x00")))

	gl.GenVertexArrays(1, &menu.vao)
	gl.GenBuffers(1, &menu.vbo)
	gl.GenBuffers(1, &menu.ebo)

	// vao
	gl.BindVertexArray(menu.vao)

	// vbo
	// specify the buffer for which the VertexAttribPointer calls apply
	gl.BindBuffer(gl.ARRAY_BUFFER, menu.vbo)

	gl.EnableVertexAttribArray(menu.position)
	gl.VertexAttribPointer(
		menu.position,
		2,
		gl.FLOAT,
		false,
		0, // no stride... yet
		gl.PtrOffset(0),
	)

	// ebo
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, menu.ebo)

	// i am guessing that order is important here
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	// ebo, vbo data
	menu.vboIndexCount = 4 * 2 // four indices (2 points per index)
	menu.eboIndexCount = 6     // 6 triangle indices for a quad
	menu.vboData = make([]float32, menu.vboIndexCount, menu.vboIndexCount)
	menu.eboData = make([]int32, menu.eboIndexCount, menu.eboIndexCount)
	menu.makeBufferData()

	// setup context
	gl.BindVertexArray(menu.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, menu.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER, glfloat_size*menu.vboIndexCount, gl.Ptr(menu.vboData), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, menu.ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, glint_size*menu.eboIndexCount, gl.Ptr(menu.eboData), gl.DYNAMIC_DRAW)
	gl.BindVertexArray(0)

	// not necesssary, but i just want to better understand using vertex arrays
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	return nil
}

func (menu *Menu) ResizeWindow(width float32, height float32) {
	menu.windowWidth = width
	menu.windowHeight = height
	for _, text := range menu.Text {
		text.ResizeWindow(width, height)
	}
	menu.ortho = mgl32.Ortho2D(0, menu.windowWidth, 0, menu.windowHeight)
}

func (menu *Menu) makeBufferData() {
	// index (0,0)
	menu.vboData[0] = menu.LowerLeft.X // position
	menu.vboData[1] = menu.LowerLeft.Y

	// index (1,0)
	menu.vboData[2] = menu.LowerLeft.X + menu.Width
	menu.vboData[3] = menu.LowerLeft.Y

	// index (1,1)
	menu.vboData[4] = menu.LowerLeft.X + menu.Width
	menu.vboData[5] = menu.LowerLeft.Y + menu.Height

	// index (0,1)
	menu.vboData[6] = menu.LowerLeft.X
	menu.vboData[7] = menu.LowerLeft.Y + menu.Height

	menu.eboData[0] = 0
	menu.eboData[1] = 1
	menu.eboData[2] = 2
	menu.eboData[3] = 0
	menu.eboData[4] = 2
	menu.eboData[5] = 3
}

func (menu *Menu) Release() {
	gl.DeleteBuffers(1, &menu.vbo)
	gl.DeleteBuffers(1, &menu.ebo)
	gl.DeleteBuffers(1, &menu.vao)
}

func (menu *Menu) Draw() bool {
	if !menu.Visible {
		return menu.Visible
	}
	gl.UseProgram(menu.program)

	gl.UniformMatrix4fv(menu.glMatrix, 1, false, &menu.ortho[0])

	gl.BindVertexArray(menu.vao)
	gl.DrawElements(gl.TRIANGLES, int32(menu.eboIndexCount), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
	for _, text := range menu.Text {
		if !text.IsHover {
			text.SetScale(1)
		}
		text.Draw()
	}
	return menu.Visible
}

func (menu *Menu) ScreenClick(xPos, yPos float64) {
	if !menu.Visible {
		return
	}
	yPos = float64(menu.windowHeight) - yPos
	if xPos > float64(menu.LowerLeft.X) && xPos < float64(menu.LowerLeft.X+menu.Width) && yPos > float64(menu.LowerLeft.Y) && yPos < float64(menu.LowerLeft.Y+menu.Height) {
		for _, text := range menu.Text {
			text.IsClicked(xPos, yPos)
		}
	}
}

func (menu *Menu) ScreenHover(xPos, yPos float64) {
	if !menu.Visible {
		return
	}
	yPos = float64(menu.windowHeight) - yPos
	if xPos > float64(menu.LowerLeft.X) && xPos < float64(menu.LowerLeft.X+menu.Width) && yPos > float64(menu.LowerLeft.Y) && yPos < float64(menu.LowerLeft.Y+menu.Height) {
		for i, text := range menu.Text {
			text.IsHovered(xPos, yPos)
			menu.Text[i] = text
		}
	}
}

func (menu *Menu) FindCenter() (lowerLeft Point) {
	windowWidthHalf := menu.windowWidth / 2
	windowHeightHalf := menu.windowHeight / 2

	menuWidthHalf := menu.Width / 2
	menuHeightHalf := menu.Height / 2

	lowerLeft.X = float32(windowWidthHalf) - menuWidthHalf
	lowerLeft.Y = float32(windowHeightHalf) - menuHeightHalf
	return
}
