package glmenu

import (
	"github.com/go-gl/mathgl/mgl32"
)

// FormatableType indicates what kind of object this is.
type FormatableType int

const (
	// FormatableLabel is a label
	FormatableLabel = 0
	// FormatableTextbox is a textbox
	FormatableTextbox = 1
)

// Padding indicates x,y padding
type Padding struct {
	X, Y float32
}

// Formatable interface definition
type Formatable interface {
	// perform click action as appropriate
	// if formatable has no reasonable click action (TextBox) returns false
	Follow() bool

	// up/down key navigation
	NavigateTo()
	NavigateAway() bool

	// is the this object something that can be interacted with by the user
	IsNoop() bool

	// rendering
	GetPosition() mgl32.Vec2
	SetPosition(v mgl32.Vec2)
	DragPosition(x, y float32)
	GetPadding() Padding
	Height() float32
	Width() float32

	// get the type
	Type() FormatableType
}
