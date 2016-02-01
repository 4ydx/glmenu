package glmenu

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Border struct {
	X, Y float32
}

type Formatable interface {
	GetPosition() mgl32.Vec2
	SetPosition(v mgl32.Vec2)
	Height() float32
	Width() float32
}
