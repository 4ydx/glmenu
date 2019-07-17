package glmenu

import (
	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

// LabelAction is a predefined action a label can perform
type LabelAction int

const (
	// Noop does nothing
	Noop LabelAction = iota
	// GotoMenu goes to a given menu
	GotoMenu
	// ExitMenu closes the current menu
	ExitMenu
	// ExitGame closes the game
	ExitGame
)

// LabelConfig is the config for a label
type LabelConfig struct {
	Padding Padding
	Action  LabelAction
	Goto    string
}

// LabelInteraction allows the user to define what happens when a label is interacted with
type LabelInteraction func(
	xPos, yPos float64,
	button MouseClick,
	isInBoundingBox bool,
)

// Label is text that can be interacted with
type Label struct {
	Config  LabelConfig
	Menu    *Menu
	Text    *v41.Text
	IsHover bool
	IsClick bool

	// public methods are expected to be defined by the user and run before the private method are called
	// if a public method is undefined, it is skipped.  Currently I have only defined onRelease as private.
	OnClick       LabelInteraction
	onRelease     LabelInteraction // set internally.  handles linking different menus together, closing menus, closing game etc.
	StopOnRelease bool             // set to true in order to prevent onRelease being called
	OnRelease     LabelInteraction
	OnHover       LabelInteraction
	OnNotHover    func()
}

// Reset the text to the original size
func (label *Label) Reset() {
	label.Text.SetScale(label.Text.ScaleMin)
}

// GetPosition returns the label's position
// This is in pixels with the origin at the center of the screen
func (label *Label) GetPosition() mgl32.Vec2 {
	return label.Text.Position
}

// GetPadding returns the padding for the label
func (label *Label) GetPadding() Padding {
	return label.Config.Padding
}

// SetString sets the label's text string
func (label *Label) SetString(str string, argv ...interface{}) {
	if len(argv) == 0 {
		label.Text.SetString(str)
	} else {
		label.Text.SetString(str, argv...)
	}
}

// IsClicked uses a bounding box to determine clicks
func (label *Label) IsClicked(xPos, yPos float64, button MouseClick) {
	mX, mY := ScreenCoordToCenteredCoord(label.Menu.Font.WindowWidth, label.Menu.Font.WindowHeight, xPos, yPos)
	inBox := InBox(mX, mY, label.Text)
	if inBox {
		label.IsClick = true
		if label.OnClick != nil {
			label.OnClick(xPos, yPos, button, inBox)
		}
	}
}

// InsidePoint returns a SCREEN COORDINATE SYSTEM point that is centered inside the label
func (label *Label) InsidePoint() (P gltext.Point) {
	// get the center point
	lowerLeft, upperRight := label.Text.GetBoundingBox()
	x := (upperRight.X + lowerLeft.X) / 2
	y := (upperRight.Y + lowerLeft.Y) / 2

	P.X = x + label.Menu.Font.WindowWidth/2
	P.Y = y + label.Menu.Font.WindowHeight/2
	return
}

// IsReleased is checked for all labels in a menu when mouseup occurs
func (label *Label) IsReleased(xPos, yPos float64, button MouseClick) {
	mX, mY := ScreenCoordToCenteredCoord(label.Menu.Font.WindowWidth, label.Menu.Font.WindowHeight, xPos, yPos)
	inBox := InBox(mX, mY, label.Text)

	// anything flagged as clicked now needs to decide whether to execute its logic based on inX && inY
	if label.IsClick {
		if label.IsHover {
			label.Text.SetColor(label.Menu.Defaults.TextHover)
		} else {
			label.Text.SetColor(label.Menu.Defaults.TextColor)
		}
		if label.OnRelease != nil {
			label.OnRelease(xPos, yPos, button, inBox)
		}
		if !label.StopOnRelease {
			label.onRelease(xPos, yPos, button, inBox)
		}
		label.StopOnRelease = false
	}
	label.IsClick = false
}

// IsHovered uses a bounding box
func (label *Label) IsHovered(xPos, yPos float64) {
	mX, mY := ScreenCoordToCenteredCoord(label.Menu.Font.WindowWidth, label.Menu.Font.WindowHeight, xPos, yPos)
	inBox := InBox(mX, mY, label.Text)
	label.IsHover = inBox
	if inBox {
		label.OnHover(xPos, yPos, MouseUnclicked, inBox)
	}
}

// Draw the label
func (label *Label) Draw() {
	label.Text.Draw()
}

// SetPosition of the label's text
func (label *Label) SetPosition(v mgl32.Vec2) {
	label.Text.SetPosition(v)
}

// DragPosition drags the label along the vector (x,y)
func (label *Label) DragPosition(x, y float32) {
	label.Text.DragPosition(x, y)
}

// Height returns the text height
func (label *Label) Height() float32 {
	return label.Text.Height()
}

// Width returns the text width
func (label *Label) Width() float32 {
	return label.Text.Width()
}

// NavigateTo activates the given label
// Used when arrow keys are used to navigate a menu
func (label *Label) NavigateTo() {
	point := label.InsidePoint()
	if label.OnHover != nil {
		label.IsHovered(float64(point.X), float64(point.Y))
	}
}

// NavigateAway if we end up needing to navigate away from this item then let the caller know
// because it might need that information.  return value of 'true'
func (label *Label) NavigateAway() bool {
	if label.IsHover {
		label.IsHover = false
		return true
	}
	return false
}

// Follow the label by performing a click/release sequence
func (label *Label) Follow() bool {
	if label.IsHover {
		point := label.InsidePoint()
		label.IsClicked(float64(point.X), float64(point.Y), MouseLeft)
		label.IsReleased(float64(point.X), float64(point.Y), MouseLeft)
		return true
	}
	return false
}

// IsNoop indicates if the label is a noop
func (label *Label) IsNoop() bool {
	return label.Config.Action == Noop
}

// Type returns the kind of this object
func (label *Label) Type() FormatableType {
	return FormatableLabel
}
