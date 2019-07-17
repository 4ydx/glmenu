package glmenu

import (
	"github.com/go-gl/mathgl/mgl32"
)

type FormatableType int

const (
	FormatableLabel   = 0
	FormatableTextbox = 1
)

type Padding struct {
	X, Y float32
}

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

	Type() FormatableType
}
