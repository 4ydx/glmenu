package main

import (
	"fmt"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"golang.org/x/image/math/fixed"
	"os"
	"runtime"
)

var useStrictCoreProfile = (runtime.GOOS == "darwin")

func keyCallback(
	w *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	/*
		if mainMenu.IsVisible && action == glfw.Release {
			if mods == glfw.ModShift {
				mainMenu.KeyRelease(key, true)
			} else {
				mainMenu.KeyRelease(key, false)
			}
		} else {
			if key == glfw.KeyM && action == glfw.Press {
				if optionMenu.IsVisible {
					optionMenu.Toggle()
				}
				mainMenu.Toggle()
			}
			if key == glfw.KeyO && action == glfw.Press {
				if !mainMenu.IsVisible {
					optionMenu.Toggle()
				}
			}
		}
	*/
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
	err = glfw.Init()
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
	window.MakeContextCurrent()
	window.SetKeyCallback(keyCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Opengl version", version)

	// load font
	fd, err := os.Open("font/luximr.ttf")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	font, err := gltext.LoadTruetype("fontconfigs")
	if err == nil {
		fmt.Println("Font loaded from disk...")
	} else {
		runesPerRow := fixed.Int26_6(16)
		runeRanges := make(gltext.RuneRanges, 0)
		runeRange := gltext.RuneRange{Low: 1, High: 128}
		runeRanges = append(runeRanges, runeRange)

		scale := fixed.Int26_6(25)
		font, err = gltext.NewTruetype(fd, scale, runeRanges, runesPerRow)
		if err != nil {
			panic(err)
		}
		err = font.Config.Save("fontconfigs")
		if err != nil {
			panic(err)
		}
	}

	// load menus
	MenuInit(window, font)
	menuManager.Show("main")

	gl.ClearColor(0, 0, 0, 0.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		xPos, yPos := window.GetCursorPos()
		menuManager.MouseHover(xPos, yPos)
		if menuManager.Draw() {
			// pause gameplay
		} else {
			// do stuff
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
	menuManager.Release()
}
