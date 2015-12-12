package glmenu

import (
	"github.com/4ydx/gltext"
)

type LabelInteraction func(
	label *Label,
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

type Label struct {
	Menu    *Menu
	Text    *gltext.Text
	Shadow  *Shadow
	IsHover bool
	IsClick bool

	// user defined
	OnClick    LabelInteraction
	OnRelease  LabelInteraction
	OnHover    LabelInteraction
	OnNotHover func(label *Label)
}

type Shadow struct {
	Label
	Offset float32
}

func (label *Label) NewShadow(offset, r, g, b float32) {
	label.Shadow = &Shadow{}
	label.Shadow.Menu = label.Menu
	label.updateShadow(offset, r, g, b)
}

func (label *Label) updateShadow(offset, r, g, b float32) {
	label.Shadow.Text = gltext.NewText(label.Menu.Font, 1.0, 1.1)
	label.Shadow.Text.SetColor(r, g, b)
	label.Shadow.Text.SetString(label.Text.String)
	label.Shadow.Text.SetPosition(label.Text.SetPositionX+offset, label.Text.SetPositionY+offset)

	label.Shadow.OnClick = label.OnClick
	label.Shadow.OnRelease = label.OnRelease
	label.Shadow.OnHover = label.OnHover
	label.Shadow.OnNotHover = label.OnNotHover
}

func (label *Label) Reset() {
	label.Text.SetScale(label.Text.ScaleMin)
	if label.Shadow != nil {
		label.Shadow.Text.SetScale(label.Text.ScaleMin)
	}
}

func (label *Label) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		label.Text.SetString(str)
	} else {
		label.Text.SetString(str, argv)
	}
	if label.Shadow != nil {
		if len(argv) == 0 {
			label.Shadow.Text.SetString(str)
		} else {
			label.Shadow.Text.SetString(str, argv)
		}
	}
}

func (label *Label) OrthoToScreenCoord() (X1 Point, X2 Point) {
	if label.Menu != nil && label.Text != nil {
		X1.X = label.Text.X1.X + label.Menu.WindowWidth/2
		X1.Y = label.Text.X1.Y + label.Menu.WindowHeight/2

		X2.X = label.Text.X2.X + label.Menu.WindowWidth/2
		X2.Y = label.Text.X2.Y + label.Menu.WindowHeight/2
	} else {
		if label.Menu == nil {
			MenuDebug("Uninitialized Menu Object")
		}
		if label.Text == nil {
			MenuDebug("Uninitialized Text Object")
		}
	}
	return
}

// IsClicked uses a bounding box to determine clicks
func (label *Label) IsClicked(xPos, yPos float64, button MouseClick) {
	// menu rendering (and text) is positioned in orthographic projection coordinates
	// but click positions are based on window coordinates
	// we have to transform them
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if inBox {
		label.IsClick = true
		if label.OnClick != nil {
			label.OnClick(label, xPos, yPos, button, inBox)
		}
	}
}

// IsReleased is checked for all labels in a menu when mouseup occurs
func (label *Label) IsReleased(xPos, yPos float64, button MouseClick) {
	// anything flagged as clicked now needs to decide whether to execute its logic based on inBox
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if label.IsClick {
		if label.OnRelease != nil {
			label.OnRelease(label, xPos, yPos, button, inBox)
		}
	}
	label.IsClick = false
}

// IsHovered uses a bounding box
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
