package glmenu

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/math/fixed"
	"os"
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

uniform vec4 background;
out vec4 fragment_color;

void main() {
  fragment_color = background;
}
` + "\x00"

type MouseClick int

const (
	MouseUnclicked MouseClick = iota
	MouseLeft
	MouseRight
	MouseCenter
)

type Menu struct {
	//trigger
	OnShow         func()
	OnEnterRelease func()

	// options
	Visible      bool
	ShowOnKey    glfw.Key
	Height       float32
	Width        float32
	IsAutoCenter bool
	lowerLeft    Point

	backgroundUniform int32
	Background        mgl32.Vec4

	// interactive objects
	Font          *gltext.Font
	Labels        []*Label
	TextBoxes     []*TextBox
	TextScaleRate float32 // increment during a scale operation

	// opengl oriented
	WindowWidth   float32
	WindowHeight  float32
	program       uint32 // shader program
	glMatrix      int32  // ortho matrix
	position      uint32 // index location
	vao           uint32
	vbo           uint32
	ebo           uint32
	vboData       []float32
	vboIndexCount int
	eboData       []int32
	eboIndexCount int
}

func (menu *Menu) NewLabel(str string) *Label {
	label := &Label{
		Menu: menu,
		Text: gltext.NewText(menu.Font, 1.0, 1.1),
	}
	menu.Labels = append(menu.Labels, label)
	label.SetString(str)
	label.Text.SetScale(1)
	return label
}

func (menu *Menu) AddTextBox(textbox *TextBox, str string, width int32, height int32, borderWidth int32) {
	textbox.Load(menu, width, height, borderWidth)
	textbox.SetString(str)
	textbox.Text.SetScale(1)
	textbox.Text.SetPosition(0, 0)
	textbox.Text.SetColor(0, 0, 0)
	menu.TextBoxes = append(menu.TextBoxes, textbox)
}

func (menu *Menu) Show() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.Visible = true
	if menu.OnShow != nil {
		menu.OnShow()
	}
}

func (menu *Menu) Hide() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.Visible = false
}

func (menu *Menu) Toggle() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.Visible = !menu.Visible
}

// Load will draw a background centered on the screen or positioned based on offsetBy values
func NewMenu(width float32, height float32, scale fixed.Int26_6, offsetBy mgl32.Vec2) (*Menu, error) {
	glfloat_size := 4
	glint_size := 4

	menu := &Menu{}
	menu.Visible = false
	menu.ShowOnKey = glfw.KeyM
	menu.Width = width
	menu.Height = height

	// load font
	fd, err := os.Open("font/luximr.ttf")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	runesPerRow := fixed.Int26_6(16)
	runeRanges := make(gltext.RuneRanges, 0)
	runeRange := gltext.RuneRange{Low: 32, High: 127}
	runeRanges = append(runeRanges, runeRange)

	menu.Font, err = gltext.NewTruetype(fd, scale, runeRanges, runesPerRow)
	if err != nil {
		return nil, err
	}

	// 2DO: make this time dependent rather than fps dependent
	menu.TextScaleRate = 0.01

	// create shader program and define attributes and uniforms
	menu.program, err = gltext.NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}
	menu.glMatrix = gl.GetUniformLocation(menu.program, gl.Str("matrix\x00"))
	menu.backgroundUniform = gl.GetUniformLocation(menu.program, gl.Str("background\x00"))
	menu.position = uint32(gl.GetAttribLocation(menu.program, gl.Str("position\x00")))

	gl.GenVertexArrays(1, &menu.vao)
	gl.GenBuffers(1, &menu.vbo)
	gl.GenBuffers(1, &menu.ebo)

	// vao
	gl.BindVertexArray(menu.vao)

	// 2DO: Change text depth to get it to render? For now this works.
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

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
	menu.lowerLeft = menu.findCenter(offsetBy)
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
	return menu, nil
}

func (menu *Menu) ResizeWindow(width float32, height float32) {
	menu.WindowWidth = width
	menu.WindowHeight = height
	menu.Font.ResizeWindow(width, height)
}

func (menu *Menu) makeBufferData() {
	// index (0,0)
	menu.vboData[0] = menu.lowerLeft.X // position
	menu.vboData[1] = menu.lowerLeft.Y

	// index (1,0)
	menu.vboData[2] = menu.lowerLeft.X + menu.Width
	menu.vboData[3] = menu.lowerLeft.Y

	// index (1,1)
	menu.vboData[4] = menu.lowerLeft.X + menu.Width
	menu.vboData[5] = menu.lowerLeft.Y + menu.Height

	// index (0,1)
	menu.vboData[6] = menu.lowerLeft.X
	menu.vboData[7] = menu.lowerLeft.Y + menu.Height

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
	for i := range menu.Labels {
		menu.Labels[i].Text.Release()
		if menu.Labels[i].Shadow != nil && menu.Labels[i].Shadow.Text != nil {
			menu.Labels[i].Shadow.Text.Release()
		}
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].Text.Release()
	}
}

func (menu *Menu) Draw() bool {
	if !menu.Visible {
		return menu.Visible
	}
	gl.UseProgram(menu.program)

	gl.UniformMatrix4fv(menu.glMatrix, 1, false, &menu.Font.OrthographicMatrix[0])
	gl.Uniform4fv(menu.backgroundUniform, 1, &menu.Background[0])

	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.BindVertexArray(menu.vao)
	gl.DrawElements(gl.TRIANGLES, int32(menu.eboIndexCount), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
	gl.Disable(gl.BLEND)

	for i := range menu.Labels {
		if !menu.Labels[i].IsHover {
			if menu.Labels[i].OnNotHover != nil {
				menu.Labels[i].OnNotHover(menu.Labels[i])
				if menu.Labels[i].Shadow != nil {
					menu.Labels[i].OnNotHover(&menu.Labels[i].Shadow.Label)
				}
			}
		}
		menu.Labels[i].Draw()
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].Draw()
	}
	return menu.Visible
}

func (menu *Menu) MouseClick(xPos, yPos float64, button MouseClick) {
	if !menu.Visible {
		return
	}
	yPos = float64(menu.WindowHeight) - yPos
	for i := range menu.Labels {
		menu.Labels[i].IsClicked(xPos, yPos, button)
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].IsClicked(xPos, yPos, button)
	}
}

func (menu *Menu) MouseRelease(xPos, yPos float64, button MouseClick) {
	if !menu.Visible {
		return
	}
	yPos = float64(menu.WindowHeight) - yPos
	for i := range menu.Labels {
		menu.Labels[i].IsReleased(xPos, yPos, button)
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].IsReleased(xPos, yPos, button)
	}
}

func (menu *Menu) MouseHover(xPos, yPos float64) {
	if !menu.Visible {
		return
	}
	yPos = float64(menu.WindowHeight) - yPos
	for i := range menu.Labels {
		if menu.Labels[i].OnHover != nil {
			menu.Labels[i].IsHovered(xPos, yPos)
		}
	}
}

func (menu *Menu) findCenter(offsetBy mgl32.Vec2) (lowerLeft Point) {
	menuWidthHalf := menu.Width / 2
	menuHeightHalf := menu.Height / 2

	lowerLeft.X = -menuWidthHalf + offsetBy.X()
	lowerLeft.Y = -menuHeightHalf + offsetBy.Y()
	return
}

func (menu *Menu) KeyRelease(key glfw.Key, withShift bool) {
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].KeyRelease(key, withShift)
	}
	if menu.OnEnterRelease != nil && key == glfw.KeyEnter {
		menu.OnEnterRelease()
	}
}
