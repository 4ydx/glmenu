package glmenu

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

var textboxVertexShader string = `
#version 330

uniform mat4 orthographic_matrix;
uniform vec2 final_position;

in vec4 centered_position;

void main() {
  vec4 center = orthographic_matrix * centered_position;
  gl_Position = vec4(center.x + final_position.x, center.y + final_position.y, center.z, center.w);
}
` + "\x00"

var textboxFragmentShader string = `
#version 330

uniform vec3 background;
out vec4 fragment_color;

void main() {
	fragment_color = vec4(background, 1);
}
` + "\x00"

type TextBoxInteraction func(
	textbox *TextBox,
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

type TextBox struct {
	Menu               *Menu
	Text               *gltext.Text
	MaxLength          int
	CursorBarFrequency int64
	Time               time.Time
	IsEdit             bool
	IsClick            bool

	// user defined
	OnClick    TextBoxInteraction
	OnRelease  TextBoxInteraction
	FilterRune func(r rune) bool

	// opengl oriented
	program          uint32
	glMatrix         int32
	position         uint32
	vao              uint32
	vbo              uint32
	ebo              uint32
	vboData          []float32
	vboIndexCount    int
	eboData          []int32
	eboIndexCount    int
	centeredPosition uint32

	backgroundUniform         int32
	borderBackground          mgl32.Vec3
	textBackground            mgl32.Vec3
	finalPositionUniform      int32
	finalPosition             mgl32.Vec2
	orthographicMatrixUniform int32

	// X1, X2: the lower left and upper right points of a box that bounds the text
	X1          Point
	X2          Point
	BorderWidth int32
	Height      int32
	Width       int32

	SetPositionX float32
	SetPositionY float32
}

func (textbox *TextBox) Load(menu *Menu, width int32, height int32, borderWidth int32) (err error) {
	textbox.Menu = menu

	// text
	textbox.CursorBarFrequency = time.Duration.Nanoseconds(500000000)
	textbox.Text = gltext.NewText(menu.Font, 1.0, 1.1)

	// border formatting
	textbox.BorderWidth = borderWidth
	textbox.Height = height
	textbox.Width = width
	textbox.X1.X = -float32(width) / 2.0
	textbox.X1.Y = -float32(height) / 2.0
	textbox.X2.X = float32(width) / 2.0
	textbox.X2.Y = float32(height) / 2.0
	textbox.borderBackground = mgl32.Vec3{1.0, 1.0, 1.0}
	textbox.textBackground = mgl32.Vec3{0.0, 0.0, 0.0}

	// create shader program and define attributes and uniforms
	textbox.program, err = gltext.NewProgram(textboxVertexShader, textboxFragmentShader)
	if err != nil {
		return err
	}

	// ebo, vbo data
	// 4 border edges (4 vertices apiece with 2 position points per index)
	// 1 background square (4 vertices x 2 positions)
	textbox.vboIndexCount = 4*4*2 + 1*4*2
	textbox.eboIndexCount = 4*6 + 1*6
	textbox.vboData = make([]float32, textbox.vboIndexCount, textbox.vboIndexCount)
	textbox.eboData = make([]int32, textbox.eboIndexCount, textbox.eboIndexCount)
	textbox.makeBufferData()

	// attributes
	textbox.centeredPosition = uint32(gl.GetAttribLocation(textbox.program, gl.Str("centered_position\x00")))

	// uniforms
	textbox.backgroundUniform = gl.GetUniformLocation(textbox.program, gl.Str("background\x00"))
	textbox.finalPositionUniform = gl.GetUniformLocation(textbox.program, gl.Str("final_position\x00"))
	textbox.orthographicMatrixUniform = gl.GetUniformLocation(textbox.program, gl.Str("orthographic_matrix\x00"))

	gl.GenVertexArrays(1, &textbox.vao)
	gl.GenBuffers(1, &textbox.vbo)
	gl.GenBuffers(1, &textbox.ebo)

	glfloatSize := int32(4)

	// vao
	gl.BindVertexArray(textbox.vao)

	// vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, textbox.vbo)

	gl.EnableVertexAttribArray(textbox.centeredPosition)
	gl.VertexAttribPointer(
		textbox.centeredPosition,
		2,
		gl.FLOAT,
		false,
		0,
		gl.PtrOffset(0),
	)
	gl.BufferData(gl.ARRAY_BUFFER, int(glfloatSize)*textbox.vboIndexCount, gl.Ptr(textbox.vboData), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, textbox.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(glfloatSize)*textbox.eboIndexCount, gl.Ptr(textbox.eboData), gl.DYNAMIC_DRAW)
	gl.BindVertexArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	return
}

// it is probably best to draw a diagram of what is happening rather than to try to read this code
// X1: lower left hand point
// X2: upper right hand point
// the textbox border is being created by drawing left and right edges whose height include the border
// the top and bottom edges horizontal width does not include the border
// keep in mind that each drawn edge is itself its own quad CCW from upper right hand corner
func (textbox *TextBox) makeBufferData() {
	// this all works because the original positioning is centered around the origin

	// left edge - positions starting at upper right CCW
	textbox.vboData[0] = textbox.X1.X
	textbox.vboData[1] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[2] = textbox.X1.X - float32(textbox.BorderWidth)
	textbox.vboData[3] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[4] = textbox.X1.X - float32(textbox.BorderWidth)
	textbox.vboData[5] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.vboData[6] = textbox.X1.X
	textbox.vboData[7] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.eboData[0], textbox.eboData[1], textbox.eboData[2], textbox.eboData[3], textbox.eboData[4], textbox.eboData[5] = 0, 1, 2, 0, 2, 3

	// top edge - intentionally leaves out the borderwidth on the x-axis
	textbox.vboData[8] = textbox.X2.X
	textbox.vboData[9] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[10] = textbox.X1.X
	textbox.vboData[11] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[12] = textbox.X1.X
	textbox.vboData[13] = textbox.X2.Y
	textbox.vboData[14] = textbox.X2.X
	textbox.vboData[15] = textbox.X2.Y
	textbox.eboData[6], textbox.eboData[7], textbox.eboData[8], textbox.eboData[9], textbox.eboData[10], textbox.eboData[11] = 4, 5, 6, 4, 6, 7

	// bottom edge - intentionally leaves out the borderwidth on the x-axis
	textbox.vboData[16] = textbox.X2.X
	textbox.vboData[17] = textbox.X1.Y
	textbox.vboData[18] = textbox.X1.X
	textbox.vboData[19] = textbox.X1.Y
	textbox.vboData[20] = textbox.X1.X
	textbox.vboData[21] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.vboData[22] = textbox.X2.X
	textbox.vboData[23] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.eboData[12], textbox.eboData[13], textbox.eboData[14], textbox.eboData[15], textbox.eboData[16], textbox.eboData[17] = 8, 9, 10, 8, 10, 11

	// right edge
	textbox.vboData[24] = textbox.X2.X + float32(textbox.BorderWidth)
	textbox.vboData[25] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[26] = textbox.X2.X
	textbox.vboData[27] = textbox.X2.Y + float32(textbox.BorderWidth)
	textbox.vboData[28] = textbox.X2.X
	textbox.vboData[29] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.vboData[30] = textbox.X2.X + float32(textbox.BorderWidth)
	textbox.vboData[31] = textbox.X1.Y - float32(textbox.BorderWidth)
	textbox.eboData[18], textbox.eboData[19], textbox.eboData[20], textbox.eboData[21], textbox.eboData[22], textbox.eboData[23] = 12, 13, 14, 12, 14, 15

	// background
	textbox.vboData[32] = textbox.X2.X
	textbox.vboData[33] = textbox.X2.Y
	textbox.vboData[34] = textbox.X1.X
	textbox.vboData[35] = textbox.X2.Y
	textbox.vboData[36] = textbox.X1.X
	textbox.vboData[37] = textbox.X1.Y
	textbox.vboData[38] = textbox.X2.X
	textbox.vboData[39] = textbox.X1.Y
	textbox.eboData[24], textbox.eboData[25], textbox.eboData[26], textbox.eboData[27], textbox.eboData[28], textbox.eboData[29] = 16, 17, 18, 16, 18, 19
}

func (textbox *TextBox) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		textbox.Text.SetString(str + "|")
	} else {
		textbox.Text.SetString(str+"|", argv)
	}
}

//TODO bar needs to be independent from the actual text
func (textbox *TextBox) Draw() {
	if time.Since(textbox.Time).Nanoseconds() > textbox.CursorBarFrequency {
		if textbox.Text.RuneCount < textbox.Text.GetLength() {
			textbox.Text.RuneCount = textbox.Text.GetLength()
		} else {
			textbox.Text.RuneCount -= 1
		}
		textbox.Time = time.Now()
	}
	if !textbox.IsEdit {
		// dont show flashing bar unless actually editing
		textbox.Text.RuneCount = textbox.Text.GetLength() - 1
	}
	gl.UseProgram(textbox.program)

	// draw
	gl.BindVertexArray(textbox.vao)

	// uniforms
	gl.Uniform2fv(textbox.finalPositionUniform, 1, &textbox.finalPosition[0])
	gl.UniformMatrix4fv(textbox.orthographicMatrixUniform, 1, false, &textbox.Menu.Font.OrthographicMatrix[0])

	// draw border - 4 * 6: four quads with six indices apiece starting at the beginning of the vbo (0)
	gl.Uniform3fv(textbox.backgroundUniform, 1, &textbox.borderBackground[0])
	gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(4*6), gl.UNSIGNED_INT, nil, int32(0))

	// draw background - start drawing after skipping the border vertices (16)
	gl.Uniform3fv(textbox.backgroundUniform, 1, &textbox.textBackground[0])
	gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(1*6), gl.UNSIGNED_INT, nil, int32(16))
	gl.BindVertexArray(0)

	textbox.Text.Draw()
}

func (textbox *TextBox) KeyRelease(key glfw.Key, withShift bool) {
	if textbox.IsEdit {
		switch key {
		case glfw.KeyBackspace:
			textbox.Backspace()
		case glfw.KeyEscape:
			textbox.IsEdit = false
		default:
			textbox.AddRune(key, withShift)
		}
	}
}

func (textbox *TextBox) AddRune(key glfw.Key, withShift bool) {
	if textbox.Text.HasRune(rune(key)) {
		processRune := true
		if textbox.FilterRune != nil {
			processRune = textbox.FilterRune(rune(key))
		}
		if processRune {
			var theRune rune
			if !withShift && key >= 65 && key <= 90 {
				theRune = rune(key) + 32
			} else {
				theRune = rune(key)
			}
			if textbox.Text.MaxRuneCount > 0 && len(textbox.Text.String) == textbox.Text.MaxRuneCount {
				// too long
			} else {
				r := []rune(textbox.Text.String)
				r = r[0 : len(r)-1] // trim the bar
				r = append(r, theRune)
				textbox.Text.SetString(string(r) + "|")
				textbox.Text.SetPosition(textbox.Text.SetPositionX, textbox.Text.SetPositionY)
			}
		}
	}
}

func (textbox *TextBox) SetPosition(x, y float32) {
	// transform to orthographic coordinates ranged -1 to 1 for the shader
	textbox.finalPosition[0] = x / (textbox.Menu.Font.WindowWidth / 2)
	textbox.finalPosition[1] = y / (textbox.Menu.Font.WindowHeight / 2)

	// used for detecting clicks, hovers, etc
	textbox.X1.X += x
	textbox.X1.Y += y
	textbox.X2.X += x
	textbox.X2.Y += y

	// used to build shadow data and for calling SetPosition again when needed
	textbox.SetPositionX = x
	textbox.SetPositionY = y
	textbox.Text.SetPosition(x, y)
}

func (textbox *TextBox) Backspace() {
	r := []rune(textbox.Text.String)
	if len(r) > 1 {
		r = r[0 : len(r)-2]
		// this will recenter the textbox on the screen
		textbox.Text.SetString(string(r) + "|")
		// this will place it back where it was previously positioned
		textbox.Text.SetPosition(textbox.Text.SetPositionX, textbox.Text.SetPositionY)
	}
}

func (textbox *TextBox) OrthoToScreenCoord() (X1 Point, X2 Point) {
	X1.X = textbox.X1.X + textbox.Menu.WindowWidth/2
	X1.Y = textbox.X1.Y + textbox.Menu.WindowHeight/2

	X2.X = textbox.X2.X + textbox.Menu.WindowWidth/2
	X2.Y = textbox.X2.Y + textbox.Menu.WindowHeight/2
	return
}

// typically called by the menu object handling the label
func (textbox *TextBox) IsClicked(xPos, yPos float64, button MouseClick) {
	// menu rendering (and text) is positioned in orthographic projection coordinates
	// but click positions are based on window coordinates
	// we have to transform them
	X1, X2 := textbox.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if inBox {
		textbox.IsClick = true
		if textbox.OnClick != nil {
			textbox.OnClick(textbox, xPos, yPos, button, inBox)
		}
	} else {
		textbox.IsEdit = false
	}
}

// typically called by the menu object handling the label
func (textbox *TextBox) IsReleased(xPos, yPos float64, button MouseClick) {
	// anything flagged as clicked now needs to decide whether to execute its logic based on inBox
	X1, X2 := textbox.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if textbox.IsClick {
		textbox.IsEdit = true
		if textbox.OnRelease != nil {
			textbox.OnRelease(textbox, xPos, yPos, button, inBox)
		}
	}
	textbox.IsClick = false
}
