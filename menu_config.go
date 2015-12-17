package glmenu

type Color struct {
	R, G, B float32
}

type MenuConfig struct {
	Color         Color
	ColorClick    Color
	ColorHover    Color
	LabelBorder   Border
	TextBoxBorder Border
}
