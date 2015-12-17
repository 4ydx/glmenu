package main

import (
	"fmt"
	"github.com/4ydx/glmenu"
	"github.com/4ydx/gltext"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"os"
)

func optionMenuInit(window *glfw.Window, font *gltext.Font) (err error) {
	//topMargin := float32(50)
	//leftMargin := float32(30)
	width, height := window.GetSize()
	optionMenu, err = glmenu.NewMenu(font, float32(width), float32(height), mgl32.Vec2{})
	if err != nil {
		fmt.Println("error loading font")
		os.Exit(1)
	}
	optionMenu.ResizeWindow(float32(width), float32(height))
	optionMenu.Background = mgl32.Vec4{1, 1, 1, 1}
	optionMenu.TextScaleRate = 0.05

	/*
		label1 := optionMenu.NewLabel("Music Volume")
		label1.NewShadow(1.5, 0, 0, 0)
		label1.Text.SetColor(0.5, 0.5, 0.5)
		label1.OnClick = func(xPos, yPos float64, button glmenu.MouseClick, isBox bool) {
			fmt.Println("clicked", xPos, yPos)
		}
		label1.OnHover = func(xPos, yPos float64, button glmenu.MouseClick, isBox bool) {
			label1.Text.AddScale(optionMenu.TextScaleRate)
			if label1.Shadow != nil {
				label1.Shadow.Text.AddScale(optionMenu.TextScaleRate)
			}
		}
		label1.OnNotHover = func() {
			label1.Text.AddScale(-optionMenu.TextScaleRate)
			if label1.Shadow != nil {
				label1.Shadow.Text.AddScale(-optionMenu.TextScaleRate)
			}
		}
	*/

	label3 := optionMenu.NewLabel("Back")
	//label3.NewShadow(1.5, 0, 0, 0)
	label3.Text.SetColor(0.5, 0.5, 0.5)
	label3.OnClick = func(xPos, yPos float64, button glmenu.MouseClick, isBox bool) {
		optionMenu.Toggle()
		mainMenu.Toggle()
	}
	label3.OnHover = func(xPos, yPos float64, button glmenu.MouseClick, isBox bool) {
		label3.Text.AddScale(optionMenu.TextScaleRate)
		if label3.Shadow != nil {
			label3.Shadow.Text.AddScale(optionMenu.TextScaleRate)
		}
	}
	label3.OnNotHover = func() {
		label3.Text.AddScale(-optionMenu.TextScaleRate)
		if label3.Shadow != nil {
			label3.Shadow.Text.AddScale(-optionMenu.TextScaleRate)
		}
	}

	//label1.Text.SetPosition(-float32(width)/2.0+leftMargin, float32(height)/2.0-topMargin)
	//label1.Text.Justify(gltext.AlignRight)
	//label3.Text.SetPosition(-float32(width)/2.0+leftMargin, float32(height)/2.0-topMargin-label1.Text.Height)
	label3.Text.Justify(gltext.AlignRight)
	return
}
