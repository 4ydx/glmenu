package menu

import (
	gltext "github.com/4ydx/gltext"
	"time"
)

// 2DO: code to render the box around the text

type TextBox struct {
	Menu               *Menu
	Text               *gltext.Text
	MaxLength          int
	CursorBarFrequency int64
	Time               time.Time
}

func (textbox *TextBox) Load(menu *Menu, font *gltext.Font) {
	textbox.CursorBarFrequency = time.Duration.Nanoseconds(500000000)
	textbox.Menu = menu
	textbox.Text = gltext.LoadText(font)
	textbox.Text.SetString("|") // the infamous flashing bar
}

func (textbox *TextBox) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		textbox.Text.SetString(str)
	} else {
		textbox.Text.SetString(str, argv)
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
	textbox.Text.Draw()
}
