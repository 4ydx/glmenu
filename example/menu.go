package main

import (
	"github.com/4ydx/glfw/v3.3/glfw"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/mathgl/mgl32"
)

// MenuInit starts menus
func MenuInit(window *glfw.Window, font *v41.Font) {
	menuManager = glmenu.NewMenuManager(window, font, glfw.KeyM, "main")

	border := glmenu.MenuBorder{
		Width: mgl32.Vec2{1, 1},
		Color: mgl32.Vec4{0.5, 0.5, 0.5, 1.0},
	}
	defaults := glmenu.MenuDefaults{
		TextColor:       mgl32.Vec3{1, 1, 1},
		TextClick:       mgl32.Vec3{250.0 / 255.0, 0, 154.0 / 255.0},
		TextHover:       mgl32.Vec3{0.6, 0.6, 0.6},
		BackgroundColor: mgl32.Vec4{0.2, 0.2, 0.2, 1.0},
		Dimensions:      mgl32.Vec2{0, 0},
		Padding:         mgl32.Vec2{10, 10},
		HoverPadding:    mgl32.Vec2{10, 10},
		Border:          border,
		TextScaleRate:   0.00,
	}

	// menu 1
	mainMenu, err := menuManager.NewMenu(window, "main", defaults, glmenu.ScreenLeft)
	if err != nil {
		panic(err)
	}
	textbox, err := mainMenu.NewTextBox("127.0.0.1", 250, 40)
	if err != nil {
		panic(err)
	}
	textbox.Text.MaxRuneCount = 16
	mainMenu.NewLabel("Options", glmenu.LabelConfig{Action: glmenu.GotoMenu, Goto: "option"})
	mainMenu.NewLabel("Dummy", glmenu.LabelConfig{Action: glmenu.Noop})
	mainMenu.NewLabel("Quit", glmenu.LabelConfig{Action: glmenu.ExitGame})

	// menu 2
	defaults = glmenu.MenuDefaults{
		BackgroundColor: mgl32.Vec4{0, 1, 1, 1},
		Dimensions:      mgl32.Vec2{200, 200},
		Padding:         mgl32.Vec2{10, 10},
		HoverPadding:    mgl32.Vec2{10, 10},
	}
	optionMenu, err := menuManager.NewMenu(window, "option", defaults, glmenu.ScreenTopCenter)
	if err != nil {
		panic(err)
	}
	optionMenu.NewLabel("Back", glmenu.LabelConfig{Action: glmenu.GotoMenu, Goto: "main"})

	// complete setup
	if err := menuManager.Finalize(glmenu.AlignRight); err != nil {
		panic(err)
	}
}
