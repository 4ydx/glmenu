package glmenu

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// ArrowCursor is a system defined arrow cursor
var ArrowCursor *glfw.Cursor

// HandCursor is a system defined hand cursor
var HandCursor *glfw.Cursor

// CursorInit initializes the cursors
func CursorInit() {
	ArrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	HandCursor = glfw.CreateStandardCursor(glfw.HandCursor)
}

// CursorDestroy releases cursor assets
func CursorDestroy() {
	if ArrowCursor != nil {
		ArrowCursor.Destroy()
	}
	if HandCursor != nil {
		HandCursor.Destroy()
	}
}
