package glmenu

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Border struct {
	X, Y float32
}

type Formatable interface {
	NavigateTo()
	NavigateAway() bool
	GetPosition() mgl32.Vec2
	SetPosition(v mgl32.Vec2)
	GetBorder() Border // the padding around the Formatable object
	Height() float32
	Width() float32
}
