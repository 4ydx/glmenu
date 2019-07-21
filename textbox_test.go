package glmenu

import (
	"runtime"
	"testing"

	"github.com/4ydx/glfw/v3.3/glfw"
	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var window *glfw.Window

func openGLContext() {
	useStrictCoreProfile := (runtime.GOOS == "darwin")

	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic("glfw error")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	if useStrictCoreProfile {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	}
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	window, err = glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
}

func TestTextBoxBackspace(t *testing.T) {
	openGLContext()

	tb := TextBox{}

	f := &v41.Font{}
	f.Config = &gltext.FontConfig{}

	text := &v41.Text{}
	text.Font = f
	text.SetString("testing")
	tb.Text = text

	text = &v41.Text{}
	text.Font = f
	text.SetString("|")
	tb.Cursor = text

	tb.CursorIndex = 1
	tb.Backspace()
	if text.String != "esting" && tb.CursorIndex != 0 {
		t.Error(tb.Text.String, tb.CursorIndex)
	}
	tb.CursorIndex = 6
	tb.Backspace()
	if text.String != "estin" && tb.CursorIndex != 5 {
		t.Error(tb.Text.String, tb.CursorIndex)
	}
	tb.Backspace()
	if text.String != "esti" && tb.CursorIndex != 4 {
		t.Error(tb.Text.String, tb.CursorIndex)
	}
}
