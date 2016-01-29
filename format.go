package glmenu

type Border struct {
	X, Y float32
}

type Formatable interface {
	SetPosition(x, y float32)
	Height() float32
	Width() float32
}
