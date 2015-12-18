package glmenu

type MenuManager struct {
	Menus []Menu
}

// ResolveNavigation connects menus together
// this must be run after all menus are prepared
func (mm *MenuManager) ResolveNavigation() {
}

func (mm *MenuManager) Clicked() {
}

func (mm *MenuManager) Toggle() {
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
