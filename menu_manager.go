package glmenu

import (
	"errors"
	"fmt"

	"github.com/4ydx/gltext/v4.1"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// MenuManager handles ineraction between related menus (navigation)
// as well as performing other conveniences
type MenuManager struct {
	Font        *v41.Font
	StartKey    glfw.Key // they key that, when pressed, will display the StartMenu
	StartMenu   string   // the name passed to each NewMenu call
	Menus       map[string]*Menu
	IsFinalized bool
}

// Finalize connects menus together and performs final formatting steps
// this must be run after all menus are prepared
func (mm *MenuManager) Finalize(align Alignment) error {
	if mm.IsFinalized {
		return errors.New("Menus have already been finalized")
	}
	for _, menu := range mm.Menus {
		menu.Finalize(align)
		for _, label := range menu.Labels {
			if label.Config.Action == GotoMenu {
				gotoMenu, ok := mm.Menus[label.Config.Goto]
				if ok {
					func(from *Menu, to *Menu, l *Label) {
						l.onRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
							if inBox {
								from.Hide()
								to.Show()
							}
						}
					}(menu, gotoMenu, label)
				}
			}
		}
	}
	mm.IsFinalized = true
	return nil
}

// IsVisible returns true when one of the menus is visible
func (mm *MenuManager) IsVisible() bool {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			return true
		}
	}
	return false
}

// MouseClick processes potential menus clicks
func (mm *MenuManager) MouseClick(xPos, yPos float64, button MouseClick) {
	yPos = float64(mm.Font.WindowHeight) - yPos

	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseClick(xPos, yPos, button)
			return
		}
	}
}

// MouseRelease processes potential mouse releases
func (mm *MenuManager) MouseRelease(xPos, yPos float64, button MouseClick) {
	yPos = float64(mm.Font.WindowHeight) - yPos

	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseRelease(xPos, yPos, button)
			return
		}
	}
}

// MouseHover processes potential mouse hover events
func (mm *MenuManager) MouseHover(xPos, yPos float64) {
	yPos = float64(mm.Font.WindowHeight) - yPos

	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseHover(xPos, yPos)
			return
		}
	}
}

// KeyRelease handles key release events
func (mm *MenuManager) KeyRelease(key glfw.Key, withShift bool) {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.KeyRelease(key, withShift)
			return
		}
	}
}

// Draw the visible menu
func (mm *MenuManager) Draw() bool {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			if menu.OnComplete != nil {
				menu.OnComplete()
			}
			return menu.Draw()
		}
	}
	return false
}

// Release opengl objects in all managed menus
func (mm *MenuManager) Release() {
	for _, menu := range mm.Menus {
		menu.Release()
	}
}

// NewMenu creates a new managed menu
func (mm *MenuManager) NewMenu(window *glfw.Window, name string, menuDefaults MenuDefaults, screenPosition ScreenPosition) (*Menu, error) {
	m, err := NewMenu(window, name, mm.Font, menuDefaults, screenPosition)
	if err != nil {
		return nil, err
	}
	m.MenuManager = mm

	if _, ok := mm.Menus[name]; ok {
		return nil, fmt.Errorf("the named menu %s already exists", name)
	}
	mm.Menus[name] = m
	return m, nil
}

// Hide all menus
func (mm *MenuManager) Hide() {
	for _, m := range mm.Menus {
		m.Hide()
	}
}

// Show the requested menu
func (mm *MenuManager) Show(name string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return fmt.Errorf("the named menu '%s' doesn't exists", name)
	}
	m.Show()
	return nil
}

// Toggle the named menu
func (mm *MenuManager) Toggle(name string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return fmt.Errorf("the named menu '%s' doesn't exists", name)
	}
	m.Toggle()
	return nil
}

// SetText for the label at the given index of the named menu
func (mm *MenuManager) SetText(name string, index int, text string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return fmt.Errorf("the named menu '%s' doesn't exists", name)
	}
	if index < 0 || index >= len(m.Labels) {
		return fmt.Errorf("the index '%d' does not exist", index)
	}
	for i, l := range m.Labels {
		if i == index {
			l.Text.SetString(text)
		}
	}
	return nil
}

// NewMenuManager handles a tree of menus that interact with one another
func NewMenuManager(font *v41.Font, startKey glfw.Key, startMenu string) *MenuManager {
	mm := &MenuManager{Font: font, StartKey: startKey, StartMenu: startMenu}
	mm.Menus = make(map[string]*Menu)
	return mm
}
