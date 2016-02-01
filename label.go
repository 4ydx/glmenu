package glmenu

import (
	"github.com/4ydx/gltext"
	"github.com/go-gl/mathgl/mgl32"
)

type LabelAction int

const (
	NOOP LabelAction = iota
	GOTO_MENU
	EXIT_MENU
	EXIT_GAME
)

type LabelConfig struct {
	Action LabelAction
	Goto   string
}

type LabelInteraction func(
	//label *Label,
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

type Label struct {
	Config  LabelConfig
	Menu    *Menu
	Text    *gltext.Text
	IsHover bool
	IsClick bool

	// defaults exist but can be user defined
	OnClick    LabelInteraction
	OnRelease  LabelInteraction
	OnHover    LabelInteraction
	OnNotHover func()
}

func (label *Label) Reset() {
	label.Text.SetScale(label.Text.ScaleMin)
}

func (label *Label) GetPosition() mgl32.Vec2 {
	return label.Text.Position
}

func (label *Label) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		label.Text.SetString(str)
	} else {
		label.Text.SetString(str, argv)
	}
}

func (label *Label) OrthoToScreenCoord() (X1 Point, X2 Point) {
	if label.Menu != nil && label.Text != nil {
		x1, x2 := label.Text.GetBoundingBox()
		X1.X = x1.X + label.Menu.WindowWidth/2
		X1.Y = x1.Y + label.Menu.WindowHeight/2

		X2.X = x2.X + label.Menu.WindowWidth/2
		X2.Y = x2.Y + label.Menu.WindowHeight/2
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
			//label.OnClick(label, xPos, yPos, button, inBox)
			label.OnClick(xPos, yPos, button, inBox)
		}
	}
}

// IsReleased is checked for all labels in a menu when mouseup occurs
func (label *Label) IsReleased(xPos, yPos float64, button MouseClick) {
	// anything flagged as clicked now needs to decide whether to execute its logic based on inBox
	X1, X2 := label.OrthoToScreenCoord()
	inBox := float32(xPos) > X1.X && float32(xPos) < X2.X && float32(yPos) > X1.Y && float32(yPos) < X2.Y
	if label.IsClick {
		if label.IsHover {
			label.Text.SetColor(label.Menu.Defaults.TextHover)
		} else {
			label.Text.SetColor(label.Menu.Defaults.TextColor)
		}
		if label.OnRelease != nil {
			label.OnRelease(xPos, yPos, button, inBox)
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
		label.OnHover(xPos, yPos, MouseUnclicked, inBox)
	}
}

func (label *Label) Draw() {
	label.Text.Draw()
}

func (label *Label) SetPosition(v mgl32.Vec2) {
	label.Text.SetPosition(v)
}

func (label *Label) Height() float32 {
	return label.Text.Height()
}

func (label *Label) Width() float32 {
	return label.Text.Width()
}
