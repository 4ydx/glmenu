package menu

import (
	gltext "github.com/4ydx/gltext"
)

type Interaction func(
	label *Label,
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

type Label struct {
	Menu       *Menu
	Text       *gltext.Text
	OnClick    Interaction
	IsClick    bool
	OnRelease  Interaction
	OnHover    Interaction
	OnNotHover func(label *Label)
	IsHover    bool
	Shadow     *Shadow
}

type Shadow struct {
	Label
	Offset float32
}

func (label *Label) AddShadow(offset, r, g, b float32) {
	label.Shadow = new(Shadow)
	label.Shadow.Menu = label.Menu
	label.UpdateShadow(offset, r, g, b)
}

func (label *Label) UpdateShadow(offset, r, g, b float32) {
	label.Shadow.Text = gltext.LoadText(label.Menu.Font)
	label.Shadow.Text.SetColor(r, g, b)
	label.Shadow.Text.SetString(label.Text.String)
	label.Shadow.Text.SetPosition(label.Text.SetPositionX+offset, label.Text.SetPositionY+offset)

	label.Shadow.OnClick = label.OnClick
	label.Shadow.OnHover = label.OnHover
	label.Shadow.OnNotHover = label.OnNotHover
}

func (label *Label) Reset() {
	label.Text.SetScale(label.Text.ScaleMin)
	if label.Shadow != nil {
		label.Shadow.Text.SetScale(label.Text.ScaleMin)
	}
}

func (label *Label) Load(menu *Menu, font *gltext.Font) {
	label.Menu = menu
	label.Text = gltext.LoadText(font)
}

func (label *Label) SetString(str string, argv ...interface{}) {
	label.Text.SetString(str, argv)
	if label.Shadow != nil {
		label.Shadow.Text.SetString(str, argv)
	}
}

func (label *Label) OrthoToScreenCoord() (X1 Point, X2 Point) {
	X1.X = label.Text.X1.X + label.Menu.WindowWidth/2
	X1.Y = label.Text.X1.Y + label.Menu.WindowHeight/2

	X2.X = label.Text.X2.X + label.Menu.WindowWidth/2
	X2.Y = label.Text.X2.Y + label.Menu.WindowHeight/2
	return
}

func (label *Label) IsClicked(xPos, yPos float64, button MouseClick) {
	// menu rendering (and text) is positioned in orthographic projection coordinates
	// but click positions are based on window coordinates
	// we have to transform them
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if inBox {
		label.IsClick = true
		label.OnClick(label, xPos, yPos, button, inBox)
	}
}

func (label *Label) IsReleased(xPos, yPos float64, button MouseClick) {
	// anything flagged as clicked now needs to decide whether to execute its logic based on inBox
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if label.IsClick {
		label.OnRelease(label, xPos, yPos, button, inBox)
	}
}

func (label *Label) IsHovered(xPos, yPos float64) {
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	label.IsHover = inBox
	if inBox {
		label.OnHover(label, xPos, yPos, MouseUnclicked, inBox)
		if label.Shadow != nil {
			label.OnHover(&label.Shadow.Label, xPos, yPos, MouseUnclicked, inBox)
		}
	}
}

func (label *Label) Draw() {
	if label.Shadow != nil {
		label.Shadow.Text.Draw()
	}
	label.Text.Draw()
}
