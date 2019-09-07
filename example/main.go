package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/4ydx/glfw/v3.3/glfw"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/gl/v4.1-core/gl"
	"golang.org/x/image/math/fixed"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

func keyCallback(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	if action != glfw.Release {
		if mods == glfw.ModShift {
			menuManager.KeyRelease(key, true)
		} else {
			menuManager.KeyRelease(key, false)
		}
	}
}

func mouseMoveCallback(w *glfw.Window, xPos, yPos float64) {
	menuManager.MouseMove(xPos, yPos)
	if menuManager.IsHovered() {
		w.SetCursor(glmenu.HandCursor)
	} else {
		w.SetCursor(glmenu.ArrowCursor)
	}
}

func mouseButtonCallback(
	w *glfw.Window,
	button glfw.MouseButton,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	xPos, yPos := w.GetCursorPos()
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		menuManager.MouseClick(xPos, yPos, glmenu.MouseLeft)
	}
	if button == glfw.MouseButtonLeft && action == glfw.Release {
		menuManager.MouseRelease(xPos, yPos, glmenu.MouseLeft)
	}
}

var window *glfw.Window
var menuManager *glmenu.MenuManager

func main() {
	var err error

	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
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

	// fullscreen
	primary := glfw.GetPrimaryMonitor()
	vm := primary.GetVideoMode()
	w, h := vm.Width, vm.Height // you should probably pick one in another manner
	window, err = glfw.CreateWindow(w, h, "Testing", primary, nil)
	// fullscreen

	// windowed
	// window, err = glfw.CreateWindow(640, 480, "Testing", nil, nil)
	// windowed

	if err != nil {
		panic(err)
	}
	glmenu.CursorInit()
	window.MakeContextCurrent()
	window.SetKeyCallback(keyCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(mouseMoveCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version", version)

	var font *v41.Font
	config, err := gltext.LoadTruetypeFontConfig("fontconfigs", "orange")

	if err == nil {
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
		fmt.Println("Font loaded from disk...")
	} else {
		fd, err := os.Open("font/orange kid.ttf")
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = fd.Close()
		}()

		runesPerRow := fixed.Int26_6(16)
		runeRanges := make(gltext.RuneRanges, 0)
		runeRanges = append(runeRanges, gltext.RuneRange{Low: 1, High: 128})

		//scale := fixed.Int26_6(25)
		scale := fixed.Int26_6(18)
		config, err = gltext.NewTruetypeFontConfig(fd, scale, runeRanges, runesPerRow, 5)
		if err != nil {
			panic(err)
		}
		err = config.Save("fontconfigs", "luximr")
		if err != nil {
			panic(err)
		}
		font, err = v41.NewFont(config)
		if err != nil {
			panic(err)
		}
	}

	// load menus
	MenuInit(window, font)
	if err := menuManager.Show("main"); err != nil {
		panic(err)
	}

	gl.ClearColor(0, 0, 0, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		xPos, yPos := window.GetCursorPos()
		menuManager.MouseHover(xPos, yPos)
		if menuManager.Draw() {
			// menuManager.Menus["main"].Drag(1, 0)
			// pause gameplay
		} else {
			// play game
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
	menuManager.Release()
	glmenu.CursorDestroy()
}
