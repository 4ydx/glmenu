package menu

import (
	gltext "github.com/4ydx/gltext"
)

type OnClick func(xPos, yPos float64) (err error)
type OnHover func(xPos, yPos float64) (err error)
type Label struct {
	Menu    *Menu
	Text    *gltext.Text
	OnClick OnClick
	OnHover OnHover
	IsHover bool
}

func (label *Label) Load(menu *Menu, font *gltext.Font) {
	label.Menu = menu
	label.Text = gltext.LoadText(font)
}

func (label *Label) SetString(str string) {
	label.Text.SetString(str)
}

func (label *Label) OrthoToScreenCoord() (X1 Point, X2 Point) {
	X1.X = label.Text.X1.X + label.Menu.WindowWidth/2
	X1.Y = label.Text.X1.Y + label.Menu.WindowHeight/2

	X2.X = label.Text.X2.X + label.Menu.WindowWidth/2
	X2.Y = label.Text.X2.Y + label.Menu.WindowHeight/2
	return
}

func (label *Label) IsClicked(xPos, yPos float64) {
	// menu rendering (and text) is positioned in orthographic projection coordinates but click positions are based on window coordinates
	// we have to transform them
	X1, X2 := label.OrthoToScreenCoord()
	if float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y {
		label.OnClick(xPos, yPos)
	}
}

func (label *Label) IsHovered(xPos, yPos float64) {
	X1, X2 := label.OrthoToScreenCoord()
	if float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y {
		label.IsHover = true
		label.OnHover(xPos, yPos)
	} else {
		label.IsHover = false
	}
}
