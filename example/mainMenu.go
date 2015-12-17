package main

import (
	"fmt"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"os"
)

func mainMenuInit(window *glfw.Window, font *gltext.Font) (err error) {
	// actually we are supposed to pass in the framebuffer sizes when creating the orthographic projection
	// this would probably require some changes though in order to track mouse movement.  simply passing
	// in "w" and "h" below into the Load methods (and resize methods) results in menues that are no longer clickable
	width, height := window.GetSize()
	fmt.Println("the window is reporting w", width, "h", height)
	//w, h := window.GetFramebufferSize()
	//fmt.Println("the window is reporting w", w, "h", h)

	mainMenu, err = glmenu.NewMenu(font, float32(width), float32(height), mgl32.Vec2{})
	mainMenu.LabelBorder = glmenu.Border{5, 5}
	if err != nil {
		fmt.Println("error loading the font")
		os.Exit(1)
	}
	mainMenu.ResizeWindow(float32(width), float32(height))
	mainMenu.Background = mgl32.Vec4{0, 0, .20, 0}

	// start
	/*
		textbox1 := mainMenu.NewTextBox("127.0.0.1", 250, 40, 1)
		textbox1.SetColor(1, 1, 1)
		textbox1.Text.MaxRuneCount = 16
	*/

	// options
	label2 := mainMenu.NewLabel("Options")
	label2.Text.SetColor(0.5, 0.5, 0.5)

	label2.OnClick = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		label2.Text.SetColor(250.0/255.0, 0, 154.0/255.0)
	}
	label2.OnRelease = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		label2.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
		if inBox {
			mainMenu.Toggle()
			optionMenu.Toggle()
		}
		if label2.IsHover {
			label2.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
		} else {
			label2.Text.SetColor(0.5, 0.5, 0.5)
		}
	}
	label2.OnHover = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		if !label2.IsClick {
			label2.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
			label2.Text.AddScale(mainMenu.TextScaleRate)
		}
	}
	label2.OnNotHover = func() {
		if !label2.IsClick {
			label2.Text.SetColor(0.5, 0.5, 0.5)
			label2.Text.AddScale(-mainMenu.TextScaleRate)
		}
	}

	// quit
	label3 := mainMenu.NewLabel("Quit")
	label3.Text.SetColor(0.5, 0.5, 0.5)

	label3.OnClick = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		label3.Text.SetColor(250.0/255.0, 0, 154.0/255.0)
	}
	label3.OnRelease = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		label3.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
		if inBox {
			window.SetShouldClose(true)
		}
		if label3.IsHover {
			label3.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
		} else {
			label3.Text.SetColor(0.5, 0.5, 0.5)
		}
	}
	label3.OnHover = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		if !label3.IsClick {
			label3.Text.SetColor(0, 250.0/255.0, 154.0/255.0)
			label3.Text.AddScale(mainMenu.TextScaleRate)
		}
	}
	label3.OnNotHover = func() {
		if !label3.IsClick {
			label3.Text.SetColor(0.5, 0.5, 0.5)
			label3.Text.AddScale(-mainMenu.TextScaleRate)
		}
	}

	label4 := mainMenu.NewLabel("Dummy")
	label4.Text.SetColor(0.5, 0.5, 0.5)

	// simple centering of values
	/*
		totalHeight := textbox1.Text.X2.Y - textbox1.Text.X1.Y +
			label2.Text.X2.Y - label2.Text.X1.Y +
			label2.Text.X2.Y - label2.Text.X1.Y
		textbox1.SetPosition(0, totalHeight/2)

		label3.Text.SetPosition(0, -totalHeight/2)
	*/

	return
}
