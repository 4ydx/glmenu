package glmenu

import (
	"github.com/4ydx/gltext"
)

// BoundingBox interface expectes a get method
type BoundingBox interface {
	GetBoundingBox() (X1, X2 gltext.Point)
}

// ScreenCoordToCenteredCoord transforms the mouse coordinate to the coordinate system with (0,0) at the center of the screen
func ScreenCoordToCenteredCoord(width, height float32, posX, posY float64) (float32, float32) {
	return float32(posX) - width/2, float32(posY) - height/2
}

// InBox determines if the mouse position, which must be relative to a coordinate system with the (0,0) origin at the center
// of the screen, falls within the bounding box of the given object
func InBox(mX, mY float32, obj BoundingBox) bool {
	lowerLeft, upperRight := obj.GetBoundingBox()
	inX := mX > lowerLeft.X && mX < upperRight.X
	inY := mY > lowerLeft.Y && mY < upperRight.Y
	return inX && inY
}
