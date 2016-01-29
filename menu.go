package glmenu

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
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

type MenuDefaults struct {
	TextColor       mgl32.Vec3
	TextHover       mgl32.Vec3
	TextClick       mgl32.Vec3
	BackgroundColor mgl32.Vec4
	Dimensions      mgl32.Vec2
}

type Menu struct {
	//trigger
	OnShow         func()
	OnEnterRelease func()

	// options
	Defaults     MenuDefaults
	IsVisible    bool
	ShowOnKey    glfw.Key
	Height       float32
	Width        float32
	IsAutoCenter bool
	lowerLeft    Point

	backgroundUniform int32
	Background        mgl32.Vec4

	// interactive objects
	Font        *gltext.Font
	Labels      []*Label
	LabelBorder Border
	TextBoxes   []*TextBox
	Formatable  []Formatable

	// increment during a scale operation
	TextScaleRate float32

	// opengl oriented
	Offset        mgl32.Vec2
	Window        *glfw.Window
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

func (menu *Menu) Finalize() {
	glfloat_size := 4
	glint_size := 4

	menu.format()
	menu.lowerLeft = menu.findCenter(menu.Offset)
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
}

// NewLabel handles spacing layout as defined on the menu level
func (menu *Menu) NewLabel(str string, config LabelConfig) *Label {
	label := &Label{
		Config: config,
		Menu:   menu,
		Text:   gltext.NewText(menu.Font, 1.0, 1.1),
	}
	menu.Labels = append(menu.Labels, label)
	menu.Formatable = append(menu.Formatable, label)

	label.SetString(str)
	label.Text.SetScale(1)
	label.Text.SetColor(menu.Defaults.TextColor)

	label.OnClick = func(xPos, yPos float64, button MouseClick, inBox bool) {
		label.Text.SetColor(menu.Defaults.TextClick)
	}
	label.OnHover = func(xPos, yPos float64, button MouseClick, inBox bool) {
		if !label.IsClick {
			label.Text.SetColor(menu.Defaults.TextHover)
			label.Text.AddScale(menu.TextScaleRate)
		}
	}
	label.OnNotHover = func() {
		if !label.IsClick {
			label.Text.SetColor(menu.Defaults.TextColor)
			label.Text.AddScale(-menu.TextScaleRate)
		}
	}
	switch config.Action {
	case EXIT_MENU:
		label.OnRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
			if inBox {
				menu.Hide()
			}
		}
	case EXIT_GAME:
		label.OnRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
			if inBox {
				menu.Window.SetShouldClose(true)
			}
		}
	}
	return label
}

func (menu *Menu) format() {
	height, width := float32(0), float32(0)
	hTotal, wTotal := float32(0), float32(0)
	length := len(menu.Formatable)
	for _, l := range menu.Formatable {
		if l.Height() > height {
			height = l.Height()
		}
		if l.Width() > width {
			width = l.Width()
		}
		hTotal += l.Height()
		wTotal += l.Width()
	}

	// not easily understood perhaps - formatting these things never is!
	// depending on the number of menu elements a vertically centered menus formatting will differ
	// - the middle object in a menu with an odd number of objects has value 0 = middleIndex-float32(i)
	//   or, in the case of even values, the two middle values are just around the center
	middleIndex := float32(math.Floor(float64(length / 2)))
	vertical := height + menu.LabelBorder.Y*2
	switch length % 2 {
	case 0:
		for i, l := range menu.Formatable {
			offset := (middleIndex-float32(i)-1)*vertical + vertical/2
			l.SetPosition(0, offset)
		}
	case 1:
		for i, l := range menu.Formatable {
			offset := (middleIndex - float32(i)) * vertical
			l.SetPosition(0, offset)
		}
	}
	// 2do: get total height... and width. make sure that the background is at least large enough to cover this area...
	// 2do: repositioning the menu in different places on the screen
}

// NewTextBox handles vertical spacing
func (menu *Menu) NewTextBox(str string, width, height float32, borderWidth int32) *TextBox {
	textbox := &TextBox{}
	textbox.Load(menu, width, height, borderWidth)
	textbox.SetString(str)
	textbox.SetColor(menu.Defaults.TextColor)
	textbox.Text.SetScale(1)

	menu.TextBoxes = append(menu.TextBoxes, textbox)
	menu.Formatable = append(menu.Formatable, textbox)
	return textbox
}

func (menu *Menu) Show() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.IsVisible = true
	if menu.OnShow != nil {
		menu.OnShow()
	}
}

func (menu *Menu) Hide() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.IsVisible = false
}

func (menu *Menu) Toggle() {
	for i := range menu.Labels {
		menu.Labels[i].Reset()
	}
	menu.IsVisible = !menu.IsVisible
}

// NewMenu creates a new menu object with a background centered on the screen or positioned using offsetBy
func NewMenu(window *glfw.Window, font *gltext.Font, defaults MenuDefaults, offsetBy mgl32.Vec2) (*Menu, error) {
	//glfloat_size := 4
	//glint_size := 4

	// i believe we are actually supposed to pass in the framebuffer sizes when creating the orthographic projection
	// this would probably require some changes though in order to track mouse movement.
	width, height := window.GetSize()
	menu := &Menu{
		Defaults:  defaults,
		Font:      font,
		IsVisible: false,
		ShowOnKey: glfw.KeyM,
		Width:     defaults.Dimensions.X(),
		Height:    defaults.Dimensions.Y(),
		Window:    window,
		Offset:    offsetBy,
	}
	menu.Background = defaults.BackgroundColor
	menu.TextScaleRate = 0.01 // 2DO: make this time dependent rather than fps dependent?
	menu.ResizeWindow(float32(width), float32(height))

	// create shader program and define attributes and uniforms
	var err error
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

	/*
		// positioning and dimensions
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
	*/
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
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].Text.Release()
	}
}

func (menu *Menu) Draw() bool {
	if !menu.IsVisible {
		return menu.IsVisible
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
				menu.Labels[i].OnNotHover()
			}
		}
		menu.Labels[i].Draw()
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].Draw()
	}
	return menu.IsVisible
}

func (menu *Menu) MouseClick(xPos, yPos float64, button MouseClick) {
	if !menu.IsVisible {
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
	if !menu.IsVisible {
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
	if !menu.IsVisible {
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
