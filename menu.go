package glmenu

import (
	"fmt"
	"image"
	"math"

	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Point struct {
	X, Y float32
}

var vertexShaderSource string = `
#version 330

uniform mat4 scale_matrix;
uniform mat4 matrix;
uniform vec2 final_position;

in vec4 position;

void main() {
  vec4 scaled = scale_matrix * matrix * position;
  gl_Position = vec4(scaled.x + final_position.x, scaled.y + final_position.y, scaled.z, scaled.w);
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
type Navigation int
type Alignment int
type ScreenPosition int

const (
	MouseUnclicked MouseClick = iota
	MouseLeft
	MouseRight
	MouseCenter

	NavigationMouse Navigation = 0
	NavigationKey              = 1

	AlignCenter Alignment = 0
	AlignRight            = 1
	AlignLeft             = 2

	ScreenCenter      ScreenPosition = 0
	ScreenTopLeft                    = 1
	ScreenTopCenter                  = 2
	ScreenTopRight                   = 3
	ScreenLeft                       = 4
	ScreenRight                      = 5
	ScreenLowerLeft                  = 6
	ScreenLowerCenter                = 7
	ScreenLowerRight                 = 8

	ScreenPadding = float32(10) // used by screen positioning calculations
)

type MenuDefaults struct {
	TextColor       mgl32.Vec3
	TextHover       mgl32.Vec3
	TextClick       mgl32.Vec3
	BackgroundColor mgl32.Vec4
	BorderColor     mgl32.Vec4
	Border          mgl32.Vec2
	Dimensions      mgl32.Vec2
	Padding         mgl32.Vec2
	HoverPadding    mgl32.Vec2

	// increment during a scale operation
	TextScaleRate float32
}

type Menu struct {
	// parent MenuManager
	*MenuManager

	Name string

	// trigger
	OnShow         func()
	OnComplete     func()
	OnEnterRelease func()

	// options
	Defaults              MenuDefaults
	IsVisible             bool
	ShowOnKey             glfw.Key
	Height                float32
	Width                 float32
	IsAutoCenter          bool
	lowerLeft, upperRight Point

	// interactive objects
	Font       *v41.Font
	Labels     []*Label
	TextBoxes  []*TextBox
	Formatable []Formatable // all labels and textboxes

	// Up/Down keypress -> NavigationVia set to "Key"
	// When in "Key", mouse navigation only happens once the mouse has been moved enough from LastMousePosition
	LastMousePosition mgl32.Vec2
	NavigationVia     Navigation
	NavigationIndex   int // once up/down arrows are pressed, determine which element needs to be entered/hovered over

	// opengl oriented
	ScreenPosition       ScreenPosition
	screenPositionOffset mgl32.Vec2
	Window               *glfw.Window
	program              uint32 // shader program
	backgroundUniform    int32
	orthographicUniform  int32 // ortho matrix
	finalPositionUniform int32 // offset to reposition based on ScreenPosition
	scaleUniform         int32 // scale matrix

	scaleIdent4   mgl32.Mat4
	scaleMatrix   mgl32.Mat4
	finalPosition mgl32.Vec2
	position      uint32 // index location
	vao           uint32
	vbo           uint32
	ebo           uint32
	vboData       []float32
	vboIndexCount int
	eboData       []int32
	eboIndexCount int
}

func (menu *Menu) Drag(x, y float32) {
	menu.screenPositionOffset[0] += x
	menu.screenPositionOffset[1] += y

	menu.finalPosition[0] = menu.screenPositionOffset.X() / (menu.Font.WindowWidth / 2)
	menu.finalPosition[1] = menu.screenPositionOffset.Y() / (menu.Font.WindowHeight / 2)

	for _, l := range menu.Formatable {
		l.DragPosition(x, y)
	}
}

func (menu *Menu) Finalize(align Alignment) {
	glfloat_size := 4
	glint_size := 4

	menu.format(align)
	menu.finalPosition[0] = menu.screenPositionOffset.X() / (menu.Font.WindowWidth / 2)
	menu.finalPosition[1] = menu.screenPositionOffset.Y() / (menu.Font.WindowHeight / 2)

	// center the menu around (0,0)
	menuWidthHalf := menu.Width / 2
	menuHeightHalf := menu.Height / 2
	menu.lowerLeft.X = -menuWidthHalf
	menu.lowerLeft.Y = -menuHeightHalf
	menu.upperRight.X = menuWidthHalf
	menu.upperRight.Y = menuHeightHalf

	// create the vbo data based on the a centered draw
	menu.makeBufferData(menu.lowerLeft)

	menu.scaleIdent4 = mgl32.Ident4()

	// create a 5 pixel border
	menu.scaleMatrix = mgl32.Scale3D(
		1+menu.Defaults.Border.X()/(menu.Width/2),
		1+menu.Defaults.Border.Y()/(menu.Height/2),
		1,
	)

	// bind data
	gl.BindVertexArray(menu.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, menu.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER, glfloat_size*menu.vboIndexCount, gl.Ptr(menu.vboData), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, menu.ebo)
	gl.BufferData(
		gl.ELEMENT_ARRAY_BUFFER, glint_size*menu.eboIndexCount, gl.Ptr(menu.eboData), gl.DYNAMIC_DRAW)
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

// NewMenuTexture accepts an image location and the internal dimensions of the smaller images embedded within
// the texture.  The image will be subdivided into evenly spaced rectangles based on the dimensions given.
// This is very beta.
func (menu *Menu) NewMenuTexture(imagePath string, dimensions mgl32.Vec2) (mt *MenuTexture, err error) {
	width, height := menu.Window.GetSize()
	mt = &MenuTexture{
		Menu:         menu,
		WindowWidth:  float32(width),
		WindowHeight: float32(height),
		Dimensions:   dimensions,
	}
	mt.ResizeWindow(float32(width), float32(height))

	mt.Image, err = gltext.LoadImage(imagePath)
	if err != nil {
		return nil, err
	}

	// Resize menuTexture to next power-of-two.
	mt.Image = gltext.Pow2Image(mt.Image).(*image.NRGBA)
	ib := mt.Image.Bounds()

	mt.textureWidth = float32(ib.Dx())
	mt.textureHeight = float32(ib.Dy())

	// generate texture
	gl.GenTextures(1, &mt.textureID)
	gl.BindTexture(gl.TEXTURE_2D, mt.textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(ib.Dx()),
		int32(ib.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(mt.Image.Pix),
	)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// create shader program and define attributes and uniforms
	mt.program, err = v41.NewProgram(menuTextureVertexShaderSource, menuTextureFragmentShaderSource)
	if err != nil {
		return mt, err
	}

	// attributes
	mt.centeredPositionAttribute = uint32(gl.GetAttribLocation(mt.program, gl.Str("centered_position\x00")))
	mt.uvAttribute = uint32(gl.GetAttribLocation(mt.program, gl.Str("uv\x00")))

	// uniforms
	mt.finalPositionUniform = gl.GetUniformLocation(mt.program, gl.Str("final_position\x00"))
	mt.orthographicMatrixUniform = gl.GetUniformLocation(mt.program, gl.Str("orthographic_matrix\x00"))
	mt.fragmentTextureUniform = gl.GetUniformLocation(mt.program, gl.Str("fragment_texture\x00"))

	return mt, nil
}

// NewLabel handles spacing layout as defined on the menu level
func (menu *Menu) NewLabel(str string, config LabelConfig) *Label {
	label := &Label{
		Config: config,
		Menu:   menu,
		Text:   v41.NewText(menu.Font, 1.0, 1.1),
	}
	menu.Labels = append(menu.Labels, label)
	menu.Formatable = append(menu.Formatable, label)

	label.SetString(str)
	label.Text.SetScale(1)
	label.Text.SetColor(menu.Defaults.TextColor)

	if config.Action == Noop {
		label.onRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {}
		return label
	}

	label.OnClick = func(xPos, yPos float64, button MouseClick, inBox bool) {
		label.Text.SetColor(menu.Defaults.TextClick)
	}
	label.OnHover = func(xPos, yPos float64, button MouseClick, inBox bool) {
		if !label.IsClick {
			label.Text.SetColor(menu.Defaults.TextHover)
			label.Text.AddScale(menu.Defaults.TextScaleRate)
		}
	}
	label.OnNotHover = func() {
		if !label.IsClick {
			label.Text.SetColor(menu.Defaults.TextColor)
			label.Text.AddScale(-menu.Defaults.TextScaleRate)
		}
	}
	switch config.Action {
	case ExitMenu:
		label.onRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
			if inBox {
				menu.Hide()
			}
		}
	case ExitGame:
		label.onRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
			if inBox {
				menu.Window.SetShouldClose(true)
			}
		}
	default:
		label.onRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {}
	}
	return label
}

func (menu *Menu) format(align Alignment) {
	var borders float32
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
		if l.Width() > wTotal {
			wTotal = l.Width() + l.GetPadding().Y*2
		}
		borders += l.GetPadding().Y * 2
	}
	hTotal = height*float32(length) + borders

	// readjust entire menu size to hold all objects
	if menu.Height < hTotal+menu.Defaults.Padding.Y() {
		menu.Height = hTotal + menu.Defaults.Padding.Y()*2
	}
	if menu.Width < wTotal+menu.Defaults.Padding.X() {
		menu.Width = wTotal + menu.Defaults.Padding.X()*2
	}

	// calculate an appropriate offset based on the screen position that was requested
	// the calculation is based on (0,0) being at the CENTER OF THE SCREEN
	switch menu.ScreenPosition {
	case ScreenTopLeft:
		menu.screenPositionOffset[0] = -(menu.Font.WindowWidth/2 - menu.Width/2) + ScreenPadding
		menu.screenPositionOffset[1] = +(menu.Font.WindowHeight/2 - menu.Height/2) - ScreenPadding
	case ScreenLeft:
		menu.screenPositionOffset[0] = -(menu.Font.WindowWidth/2 - menu.Width/2) + ScreenPadding
		menu.screenPositionOffset[1] = 0
	case ScreenLowerLeft:
		menu.screenPositionOffset[0] = -(menu.Font.WindowWidth/2 - menu.Width/2) + ScreenPadding
		menu.screenPositionOffset[1] = -(menu.Font.WindowHeight/2 - menu.Height/2) + ScreenPadding
	case ScreenTopCenter:
		menu.screenPositionOffset[0] = 0
		menu.screenPositionOffset[1] = +(menu.Font.WindowHeight/2 - menu.Height/2) - ScreenPadding
	case ScreenLowerCenter:
		menu.screenPositionOffset[0] = 0
		menu.screenPositionOffset[1] = -(menu.Font.WindowHeight/2 - menu.Height/2) + ScreenPadding
	case ScreenTopRight:
		menu.screenPositionOffset[0] = +(menu.Font.WindowWidth/2 - menu.Width/2) - ScreenPadding
		menu.screenPositionOffset[1] = +(menu.Font.WindowHeight/2 - menu.Height/2) - ScreenPadding
	case ScreenRight:
		menu.screenPositionOffset[0] = +(menu.Font.WindowWidth/2 - menu.Width/2) - ScreenPadding
		menu.screenPositionOffset[1] = 0
	case ScreenLowerRight:
		menu.screenPositionOffset[0] = +(menu.Font.WindowWidth/2 - menu.Width/2) - ScreenPadding
		menu.screenPositionOffset[1] = -(menu.Font.WindowHeight/2 - menu.Height/2) + ScreenPadding
	}

	// depending on the number of menu elements a vertically centered menus formatting will differ
	// - the middle object in a menu with an odd number of objects has value 0 = middleIndex-float32(i)
	//   or, in the case of even values, the two middle values are just around the center
	middleIndex := float32(math.Floor(float64(length / 2)))
	if length%2 == 0 {
		// even number of objects to vertically align
		for i, l := range menu.Formatable {
			vertical := height + l.GetPadding().Y*2
			yOffset := (middleIndex-float32(i)-1)*vertical + vertical/2
			xOffset := float32(0)
			if align != AlignCenter {
				if l.Type() == FormatableLabel {
					xOffset = -((menu.Width-menu.Defaults.Padding.X()*2-menu.Defaults.HoverPadding.X())/2 - l.Width()/2)
				} else {
					xOffset = -((menu.Width-menu.Defaults.Padding.X()*2)/2 - l.Width()/2)
				}
				if align == AlignRight {
					xOffset = -xOffset
				}
			}
			l.SetPosition(mgl32.Vec2{xOffset + l.GetPosition().X() + menu.screenPositionOffset.X(), yOffset + l.GetPosition().Y() + menu.screenPositionOffset.Y()})
		}
	} else {
		for i, l := range menu.Formatable {
			vertical := height + l.GetPadding().Y*2
			yOffset := (middleIndex - float32(i)) * vertical
			xOffset := float32(0)
			if align != AlignCenter {
				if l.Type() == FormatableLabel {
					xOffset = -((menu.Width-menu.Defaults.Padding.X()*2-menu.Defaults.HoverPadding.X())/2 - l.Width()/2)
				} else {
					xOffset = -((menu.Width-menu.Defaults.Padding.X()*2)/2 - l.Width()/2)
				}
				if align == AlignRight {
					xOffset = -xOffset
				}
			}
			l.SetPosition(mgl32.Vec2{xOffset + l.GetPosition().X() + menu.screenPositionOffset.X(), yOffset + l.GetPosition().Y() + menu.screenPositionOffset.Y()})
		}
	}
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
		menu.NavigationVia = NavigationMouse
		menu.NavigationIndex = -1
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
func NewMenu(window *glfw.Window, name string, font *v41.Font, defaults MenuDefaults, screenPosition ScreenPosition) (*Menu, error) {
	// i believe we are actually supposed to pass in the framebuffer sizes when creating the orthographic projection
	// this would probably require some changes though in order to track mouse movement.
	width, height := window.GetSize()
	menu := &Menu{
		Name:           name,
		Defaults:       defaults,
		Font:           font,
		IsVisible:      false,
		ShowOnKey:      glfw.KeyM,
		Width:          defaults.Dimensions.X(),
		Height:         defaults.Dimensions.Y(),
		Window:         window,
		ScreenPosition: screenPosition,
	}
	menu.Font.ResizeWindow(float32(width), float32(height))

	// reasonable default is to follow the first followable element when hitting enter i suppose
	menu.OnEnterRelease = func() {
		if menu.IsVisible {
			for i := range menu.Formatable {
				if menu.Formatable[i].Follow() {
					return
				}
			}
		}
	}

	// create shader program and define attributes and uniforms
	var err error
	menu.program, err = v41.NewProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}
	menu.scaleUniform = gl.GetUniformLocation(menu.program, gl.Str("scale_matrix\x00"))
	menu.orthographicUniform = gl.GetUniformLocation(menu.program, gl.Str("matrix\x00"))
	menu.finalPositionUniform = gl.GetUniformLocation(menu.program, gl.Str("final_position\x00"))
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

	// binding of the data now in the Finalize method
	return menu, nil
}

func (menu *Menu) makeBufferData(lowerLeft Point) {
	// index (0,0)
	menu.vboData[0] = lowerLeft.X // position
	menu.vboData[1] = lowerLeft.Y

	// index (1,0)
	menu.vboData[2] = lowerLeft.X + menu.Width
	menu.vboData[3] = lowerLeft.Y

	// index (1,1)
	menu.vboData[4] = lowerLeft.X + menu.Width
	menu.vboData[5] = lowerLeft.Y + menu.Height

	// index (0,1)
	menu.vboData[6] = lowerLeft.X
	menu.vboData[7] = lowerLeft.Y + menu.Height

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
	if !menu.MenuManager.IsFinalized {
		panic(fmt.Sprintf("menu %s's MenuManager must be finalized prior to drawing!", menu.Name))
	}
	if !menu.IsVisible {
		return menu.IsVisible
	}
	gl.UseProgram(menu.program)

	for i := 0; i < 2; i++ {
		// i == 0 is background draw at higher scale, producing a border around the menu
		if i == 0 {
			gl.UniformMatrix4fv(menu.scaleUniform, 1, false, &menu.scaleMatrix[0])
			gl.Uniform4fv(menu.backgroundUniform, 1, &menu.Defaults.BorderColor[0])
		} else {
			gl.UniformMatrix4fv(menu.scaleUniform, 1, false, &menu.scaleIdent4[0])
			gl.Uniform4fv(menu.backgroundUniform, 1, &menu.Defaults.BackgroundColor[0])
		}
		gl.Uniform2fv(menu.finalPositionUniform, 1, &menu.finalPosition[0])
		gl.UniformMatrix4fv(menu.orthographicUniform, 1, false, &menu.Font.OrthographicMatrix[0])

		gl.Enable(gl.BLEND)
		gl.BlendEquation(gl.FUNC_ADD)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		gl.BindVertexArray(menu.vao)
		gl.DrawElements(gl.TRIANGLES, int32(menu.eboIndexCount), gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)
		gl.Disable(gl.BLEND)
	}

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
	mX, mY := ScreenCoordToCenteredCoord(menu.Font.WindowWidth, menu.Font.WindowHeight, xPos, yPos)
	if inBox := InBox(mX, mY, menu); !inBox {
		return
	}
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
	dist := math.Sqrt(math.Pow(float64(menu.LastMousePosition[0])-xPos, 2) + math.Pow(float64(menu.LastMousePosition[1])-yPos, 2))
	menu.LastMousePosition[0] = float32(xPos)
	menu.LastMousePosition[1] = float32(yPos)

	if dist > 1 {
		// a bit of mouse movement will reenable mouse position evaluation
		menu.NavigationVia = NavigationMouse
		menu.NavigationIndex = -1
	}
	if menu.NavigationVia == NavigationKey {
		menu.Formatable[menu.NavigationIndex].NavigateTo()
		return
	}
	for i := range menu.Labels {
		if menu.Labels[i].OnHover != nil {
			menu.Labels[i].IsHovered(xPos, yPos)
		}
	}
}

func (menu *Menu) KeyRelease(key glfw.Key, withShift bool) {
	if key == glfw.KeyUp || key == glfw.KeyDown {
		for i := range menu.Formatable {
			if menu.Formatable[i].NavigateAway() {
				menu.NavigationIndex = i
			}
		}
		// adjust endpoints skipping objects that are NOOP
		for {
			if key == glfw.KeyUp {
				menu.NavigationIndex -= 1
			} else {
				menu.NavigationIndex += 1
			}
			if menu.NavigationIndex < 0 || menu.NavigationIndex == len(menu.Formatable) || !menu.Formatable[menu.NavigationIndex].IsNoop() {
				if menu.NavigationIndex < -1 {
					menu.NavigationIndex = -1
				}
				// once we have found something that can be interacted with or have gone beyond our boundaries, break the for loop
				break
			}
		}
		// if we ended up going beyond a boundary try to go back into menu entries to find something that can be interacted with
		// the worst case scenario is that everything is NOOP in which case we just need to exit
		if menu.NavigationIndex < 0 {
			// went up the menu too much now go back down
			for {
				menu.NavigationIndex++
				if menu.NavigationIndex == len(menu.Formatable) || !menu.Formatable[menu.NavigationIndex].IsNoop() {
					break
				}
			}
		} else if menu.NavigationIndex == len(menu.Formatable) {
			// went down the menu too much now go back up
			for {
				menu.NavigationIndex--
				if menu.NavigationIndex < 0 || !menu.Formatable[menu.NavigationIndex].IsNoop() {
					break
				}
			}
		}
		if menu.NavigationIndex < 0 {
			menu.NavigationIndex = 0
		}
		if menu.NavigationIndex == len(menu.Formatable) {
			menu.NavigationIndex -= 1
		}

		// perform necessary visual changes as we navigate to the next place
		for i := range menu.Formatable {
			if i == menu.NavigationIndex {
				menu.Formatable[i].NavigateTo()
			}
		}
		menu.NavigationVia = NavigationKey
	}
	for i := range menu.TextBoxes {
		menu.TextBoxes[i].KeyRelease(key, withShift)
	}
	if menu.OnEnterRelease != nil && key == glfw.KeyEnter {
		menu.OnEnterRelease()
	}
}

func (menu *Menu) GetBoundingBox() (X1, X2 gltext.Point) {
	x, y := menu.screenPositionOffset.X(), menu.screenPositionOffset.Y()

	X1.X = menu.lowerLeft.X + x
	X1.Y = menu.lowerLeft.Y + y

	X2.X = menu.upperRight.X + x
	X2.Y = menu.upperRight.Y + y
	return
}
