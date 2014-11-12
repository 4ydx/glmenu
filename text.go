package menu

import (
	gltext "github.com/4ydx/gltext"
	"os"
)

type OnClick func(xPos, yPos float64) (err error)
type OnHover func(xPos, yPos float64) (err error)
type Text struct {
	*gltext.Font
	*gltext.Text
	OnClick OnClick
	OnHover OnHover
	IsHover bool
}

func (text *Text) Load(scale, low, high int32) (err error) {
	fd, err := os.Open("font/luximr.ttf")
	if err != nil {
		return
	}
	defer fd.Close()

	text.Font, err = gltext.LoadTruetype(fd, scale, 32, 127)
	if err != nil {
		return
	}
	text.Text, err = gltext.LoadText(text.Font)
	if err != nil {
		return
	}
	return nil
}

func (text *Text) SetString(str string) (gltext.Point, gltext.Point) {
	return text.Text.SetString(text.Font, str)
}

func (text *Text) IsClicked(xPos, yPos float64) {
	if float32(xPos) > text.X1.X && float32(xPos) < text.X2.X && float32(yPos) > text.X1.Y && float32(yPos) < text.X2.Y {
		text.OnClick(xPos, yPos)
	}
}

func (text *Text) IsHovered(xPos, yPos float64) {
	if float32(xPos) > text.X1.X && float32(xPos) < text.X2.X && float32(yPos) > text.X1.Y && float32(yPos) < text.X2.Y {
		text.IsHover = true
		text.OnHover(xPos, yPos)
	} else {
		text.IsHover = false
	}
}
