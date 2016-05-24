package glmenu

import (
	"errors"
	"fmt"
	"github.com/4ydx/gltext"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type MenuManager struct {
	Font       *gltext.Font
	StartKey   glfw.Key // they key that, when pressed, will display the StartMenu
	StartMenu  string   // the name passed to each NewMenu call
	Menus      map[string]*Menu
	IsResolved bool
}

// Finalize connects menus together and performs final formatting steps
// this must be run after all menus are prepared
func (mm *MenuManager) Finalize() error {
	if mm.IsResolved {
		return errors.New("Menus have already been resolved")
	}
	for _, menu := range mm.Menus {
		menu.Finalize()
		for _, label := range menu.Labels {
			if label.Config.Action == GOTO_MENU {
				gotoMenu, ok := mm.Menus[label.Config.Goto]
				if ok {
					func(m *Menu, to *Menu, l *Label) {
						l.OnRelease = func(xPos, yPos float64, button MouseClick, inBox bool) {
							if inBox {
								m.Hide()
								to.Show()
							}
						}
					}(menu, gotoMenu, label)
				}
			}
		}
	}
	mm.IsResolved = true
	return nil
}

func (mm *MenuManager) IsVisible() bool {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			return true
		}
	}
	return false
}

// Clicked resolves menus that have been clicked
func (mm *MenuManager) MouseClick(xPos, yPos float64, button MouseClick) {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseClick(xPos, yPos, button)
		}
	}
}

func (mm *MenuManager) MouseRelease(xPos, yPos float64, button MouseClick) {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseRelease(xPos, yPos, button)
		}
	}
}

func (mm *MenuManager) MouseHover(xPos, yPos float64) {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			menu.MouseHover(xPos, yPos)
		}
	}
}

func (mm *MenuManager) Draw() bool {
	for _, menu := range mm.Menus {
		if menu.IsVisible {
			return menu.Draw()
		}
	}
	return false
}

func (mm *MenuManager) Release() {
	for _, menu := range mm.Menus {
		menu.Release()
	}
}

func (mm *MenuManager) NewMenu(window *glfw.Window, name string, menuDefaults MenuDefaults, offsetBy mgl32.Vec2) (*Menu, error) {
	m, err := NewMenu(window, mm.Font, menuDefaults, offsetBy)
	if err != nil {
		return nil, err
	}
	if _, ok := mm.Menus[name]; ok {
		return nil, errors.New(fmt.Sprintf("The named menu %s already exists.", name))
	}
	mm.Menus[name] = m
	return m, nil
}

func (mm *MenuManager) Show(name string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return errors.New(fmt.Sprintf("The named menu '%s' doesn't exists.", name))
	}
	m.Show()
	return nil
}

func (mm *MenuManager) Toggle(name string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return errors.New(fmt.Sprintf("The named menu '%s' doesn't exists.", name))
	}
	m.Toggle()
	return nil
}

func (mm *MenuManager) SetText(name string, index int, text string) error {
	m, ok := mm.Menus[name]
	if !ok {
		return errors.New(fmt.Sprintf("The named menu '%s' doesn't exists.", name))
	}
	for i, l := range m.Labels {
		if i == index {
			l.Text.SetString(text)
		}
	}
	return nil
}

// NewMenuManager handles a tree of menus that interact with one another
func NewMenuManager(font *gltext.Font, startKey glfw.Key, startMenu string) *MenuManager {
	mm := &MenuManager{Font: font, StartKey: startKey, StartMenu: startMenu}
	mm.Menus = make(map[string]*Menu)
	return mm
}

/*
menuManager := NewMenuManager()

m1 := menuManager.NewMenu("Name1", menuConfig)
m1.NewLabel(..., "")
m1.NewLabel(..., "Name2")

m2 := menuManager.NewMenu("Name2", menuConfig)
m2.NewLabel(..., "")
m2.NewLabel(..., "Name1")

menuManager.ResolveNavigation()
*/
