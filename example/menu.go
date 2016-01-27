package main

import (
	"fmt"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"os"
)

func MenuInit(window *glfw.Window, font *gltext.Font) {
	menuManager = glmenu.NewMenuManager(font, glfw.KeyM)

	defaults := glmenu.MenuDefaults{
		TextColor:       mgl32.Vec3{1, 1, 1},
		TextClick:       mgl32.Vec3{250.0 / 255.0, 0, 154.0 / 255.0},
		TextHover:       mgl32.Vec3{0, 250.0 / 255.0, 154.0 / 255.0},
		BackgroundColor: mgl32.Vec4{0.5, 0.5, 0.5, 1.0},
	}

	// menu 1
	mainMenu, err := menuManager.NewMenu(window, "main", defaults, mgl32.Vec2{})
	if err != nil {
		fmt.Println("error loading the font")
		os.Exit(1)
	}
	textbox := mainMenu.NewTextBox("127.0.0.1", 250, 40, 1)
	textbox.Text.MaxRuneCount = 16
	label := mainMenu.NewLabel("Options", glmenu.NOOP)
	mainMenu.NewLabel("Quit", glmenu.EXIT_GAME)
	mainMenu.NewLabel("Dummy", glmenu.NOOP)

	// menu 2
	optionMenu, err := menuManager.NewMenu(window, "option", glmenu.MenuDefaults{BackgroundColor: mgl32.Vec4{1, 1, 1, 1}}, mgl32.Vec2{})
	if err != nil {
		fmt.Println("error loading font")
		os.Exit(1)
	}
	optionMenu.NewLabel("Back", glmenu.NOOP)

	// navigation
	label.OnRelease = func(xPos, yPos float64, button glmenu.MouseClick, inBox bool) {
		if inBox {
			mainMenu.Hide()
			optionMenu.Show()
		}
	}
}
