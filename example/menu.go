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
	menuManager = glmenu.NewMenuManager(font, glfw.KeyM, "main")

	defaults := glmenu.MenuDefaults{
		TextColor:       mgl32.Vec3{1, 1, 1},
		TextClick:       mgl32.Vec3{250.0 / 255.0, 0, 154.0 / 255.0},
		TextHover:       mgl32.Vec3{0, 250.0 / 255.0, 154.0 / 255.0},
		BackgroundColor: mgl32.Vec4{0.5, 0.5, 0.5, 1.0},
		Dimensions:      mgl32.Vec2{10, 10},
		Border:          10,
	}

	// menu 1
	mainMenu, err := menuManager.NewMenu(window, "main", defaults, mgl32.Vec2{100, 0})
	//mainMenu, err := menuManager.NewMenu(window, "main", defaults, mgl32.Vec2{0, 0})
	if err != nil {
		fmt.Println("error loading the font")
		os.Exit(1)
	}
	// 9 different embedded images within image.jpg with indices 0 - 8 running from upper left to lower right
	mainMenu.NewMenuTexture("texture/image.jpg", mgl32.Vec2{3, 3})

	textbox := mainMenu.NewTextBox("127.0.0.1", 250, 40, 1)
	textbox.Text.MaxRuneCount = 16
	mainMenu.NewLabel("Options", glmenu.LabelConfig{Action: glmenu.GOTO_MENU, Goto: "option"})
	mainMenu.NewLabel("Quit", glmenu.LabelConfig{Action: glmenu.EXIT_GAME})
	mainMenu.NewLabel("Dummy", glmenu.LabelConfig{Action: glmenu.NOOP})

	// menu 2
	optionMenu, err := menuManager.NewMenu(window, "option", glmenu.MenuDefaults{BackgroundColor: mgl32.Vec4{0, 1, 1, 1}, Dimensions: mgl32.Vec2{200, 200}}, mgl32.Vec2{})
	if err != nil {
		fmt.Println("error loading font")
		os.Exit(1)
	}
	optionMenu.NewLabel("Back", glmenu.LabelConfig{Action: glmenu.GOTO_MENU, Goto: "main"})

	// complete setup
	menuManager.Finalize()
}
