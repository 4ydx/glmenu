package menu

import (
	gltext "github.com/4ydx/gltext"
)

type OnClick func(xPos, yPos float64) (err error)
type OnHover func(xPos, yPos float64) (err error)
type Label struct {
	Text    *gltext.Text
	Menu    *Menu
	OnClick OnClick
	OnHover OnHover
	IsHover bool
}

func (label *Label) Load(font *gltext.Font) {
	label.Text = gltext.LoadText(font)
}

func (label *Label) SetString(str string) (gltext.Point, gltext.Point) {
	return label.Text.SetString(str)
}

func (label *Label) IsClicked(xPos, yPos float64) {
	if float32(xPos) > label.Text.X1.X && float32(xPos) < label.Text.X2.X && float32(yPos) > label.Text.X1.Y && float32(yPos) < label.Text.X2.Y {
		label.OnClick(xPos, yPos)
	}
}

func (label *Label) IsHovered(xPos, yPos float64) {
	if float32(xPos) > label.Text.X1.X && float32(xPos) < label.Text.X2.X && float32(yPos) > label.Text.X1.Y && float32(yPos) < label.Text.X2.Y {
		label.IsHover = true
		label.OnHover(xPos, yPos)
	} else {
		label.IsHover = false
	}
}
