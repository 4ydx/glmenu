package menu

import (
	gltext "github.com/4ydx/gltext"
	"time"
)

// 2DO: code to render the box around the text

type TextBoxInteraction func(
	textbox *TextBox,
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

type TextBox struct {
	Menu               *Menu
	Text               *gltext.Text
	OnClick            TextBoxInteraction
	IsClick            bool
	OnRelease          TextBoxInteraction
	MaxLength          int
	CursorBarFrequency int64
	Time               time.Time
	IsEdit             bool
}

func (textbox *TextBox) Load(menu *Menu, font *gltext.Font) {
	textbox.CursorBarFrequency = time.Duration.Nanoseconds(500000000)
	textbox.Menu = menu
	textbox.Text = gltext.LoadText(font)
}

func (textbox *TextBox) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		textbox.Text.SetString(str + "|")
	} else {
		textbox.Text.SetString(str+"|", argv)
	}
}

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
	// 2DO: draw the border

	textbox.Text.Draw()
}

func (textbox *TextBox) DeleteCharacter() {
	r := []rune(textbox.Text.String)
	r = r[0 : len(r)-2]
	textbox.Text.SetString(string(r) + "|")
}

func (textbox *TextBox) OrthoToScreenCoord() (X1 Point, X2 Point) {
	X1.X = textbox.Text.X1.X + textbox.Menu.WindowWidth/2
	X1.Y = textbox.Text.X1.Y + textbox.Menu.WindowHeight/2

	X2.X = textbox.Text.X2.X + textbox.Menu.WindowWidth/2
	X2.Y = textbox.Text.X2.Y + textbox.Menu.WindowHeight/2
	return
}

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
